#!/usr/bin/env node

/* eslint-env node */

/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * OpenAPI-based form config generator
 *
 * Usage:
 *   pnpm generate:form-config <methodName>
 *
 * Example:
 *   pnpm generate:form-config mountsEnableSecretsEngine
 *
 * The API class is automatically determined from the OpenAPI spec's tags.
 */

import fs from 'node:fs';
import path from 'node:path';
import { execSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { prepFormConfig, generateConfigContent } from '../app/utils/form-config-generator.js';
import { dasherize } from '@ember/string';

// ES module workaround to get absolute paths for __dirname,
// ensuring this script can work regardless of where it's run from
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const OPENAPI_PATH = path.join(__dirname, '../node_modules/@hashicorp/vault-client-typescript/openapi.json');
const OUTPUT_DIR = path.join(__dirname, '../app/forms/v2/generated');

const normalize = (msg) => {
  return msg.trim().replace(/\n\s+/g, '\n');
};

/**
 * Compose functions left-to-right, passing output of each as input to the next
 * Example: pipe(fn1, fn2, fn3)(x) === fn3(fn2(fn1(x)))
 */
const pipe = (...fns) => {
  return (initial) => fns.reduce((value, fn) => fn(value), initial);
};

const parseArgs = () => {
  const args = process.argv.slice(2);
  const method = args[0];

  if (!method || method.startsWith('--')) {
    const err = normalize(`
      ❌ Missing required argument: methodName
      💡 Usage: pnpm generate:form-config <methodName>
      💡 Example: pnpm generate:form-config mountsEnableSecretsEngine
    `);

    console.error(err);
    process.exit(1);
  }

  return { method };
};

const loadOpenAPISpec = () => {
  try {
    const spec = JSON.parse(fs.readFileSync(OPENAPI_PATH, 'utf-8'));
    console.log(`✅ Found openapi.json with ${Object.keys(spec.paths).length} paths...`);
    return spec;
  } catch (error) {
    console.error(`❌ Error loading openapi.json: ${error.message}`);
    process.exit(1);
  }
};

const prepFormConfigWithLogging = (spec, methodName) => {
  const operationId = dasherize(methodName);
  console.log(`🔍 Searching for ${operationId} operation in openapi.json...`);

  const config = prepFormConfig(spec, methodName);
  if (!config) {
    const err = normalize(`
      ❌ Operation "${operationId}" not found in openapi.json
      💡 If this is a plugin-based method, ensure the plugin is enabled and regenerate openapi.json
    `);
    console.error(err);
    process.exit(1);
  }

  if (!config.apiClass) {
    const err = normalize(`
      ❌ Could not determine API class for "${operationId}"
      💡 The operation may be missing tags in the OpenAPI spec
    `);
    console.error(err);
    process.exit(1);
  }

  console.log(`🔨 Building form config for ${config.name}... \n`);
  return config;
};

const writeAndFormat = (content, methodName) => {
  const filename = `${dasherize(methodName)}-config.ts`;
  const filePath = path.join(OUTPUT_DIR, filename);

  fs.writeFileSync(filePath, content, 'utf-8');
  execSync(`pnpm prettier --write "${filePath}"`, { stdio: 'pipe' });

  return filename;
};

const main = () => {
  const startTime = Date.now();
  const { method } = parseArgs();
  console.log(`⚡️ Generating form config for ${method} method...\n`);

  const filename = pipe(
    (spec) => prepFormConfigWithLogging(spec, method),
    (config) => generateConfigContent(config),
    (content) => writeAndFormat(content, method)
  )(loadOpenAPISpec());

  const duration = ((Date.now() - startTime) / 1000).toFixed(2);
  const msg = normalize(`
    ✨ Done! Completed in ${duration}s
    Output: app/forms/v2/generated/${filename}
  `);

  console.log(msg);
};

main();
