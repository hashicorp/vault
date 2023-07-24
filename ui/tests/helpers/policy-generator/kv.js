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

export const secretPathCreateReadUpdate = (backend, secretPath) => {
  return `
    path "${backend}/metadata/*" {
      capabilities = ["list"]
    }
    path "${backend}/data/${secretPath}" {
      capabilities = ["create", "read", "update"]
    }
  `;
};

export const dataCRUDandMetadataUpdateDelete = (backend) => {
  return `
    path "${backend}/metadata/*" {
      capabilities = ["update","delete"]
    }
    path "${backend}/data/*" {
      capabilities = ["create", "read", "update", "delete"]
    }
  `;
};
