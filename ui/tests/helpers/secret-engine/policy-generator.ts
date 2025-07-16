/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// This policy can mount a secret engine
// and list and create oidc keys, relevant for setting identity_key_token for WIF
export const adminOidcCreateRead = (mountPath: string) => {
  return `
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
    path "identity/oidc/key/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
   path "${mountPath}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};

// This policy can mount the engine
// But does not have access to oidc/key list or read
export const adminOidcCreate = (mountPath: string) => {
  return `
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
    path "${mountPath}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
    path "identity/oidc/key/*" {
      capabilities = ["create", "update"]
    },
  `;
};
