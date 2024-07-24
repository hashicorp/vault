/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// This is a policy that can mount a secret engine and list and create oidc keys
// Relevant for setting identity_key_token for WIF
export const adminPolicy = () => {
  return `
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
    path "identity/oidc/key/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};

// User can mount the engine
// But does not have access to oidc/key list or create
export const noOidcAdminPolicy = () => {
  return `
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};
