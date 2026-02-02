/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import { formatEot } from './formatters';

export enum PolicyTypes {
  ACL = 'acl',
  RGP = 'rgp',
  EGP = 'egp',
}

export const ACL_CAPABILITIES = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'] as const;
export type AclCapability = (typeof ACL_CAPABILITIES)[number]; // 'create' | 'read' | 'update' | 'delete' | 'list' | 'patch' | 'sudo'

export class PolicyStanza {
  @tracked capabilities: Set<AclCapability> = new Set();
  @tracked path;

  constructor({ path = '' } = {}) {
    this.path = path;
  }

  get preview() {
    return aclTemplate(this.path, Array.from(this.capabilities));
  }
}

export const formatStanzas = (stanzas: PolicyStanza[]) => stanzas.map((s) => s.preview).join('\n');

export const policySnippetArgs = (policyName: string, policy: string) => {
  const formattedPolicy = formatEot(policy);
  const resourceArgs = { name: `"${policyName}"`, policy: formattedPolicy };
  return {
    terraform: { resource: 'vault_policy', resourceArgs },
    cli: { command: `policy write ${policyName}`, content: `- ${formattedPolicy}` },
  };
};

/**
 * Formats an ACL policy stanza in HCL
 * @param path - The Vault API path the policy applies to (e.g., "secret/data/*")
 * @param capabilities - Array of capabilities (e.g., '"read", "list"')
 * @returns A formatted HCL policy string
 */
export const aclTemplate = (path: string, capabilities: AclCapability[]) => {
  const formatted = formatCapabilities(capabilities);
  // Indentions below are intentional so policy renders prettily in code editor
  return `path "${path}" {
    capabilities = [${formatted}]
}`;
};

// returns a string with each capability wrapped in double quotes => ["create", "read"]
export const formatCapabilities = (capabilities: AclCapability[]) => {
  // Filter from ACL_CAPABILITIES to list capabilities in consistent order
  const allowed = ACL_CAPABILITIES.filter((p) => capabilities.includes(p));
  return allowed.length ? allowed.map((c) => `"${c}"`).join(', ') : '';
};

// Type Guards
export const isAclCapability = (value: string): value is AclCapability =>
  ACL_CAPABILITIES.includes(value as AclCapability);
