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
  // "delete" capability on this path can delete latest version
  return `
    path "${backend}/data/${secretPath}" {
      capabilities = [${format(capabilities)}]
    }
  `;
};

export const metadataPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  // "delete" capability on this path can destroy all versions
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

export const deleteVersionsPolicy = ({ backend, secretPath = '*' }) => {
  return `
    path "${backend}/delete/${secretPath}" {
        capabilities = ["update"]
    }
  `;
};
export const undeleteVersionsPolicy = ({ backend, secretPath = '*' }) => {
  return `
    path "${backend}/undelete/${secretPath}" {
        capabilities = ["update"]
    }
  `;
};
export const destroyVersionsPolicy = ({ backend, secretPath = '*' }) => {
  return `
    path "${backend}/destroy/${secretPath}" {
        capabilities = ["update"]
    }
  `;
};

// Personas for reuse in workflow tests
export const personas = {
  admin: (backend) => adminPolicy(backend),
  dataReader: (backend) => dataPolicy({ backend, capabilities: ['read'] }),
  dataListReader: (backend) =>
    dataPolicy({ backend, capabilities: ['read', 'delete'] }) + metadataListPolicy(backend),
  metadataMaintainer: (backend) =>
    metadataListPolicy(backend) +
    metadataPolicy({ backend, capabilities: ['create', 'read', 'update', 'list'] }) +
    deleteVersionsPolicy({ backend }) +
    undeleteVersionsPolicy({ backend }) +
    destroyVersionsPolicy({ backend }),
  secretCreator: (backend) =>
    dataPolicy({ backend, capabilities: ['create', 'update'] }) +
    metadataPolicy({ backend, capabilities: ['delete'] }),
};
