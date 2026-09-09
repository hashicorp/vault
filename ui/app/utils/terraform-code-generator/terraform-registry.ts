/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { sysPoliciesAclNameMapping } from '../terraform-mappings/sys-policies-acl-name-mapping';

// Types
export interface TerraformBlock {
  type: 'resource' | 'variable';
  content: string;
}

export interface SingleBlockDescriptor<T extends Record<string, unknown>> {
  multiBlock: false;
  mapping: (payload: T) => string;
}

export interface MultiBlockDescriptor<T extends Record<string, unknown>> {
  multiBlock: true;
  mapping: (payload: T) => TerraformBlock[];
}

export type TerraformDescriptor<T extends Record<string, unknown>> =
  | SingleBlockDescriptor<T>
  | MultiBlockDescriptor<T>;

// Registry

const registry: Record<string, TerraformDescriptor<Record<string, unknown>>> = {};

/** Look up a descriptor by feature key. Returns undefined if not registered. */
export const getTerraformDescriptor = (featureKey: string) => registry[featureKey];

// Renderers

/**
 * Renders a TerraformBlock array to a single HCL string.
 * Variables are always emitted before resources; within each group,
 * order matches the array (mapping functions control dependency order).
 */
export const renderTerraformBlocks = (blocks: TerraformBlock[]): string => {
  return blocks.map((b) => b.content).join('\n\n');
};

/**
 * Produces a Terraform cross-resource reference string.
 * Only meaningful in multiBlock: true mappings where localId is derived
 * from the payload rather than using the '<local-id>' placeholder.
 *
 * @example
 * ref('vault_mount', 'kv', 'path') // → "vault_mount.kv.path"
 * ref('vault_mount', 'kv')         // → "vault_mount.kv"
 */
export const ref = (resourceType: string, localId: string, attribute?: string) => {
  const base = `${resourceType}.${localId}`;
  return attribute ? `${base}.${attribute}` : base;
};

/**  Registry */

registry['policies/acl'] = {
  multiBlock: false,
  mapping: sysPoliciesAclNameMapping as unknown as (payload: Record<string, unknown>) => string,
};
