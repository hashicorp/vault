/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// This is a policy can both mount and list and create oidc keys
export const adminPolicy = (mountPath: string) => {
  return `
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
    path "identity/oidc/key/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};

// user can mount the engine
// but user does not have access to oidc/key list or create
export const noOidcAdminPolicy = (mountPath: string) => {
  return `
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};
