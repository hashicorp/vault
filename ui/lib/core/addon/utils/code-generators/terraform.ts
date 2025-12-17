/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { typeOf } from '@ember/utils';

// Yes, this seems silly but pretty formatting these snippets was a...journey.
// Hopefully these consts make it easier for whoever comes next.
const TWO_SPACES = `  `;
const FOUR_SPACES = `    `;

export interface TerraformTemplateArgs {
  resource?: string; // The Terraform resource type (e.g., "vault_auth_backend", "vault_mount")
  localId?: string; // The local identifier/label for this resource instance. Used in Terraform commands like `terraform import vault_auth_backend.example github` (<- is "example" here)
  resourceArgs?: Record<string, unknown>; //  Key/value pairs that build the terraform configuration, string interpolations should be wrapped in double `""` quotes, e.g. resourceArgs = { name: `"${name}"` };
}

/**
 * Generates a Terraform resource block for Vault providers.
 * @see https://registry.terraform.io/providers/hashicorp/vault/latest/docs
 * @example
 * ```
 * terraformTemplate({
 *   resource: 'vault_mount',
 *   localId: 'kvv2-example',
 *   resourceArgs: {
 *     path: '"my-kv-path"',
 *     type: '"kv-v2"',
 *     description: '"This is an example KV Version 2 secret engine mount"',
 *     options: {
 *       version: '"2"',
 *       type: '"kv-v2"'
 *     }
 *   }
 * })
 * ```
 *
 * Output:
 * ```
 * resource "vault_mount" "kvv2-example" {
 *   path = "my-kv-path"
 *   type = "kv-v2"
 *   description = "This is an example KV Version 2 secret engine mount"
 *   options = {
 *     version = "2"
 *     type = "kv-v2"
 *   }
 * }
 * ```
 */
export const terraformTemplate = ({
  resource = '<resource name>',
  localId = '<local identifier>',
  resourceArgs = {},
}: TerraformTemplateArgs = {}) => {
  const formattedContent = formatTerraformResourceArgs(resourceArgs);
  return `resource "${resource}" "${localId}" {
${formattedContent.join('\n')}
}`;
};

export const formatTerraformResourceArgs = (resourceArgs: Record<string, unknown> = {}) => {
  const formattedArgs = [];
  for (const [key, value] of Object.entries(resourceArgs)) {
    // Handle nested objects (like "options" above)
    if (typeOf(value) === 'object') {
      const formattedValue = formatNestedObject(key, value as Record<string, unknown>);
      formattedArgs.push(formattedValue);
      continue;
    }
    // Additional spaces before "key" are so arguments are indented over
    formattedArgs.push(formatKvPairs(TWO_SPACES, key, value));
  }
  return formattedArgs;
};

const formatNestedObject = (objKey: string, objValue: Record<string, unknown>): string => {
  const formatted = [];
  for (const [key, value] of Object.entries(objValue)) {
    // More spaces to indent even more because this block is nested
    formatted.push(formatKvPairs(FOUR_SPACES, key, value));
  }
  const formatObjValue = `{
${formatted.join('\n')}
  }`;
  return formatKvPairs(TWO_SPACES, objKey, formatObjValue);
};

const formatKvPairs = (indent: string, key: string, value: unknown) => `${indent}${key} = ${value}`;
