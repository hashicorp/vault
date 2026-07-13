/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

// ⚠️  AUTO-GENERATED from Vault OpenAPI + Terraform provider schema
// API path:           /sys/policies/acl/{name}
// Terraform resource: vault_policy
// Generated:          2026-07-08
//
// Steps before registering:
//   1. Review matched fields — remove any not relevant to your feature
//   2. Resolve any object/array // TODO items manually
//   3. Verify against Terraform provider docs:
//      https://registry.terraform.io/providers/hashicorp/vault/latest/docs/resources/policy
//   4. Add the registry entry at the bottom of app/utils/terraform-registry.ts
//
// ℹ️  In vault_policy but not in Vault API request body (omitted):
//    - allow_overwrite

import { terraformResourceTemplate } from 'core/utils/code-generators/terraform';
import { formatEot } from 'core/utils/code-generators/formatters';

export interface SysPoliciesAclNamePayload {
  name: string;
  policy: string;
}

export const sysPoliciesAclNameMapping = (payload: SysPoliciesAclNamePayload): string => {
  return terraformResourceTemplate({
    resource: 'vault_policy',
    localId: '<local identifier>',
    resourceArgs: {
      name: `"${payload.name || '<policy name>'}"`,
      policy: formatEot(payload.policy),
    },
  });
};
