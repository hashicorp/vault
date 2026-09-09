/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { formatArgsFromPayload } from './formatters';

// The CliTemplateArgs intentionally does NOT include "namespace" because the location for CLI commands is not consistent.
// Additionally, namespace can be specified via the environment variable `VAULT_NAMESPACE` and passing a flag is unnecessary.

export interface CliTemplateArgs {
  command?: string; // The CLI command that comes after "vault" (e.g., "policy write my-policy")
  content?: string; // Vault CLI commands are not consistent so this is essentially whatever needs to come after the command, usually param flags e.g. `-path="some-path"` or `-ttl='24h'`
}

export const cliTemplate = ({ command = '', content = '' }: CliTemplateArgs = {}) => {
  // Show placeholders when no args are provided
  if (!command && !content) {
    return 'vault <command> [args]';
  }

  const segments = ['vault', command, content].filter(Boolean);
  return segments.join(' ');
};

// generate a CLI command with args from a generic object or form payload
// the payload object will be converted to CLI flags in the format of `-key=value` and appended to the command
export const generateCliCommand = (command = '', payload: Record<string, unknown> = {}) => {
  const filteredArgs = formatArgsFromPayload(payload);
  const content = Object.entries(filteredArgs)
    .map(([key, value]) => {
      // For boolean flags, include the flag without a value if true, and omit if false
      if (typeof value === 'boolean') {
        return value ? `-${key}` : '';
      }
      // For array values, join them with commas (e.g., `-key=value1,value2`)
      if (Array.isArray(value)) {
        return `-${key}=${value.join(',')}`;
      }
      // for nested objects, repeat the flag for each key-value pair (e.g., `-options="version=2" -options="type=kv-v2"`)
      if (typeof value === 'object' && value !== null) {
        return Object.entries(value)
          .map(([nestedKey, nestedValue]) => `-${key}="${nestedKey}=${nestedValue}"`)
          .join(' ');
      }
      // For other types, include the flag with its value
      return `-${key}=${value}`;
    })
    .filter(Boolean) // Remove any empty strings resulting from false boolean flags
    .join(' ');

  return cliTemplate({ command, content });
};

export const generateCliWriteCommand = (path: string, payload: Record<string, unknown> = {}) =>
  generateCliCommand(`write ${path}`, payload);
