/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Maps Vault API path patterns to their corresponding Terraform resource type.
 *
 * This is the only piece of information the OpenAPI-driven generator cannot
 * derive automatically — the OpenAPI spec has no concept of Terraform resource
 * names. Add an entry here before running the generator for a new feature:
 *
 *   pnpm generate:terraform-mapping /sys/your/path/{name}
 */
export const TERRAFORM_RESOURCE_MAP: Record<string, string> = {
  '/sys/policies/acl/{name}': 'vault_policy',
};
