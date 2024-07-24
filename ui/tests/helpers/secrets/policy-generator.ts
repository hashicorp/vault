/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// returns a string with each capability wrapped in double quotes => ["create", "read"]
// This is a policy can both mount a secret engine
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

// user can mount the engine
// but does not have access to oidc/key list or create
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
