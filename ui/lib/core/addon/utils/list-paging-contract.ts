/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Immutable record of Vault LIST endpoints that return a full key set with no server-side
 * page or limit parameters. UI callers use this to lock the session-cache fallback when
 * paginating in the browser.
 *
 * Mapping chain (each entry):
 *   Ember route → generated @hashicorp/vault-client-typescript client → HTTP handler → Go list API
 *
 * Entries:
 * - `acl`: GET /v1/sys/policies/acl with list=true → api/sys_policy.go ListPoliciesWithContext
 * - `pki-keys`: LIST pki/keys → builtin/logical/pki/path_fetch_keys.go pathListKeysHandler
 *   (returns keys and key_info in one response)
 *
 * Endpoint Group Policy (EGP) and Role Group Policy (RGP) lists use the same full-key LIST
 * shape but are intentionally out of scope for this contract.
 */

export type ListPagingMode = 'list-true-only';

export type ListPagingFallback = 'session-cache';

export type ListPagingContractEntry = {
  paging: ListPagingMode;
  fallback: ListPagingFallback;
};

export const LIST_PAGING_CONTRACT = {
  acl: {
    paging: 'list-true-only',
    fallback: 'session-cache',
  },
  'pki-keys': {
    paging: 'list-true-only',
    fallback: 'session-cache',
  },
} as const satisfies Record<string, ListPagingContractEntry>;

export type ListPagingContract = typeof LIST_PAGING_CONTRACT;
