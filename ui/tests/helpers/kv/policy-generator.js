/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ACL_CAPABILITIES, formatCapabilities } from 'core/utils/code-generators/policy';

const root = ACL_CAPABILITIES;

export const adminPolicy = (backend) => {
  return `
    path "${backend}/*" {
      capabilities = [${formatCapabilities(root)}]
    },
  `;
};

export const dataPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  // "delete" capability on this path can delete latest version
  return `
    path "${backend}/data/${secretPath}" {
      capabilities = [${formatCapabilities(capabilities)}]
    }
  `;
};

export const subkeysPolicy = ({ backend, secretPath = '*' }) => {
  return `
    path "${backend}/subkeys/${secretPath}" {
      capabilities = ["read"]
    }
  `;
};

export const dataNestedPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  return `
    path "${backend}/data/app/${secretPath}" {
      capabilities = [${formatCapabilities(capabilities)}]
    }
  `;
};

export const metadataPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  // "delete" capability on this path can destroy all versions
  return `
    path "${backend}/metadata/${secretPath}" {
        capabilities = [${formatCapabilities(capabilities)}]
    }
  `;
};

export const metadataNestedPolicy = ({ backend, secretPath = '*', capabilities = root }) => {
  return `
    path "${backend}/metadata/app/${secretPath}" {
        capabilities = [${formatCapabilities(capabilities)}]
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
  admin: (backend) => adminPolicy(backend) + subkeysPolicy({ backend }),
  dataReader: (backend) => dataPolicy({ backend, capabilities: ['read'] }),
  dataListReader: (backend) =>
    dataPolicy({ backend, capabilities: ['read', 'delete'] }) + metadataListPolicy(backend),
  metadataMaintainer: (backend) =>
    metadataPolicy({ backend, capabilities: ['create', 'read', 'update', 'list'] }) +
    deleteVersionsPolicy({ backend }) +
    undeleteVersionsPolicy({ backend }) +
    destroyVersionsPolicy({ backend }),
  secretNestedCreator: (backend) =>
    dataNestedPolicy({ backend, capabilities: ['create', 'update'] }) +
    metadataNestedPolicy({ backend, capabilities: ['list', 'delete'] }),
  secretCreator: (backend) =>
    dataPolicy({ backend, capabilities: ['create', 'update'] }) +
    metadataPolicy({ backend, capabilities: ['delete'] }),
  secretPatcher: (backend) =>
    // this persona should never have data "read"
    dataPolicy({ backend, capabilities: ['patch'] }) +
    metadataPolicy({ backend, capabilities: ['list', 'read'] }) +
    subkeysPolicy({ backend }) +
    // granting patcher persona delete/destroy capabilities because this matches policies from real customer use cases
    deleteVersionsPolicy({ backend }) +
    destroyVersionsPolicy({ backend }),
};
