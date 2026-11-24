#!/usr/bin/env node

/* eslint-env node */

/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Terraform mapping scaffold generator.
 *
 * Runs in one of two modes depending on the argument:
 *
 * OpenAPI mode  — pass a Vault API path (starts with "/"):
 *   Reads field definitions from the Vault OpenAPI spec and the Terraform
 *   provider schema, cross-references them, and emits a fully typed scaffold.
 *   Requires terraform CLI and an entry in app/utils/terraform-resource-map.ts.
 *   Falls back to interactive mode if any prerequisite is missing.
 *
 *   pnpm generate:terraform-mapping /sys/policies/acl/{name}
 *
 * Interactive mode — pass a camelCase method name (or omit the argument):
 *   Prompts for the Terraform resource type, feature key, and field
 *   definitions one at a time. No external tools required.
 *
 *   pnpm generate:terraform-mapping mountsEnableSecretsEngine
 */

import fs from 'node:fs';
import path from 'node:path';
import readline from 'node:readline';
import { execSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { createRequire } from 'node:module';
import { dasherize, camelize } from '@ember/string';
import {
  crossReference,
  generateInteractiveScaffold,
  generateOpenApiScaffold,
} from '../app/utils/terraform-code-generator/mapping-generator.ts';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const UI_DIR = path.join(__dirname, '..');
const OUTPUT_DIR = path.join(UI_DIR, 'app/utils/terraform-mappings');
const RESOURCE_MAP_PATH = path.join(UI_DIR, 'app/utils/terraform-code-generator/terraform-resource-map.ts');
// Kept outside the repo so it is never accidentally committed
const TF_PROBE_DIR = path.join('/tmp', 'terraform-schema-probe');

const normalize = (msg) => msg.trim().replace(/\n\s+/g, '\n');

const log = {
  info: (msg) => process.stdout.write(msg + '\n'),
  warn: (msg) => process.stderr.write(msg + '\n'),
  error: (msg) => process.stderr.write(msg + '\n'),
};

// ---------------------------------------------------------------------------
// File I/O
// ---------------------------------------------------------------------------

const writeAndFormat = (content, filename) => {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  const filePath = path.join(OUTPUT_DIR, filename);
  fs.writeFileSync(filePath, content, 'utf-8');
  try {
    execSync(`pnpm prettier --write "${filePath}"`, { cwd: UI_DIR, stdio: 'pipe' });
  } catch {
    log.warn('⚠️  Prettier formatting skipped — format manually before committing.');
  }
  return filePath;
};

// ---------------------------------------------------------------------------
// OpenAPI mode
// ---------------------------------------------------------------------------

const readResourceMap = () => {
  if (!fs.existsSync(RESOURCE_MAP_PATH)) return {};
  const src = fs.readFileSync(RESOURCE_MAP_PATH, 'utf-8');
  const matches = [...src.matchAll(/'([^']+)':\s*'([^']+)'/g)];
  return Object.fromEntries(matches.map(([, k, v]) => [k, v]));
};

const loadOpenApiSpec = () => {
  const require = createRequire(import.meta.url);
  try {
    const pkgPath = require.resolve('@hashicorp/vault-client-typescript/openapi.json');
    return JSON.parse(fs.readFileSync(pkgPath, 'utf-8'));
  } catch {
    return null;
  }
};

// OpenAPI specs avoid repeating type definitions inline — instead a field points to a shared
// definition elsewhere in the document: { "$ref": "#/components/schemas/PolicyWriteRequest" }.
// resolveRef strips the leading "#/", splits on "/" to get individual keys, and walks the spec
// object one key at a time to return the actual schema definition.
const resolveRef = (spec, ref) => {
  const parts = ref.replace('#/', '').split('/');
  return parts.reduce((obj, key) => obj?.[key], spec);
};

const extractOpenApiFields = (spec, apiPath) => {
  const pathDef = spec.paths?.[apiPath];
  if (!pathDef) return null;

  const fields = [];

  for (const param of pathDef.parameters ?? []) {
    if (param.in === 'path') {
      fields.push({
        name: param.name,
        openApiType: param.schema?.type ?? 'string',
        description: param.description ?? '',
        required: param.required ?? false,
        source: 'path-param',
      });
    }
  }

  // Drill into the POST operation's requestBody to get the JSON schema for the payload.
  // If the schema is a $ref pointer, resolveRef fetches the actual definition.
  // Then iterate the schema's properties to collect each field's name, type, description,
  // and whether it appears in the required array.
  const postBody = pathDef.post?.requestBody?.content?.['application/json']?.schema;
  if (postBody) {
    const resolved = postBody.$ref ? resolveRef(spec, postBody.$ref) : postBody;
    const props = resolved?.properties ?? {};
    const requiredFields = resolved?.required ?? [];
    for (const [name, def] of Object.entries(props)) {
      fields.push({
        name,
        openApiType: def.type ?? 'string',
        description: def.description ?? '',
        required: requiredFields.includes(name),
        source: 'request-body',
      });
    }
  }

  return fields;
};

const isTerraformAvailable = () => {
  try {
    execSync('terraform version', { stdio: 'pipe' });
    return true;
  } catch {
    return false;
  }
};

const fetchTerraformSchema = (tfResource) => {
  log.info(`Fetching Terraform provider schema for ${tfResource}...`);

  fs.mkdirSync(TF_PROBE_DIR, { recursive: true });
  fs.writeFileSync(
    path.join(TF_PROBE_DIR, 'main.tf'),
    `terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.0"
    }
  }
}
`
  );

  if (!fs.existsSync(path.join(TF_PROBE_DIR, '.terraform'))) {
    log.info('   Initializing Terraform provider (one-time download)...');
    execSync('terraform init -no-color', { cwd: TF_PROBE_DIR, stdio: 'pipe' });
  }

  const schemaJson = execSync('terraform providers schema -json', {
    cwd: TF_PROBE_DIR,
    stdio: 'pipe',
  }).toString();

  const schema = JSON.parse(schemaJson);
  return (
    schema.provider_schemas?.['registry.terraform.io/hashicorp/vault']?.resource_schemas?.[tfResource]?.block
      ?.attributes ?? null
  );
};

const runOpenApiMode = async (apiPath) => {
  log.info(`\n⚡️ OpenAPI Terraform mapping generator\n   API path: ${apiPath}\n`);

  // 1. Check resource map
  const resourceMap = readResourceMap();
  const tfResource = resourceMap[apiPath];
  if (!tfResource) {
    log.warn(
      normalize(`
      ⚠️  No Terraform resource found for "${apiPath}" in app/utils/terraform-resource-map.ts.
      💡 Add an entry:
         '${apiPath}': 'vault_<resource_type>',
      💡 Then re-run, or continue in interactive mode.
    `)
    );
    return false;
  }
  log.info(`🔍 Resolved: ${apiPath} → ${tfResource}`);

  // 2. Check OpenAPI spec
  const spec = loadOpenApiSpec();
  if (!spec) {
    log.warn(
      normalize(`
      ⚠️  Could not load @hashicorp/vault-client-typescript/openapi.json.
      💡 The package is listed in package.json but may not be installed yet.
      💡 Run: pnpm install
      💡 Falling back to interactive mode.
    `)
    );
    return false;
  }
  const openApiFields = extractOpenApiFields(spec, apiPath);
  if (!openApiFields) {
    log.warn(`⚠️  Path "${apiPath}" not found in OpenAPI spec. Falling back to interactive mode.`);
    return false;
  }
  log.info(`📋 OpenAPI fields: ${openApiFields.map((f) => f.name).join(', ')}`);

  // 3. Check terraform CLI
  if (!isTerraformAvailable()) {
    log.warn(
      normalize(`
      ⚠️  terraform not found in PATH — cannot fetch provider schema.
      💡 Install: brew install terraform
      💡 Falling back to interactive mode.
    `)
    );
    return false;
  }

  // 4. Fetch TF schema
  let tfAttributes;
  try {
    tfAttributes = fetchTerraformSchema(tfResource);
  } catch (e) {
    log.warn(`⚠️  Failed to fetch Terraform schema: ${e.message}. Falling back to interactive mode.`);
    return false;
  }
  if (!tfAttributes) {
    log.warn(
      `⚠️  Resource "${tfResource}" not found in Terraform provider schema. Falling back to interactive mode.`
    );
    return false;
  }

  // 5. Cross-reference and emit
  const { matched, inOpenApiOnly, inTfOnly } = crossReference(openApiFields, tfAttributes);

  if (inOpenApiOnly.length) {
    log.info(
      `\n⚠️  In Vault API but not in ${tfResource} (omitted):\n   ${inOpenApiOnly
        .map((f) => f.name)
        .join(', ')}`
    );
  }
  if (inTfOnly.length) {
    log.info(
      `\nℹ️  In ${tfResource} but not in Vault API request body (omitted):\n   ${inTfOnly.join(', ')}`
    );
  }
  log.info(`\n✅ Matched fields: ${matched.map((f) => `${f.name} (${f.fieldType})`).join(', ')}`);

  const operationId = apiPath.replace(/\//g, '-').replace(/[{}]/g, '').replace(/^-/, '').replace(/-$/, '');
  const content = generateOpenApiScaffold({ apiPath, tfResource, matched, inOpenApiOnly, inTfOnly });
  const filename = `${dasherize(operationId)}-mapping.ts`;
  writeAndFormat(content, filename);
  return { filename };
};

// ---------------------------------------------------------------------------
// Interactive mode - default fallback if OpenAPI mode fails or no API path is provided
// ---------------------------------------------------------------------------

let rl;
const getReadline = () => {
  if (!rl) rl = readline.createInterface({ input: process.stdin, output: process.stdout });
  return rl;
};
const ask = (question) => new Promise((resolve) => getReadline().question(question, resolve));

const FIELD_TYPES = ['string', 'boolean', 'number', 'object', 'array', 'heredoc'];

const collectInputs = async (method) => {
  log.info(`\n⚡️ Terraform mapping generator (interactive)\n   Method: ${method}\n`);

  const tfResource = (await ask('Terraform resource type (e.g. vault_mount): ')).trim();
  if (!tfResource) {
    log.error('❌ Terraform resource type is required.');
    getReadline().close();
    process.exit(1);
  }

  const featureKey = (await ask('Registry feature key (e.g. secrets/kv): ')).trim();
  if (!featureKey) {
    log.error('❌ Feature key is required.');
    getReadline().close();
    process.exit(1);
  }

  const fields = [];
  log.info('\n📋 Add fields one at a time. Press enter with no name when done.\n');

  // eslint-disable-next-line no-constant-condition
  while (true) {
    const name = (await ask('  Field name (or enter to finish): ')).trim();
    if (!name) break;

    const rawType = (await ask(`  Type (${FIELD_TYPES.join(' | ')}) [string]: `)).trim() || 'string';
    const type = FIELD_TYPES.includes(rawType) ? rawType : 'string';
    const required = (await ask('  Required? (y/N): ')).trim().toLowerCase() === 'y';

    fields.push({ name, type, required });
    log.info(`  ✅ Added ${name}: ${type}\n`);
  }

  getReadline().close();

  if (!fields.length) {
    log.error('❌ At least one field is required.');
    getReadline().close();
    process.exit(1);
  }

  return { method, tfResource, featureKey, fields };
};

const runInteractiveMode = async (method) => {
  const inputs = await collectInputs(method);
  const content = generateInteractiveScaffold(inputs);
  const filename = `${dasherize(method)}-mapping.ts`;
  writeAndFormat(content, filename);
  return { filename };
};

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

const main = async () => {
  const startTime = Date.now();
  const arg = process.argv[2];

  let result;

  if (arg?.startsWith('/')) {
    // Looks like an API path — try OpenAPI mode, fall back to interactive
    result = await runOpenApiMode(arg);
    if (!result) {
      // Derive a method name from the path for the interactive fallback
      const fallbackMethod = camelize(
        arg.replace(/\//g, '-').replace(/[{}]/g, '').replace(/^-/, '').replace(/-$/, '')
      );
      log.info(`\n🔄 Continuing in interactive mode (method: ${fallbackMethod})\n`);
      result = await runInteractiveMode(fallbackMethod);
    }
  } else {
    // Method name or no arg — interactive mode directly
    const method = arg || 'terraformMapping';
    if (!arg) {
      log.info(
        normalize(`
        💡 No argument provided. Running in interactive mode.
        💡 Tip: pass a Vault API path for OpenAPI-driven generation:
           pnpm generate:terraform-mapping /sys/policies/acl/{name}
      `)
      );
    }
    result = await runInteractiveMode(method);
    getReadline().close();
  }

  const duration = ((Date.now() - startTime) / 1000).toFixed(2);
  log.info(
    normalize(`
    ✨ Done! Completed in ${duration}s
    Output: app/utils/terraform-mappings/${result.filename}
    Next: review the scaffold, trim unneeded fields, then register in terraform-registry.ts
  `)
  );
};

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  main();
}
