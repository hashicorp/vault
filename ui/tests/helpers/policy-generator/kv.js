/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const adminPolicy = (backend) => {
  return `
    path "${backend}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};

// DATA POLICIES
export const dataSecretPathCreateReadUpdate = (backend, secretPath) => {
  return `
    path "${backend}/data/${secretPath}" {
      capabilities = ["create", "read", "update"]
    }
  `;
};

// METADATA POLICIES
export const metadataListOnly = (backend) => {
  return `
    path "${backend}/metadata/*" {
      capabilities = ["list"]
    }
  `;
};
