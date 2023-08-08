/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const root = ['create', 'read', 'update', 'delete', 'list'];

// returns a string with each capability wrapped in double quotes => ["create", "read"]
const format = (array) => array.map((c) => `"${c}"`).join(', ');

export const adminPolicy = (backend) => {
  return `
    path "${backend}/*" {
      capabilities = [${format(root)}]
    },
  `;
};

export const dataPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  return `
    path "${backend}/data/${secretPath}" {
      capabilities = [${format(capabilities)}]
    }
  `;
};

export const metadataPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  return `
    path "${backend}/metadata/${secretPath}" {
        capabilities = [${format(capabilities)}]
    }
  `;
};

export const metadataListPolicy = (backend) => {
  return `
    path "${backend}/metadata" {
        capabilities = ["list"]
    }
  `;
};
