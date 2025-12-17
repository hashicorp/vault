/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

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
