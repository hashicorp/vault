/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { dasherize, camelize, classify } from '@ember/string';

export type FieldType = 'string' | 'boolean' | 'number' | 'object' | 'array' | 'heredoc';

export interface OpenApiField {
  name: string;
  openApiType: string;
  description: string;
  required: boolean;
  source: 'path-param' | 'request-body';
  fieldType?: FieldType;
}

export interface InteractiveField {
  name: string;
  type: FieldType;
  required: boolean;
  fieldType?: FieldType;
}

export interface CrossReferenceResult {
  matched: OpenApiField[];
  inOpenApiOnly: OpenApiField[];
  inTfOnly: string[];
}

export interface InteractiveScaffoldArgs {
  method: string;
  tfResource: string;
  featureKey: string;
  fields: InteractiveField[];
}

export interface OpenApiScaffoldArgs {
  apiPath: string;
  tfResource: string;
  matched: OpenApiField[];
  inOpenApiOnly: OpenApiField[];
  inTfOnly: string[];
}

// Maps internal field types to their TypeScript equivalents
export const TS_TYPE: Record<FieldType, string> = {
  string: 'string',
  boolean: 'boolean',
  number: 'number',
  object: 'Record<string, unknown>',
  array: 'unknown[]',
  heredoc: 'string',
};

// Provider-managed attributes that should never appear in the generated payload
export const PROVIDER_MANAGED = new Set(['id', 'accessor', 'namespace']);

// Fields whose values should be emitted as heredoc (formatEot) rather than quoted strings
export const HEREDOC_FIELD_NAMES = new Set(['policy', 'rules', 'content', 'config']);

// Maps an OpenAPI type string to an internal field type used for code generation
export const openApiTypeToFieldType = (openApiType: string, fieldName: string): FieldType => {
  if (openApiType === 'boolean') return 'boolean';
  if (openApiType === 'integer' || openApiType === 'number') return 'number';
  if (openApiType === 'object') return 'object';
  if (openApiType === 'array') return 'array';
  if (openApiType === 'string' && HEREDOC_FIELD_NAMES.has(fieldName)) return 'heredoc';
  return 'string';
};

// Generates a single resourceArgs line for a given field based on its type
export const resourceArgLine = ({
  name,
  type,
  fieldType,
}: {
  name: string;
  type?: FieldType;
  fieldType?: FieldType;
}): string => {
  const t = fieldType ?? type;
  if (t === 'string') return `      ${name}: \`"\${payload.${name}}"\`,`;
  if (t === 'heredoc') return `      ${name}: formatEot(payload.${name}),`;
  if (t === 'boolean' || t === 'number') return `      ${name}: payload.${name},`;
  return `      // TODO: ${name} — ${t} type, requires manual handling`;
};

// Cross-references OpenAPI fields against Terraform attributes, returning only
// fields present in both. Provider-managed fields (id, accessor, namespace) are
// always excluded from the Terraform side.
export const crossReference = (
  openApiFields: OpenApiField[],
  tfAttributes: Record<string, unknown>
): CrossReferenceResult => {
  const tfFieldNames = new Set(Object.keys(tfAttributes).filter((k) => !PROVIDER_MANAGED.has(k)));
  const openApiFieldNames = new Set(openApiFields.map((f) => f.name));

  const matched = openApiFields
    .filter((f) => tfFieldNames.has(f.name))
    .map((f) => ({ ...f, fieldType: openApiTypeToFieldType(f.openApiType, f.name) }));
  const inOpenApiOnly = openApiFields.filter((f) => !tfFieldNames.has(f.name));
  const inTfOnly = [...tfFieldNames].filter((n) => !openApiFieldNames.has(n));

  return { matched, inOpenApiOnly, inTfOnly };
};

// Derives a suggested registry feature key from the API path by stripping the /sys/ prefix,
// URL path parameters (e.g. /{name}), and any trailing slash.
// e.g. /sys/policies/acl/{name} → policies/acl
export const featureKeyHint = (apiPath: string): string =>
  apiPath
    .replace(/^\/sys\//, '')
    .replace(/\/{[^}]+}/g, '')
    .replace(/\/$/, '');

// Generates the full TypeScript scaffold string for interactive mode
export const generateInteractiveScaffold = ({
  method,
  tfResource,
  featureKey,
  fields,
}: InteractiveScaffoldArgs): string => {
  const interfaceName = `${classify(method)}Payload`;
  const mappingName = `${camelize(method)}Mapping`;
  const dasherized = dasherize(method);
  const docsResource = tfResource.replace('vault_', '');
  const hasHeredoc = fields.some((f) => f.type === 'heredoc');

  const interfaceLines = fields
    .map((f) => `  ${f.name}${f.required ? '' : '?'}: ${TS_TYPE[f.type] ?? 'unknown'};`)
    .join('\n');
  const argLines = fields.map(resourceArgLine).join('\n');
  const formattersImport = hasHeredoc
    ? `import { formatEot } from 'core/utils/code-generators/formatters';\n`
    : '';

  return `/**
 * Copyright IBM Corp. 2016, ${new Date().getFullYear()}
 * SPDX-License-Identifier: BUSL-1.1
 */

// ⚠️  AUTO-GENERATED SCAFFOLD — review and trim before registering
// Run: pnpm generate:terraform-mapping ${method}
//
// Steps before registering:
//   1. Remove fields not relevant to the Terraform resource
//   2. Resolve any object/array // TODO items
//   3. Verify field names match the Terraform provider docs:
//      https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/${docsResource}
//   4. Add the registry entry at the bottom of app/utils/terraform-registry.ts

import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
${formattersImport}
export interface ${interfaceName} {
${interfaceLines}
}

export const ${mappingName} = (payload: ${interfaceName}): string => {
  return terraformResourceTemplate({
    resource: '${tfResource}',
    localId: '<local-id>',
    resourceArgs: {
${argLines}
    },
  });
};

// Uncomment and move to app/utils/terraform-registry.ts when ready:
// import { ${mappingName} } from './terraform-mappings/${dasherized}-mapping';
// registry['${featureKey}'] = { multiBlock: false, mapping: ${mappingName} };
`;
};

// Generates the full TypeScript scaffold string for OpenAPI mode
export const generateOpenApiScaffold = ({
  apiPath,
  tfResource,
  matched,
  inOpenApiOnly,
  inTfOnly,
}: OpenApiScaffoldArgs): string => {
  const operationId = apiPath.replace(/\//g, '-').replace(/[{}]/g, '').replace(/^-/, '').replace(/-$/, '');

  const interfaceName = `${classify(camelize(operationId))}Payload`;
  const mappingName = `${camelize(operationId)}Mapping`;
  const dasherized = dasherize(operationId);
  const docsResource = tfResource.replace('vault_', '');
  const hasHeredoc = matched.some((f) => f.fieldType === 'heredoc');

  const interfaceLines = matched
    .map((f) => `  ${f.name}${f.required ? '' : '?'}: ${TS_TYPE[f.fieldType as FieldType] ?? 'unknown'};`)
    .join('\n');
  const argLines = matched.map(resourceArgLine).join('\n');
  const formattersImport = hasHeredoc
    ? `import { formatEot } from 'core/utils/code-generators/formatters';\n`
    : '';

  const openApiOnlyNote =
    inOpenApiOnly.length > 0
      ? `//\n// ⚠️  In Vault API but not in ${tfResource} (omitted):\n${inOpenApiOnly
          .map((f) => `//    - ${f.name}: ${f.description || f.openApiType}`)
          .join('\n')}\n`
      : '';
  const tfOnlyNote =
    inTfOnly.length > 0
      ? `//\n// ℹ️  In ${tfResource} but not in Vault API request body (omitted):\n${inTfOnly
          .map((n) => `//    - ${n}`)
          .join('\n')}\n`
      : '';

  const date = new Date();
  return `/**
 * Copyright IBM Corp. 2016, ${date.getFullYear()}
 * SPDX-License-Identifier: BUSL-1.1
 */

// ⚠️  AUTO-GENERATED from Vault OpenAPI + Terraform provider schema
// API path:           ${apiPath}
// Terraform resource: ${tfResource}
// Generated:          ${date.toISOString().slice(0, 10)}
//
// Steps before registering:
//   1. Review matched fields — remove any not relevant to your feature
//   2. Resolve any object/array // TODO items manually
//   3. Verify against Terraform provider docs:
//      https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/${docsResource}
//   4. Add the registry entry at the bottom of app/utils/terraform-registry.ts
${openApiOnlyNote}${tfOnlyNote}
import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
${formattersImport}
export interface ${interfaceName} {
${interfaceLines}
}

export const ${mappingName} = (payload: ${interfaceName}): string => {
  return terraformResourceTemplate({
    resource: '${tfResource}',
    localId: '<local-id>',
    resourceArgs: {
${argLines}
    },
  });
};

// Uncomment and move to app/utils/terraform-registry.ts when ready:
// import { ${mappingName} } from './terraform-mappings/${dasherized}-mapping';
// registry['${featureKeyHint(apiPath)}'] = { multiBlock: false, mapping: ${mappingName} };
`;
};
