/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { sanitizePath, sanitizeStart } from 'core/utils/sanitize-path';

import type ApiService from 'vault/services/api';
import type NamespaceService from 'vault/services/namespace';

interface Capabilities {
  canCreate: boolean;
  canDelete: boolean;
  canList: boolean;
  canPatch: boolean;
  canRead: boolean;
  canSudo: boolean;
  canUpdate: boolean;
}

interface MultipleCapabilities {
  [key: string]: Capabilities;
}

type CapabilityTypes = 'root' | 'sudo' | 'deny' | 'create' | 'read' | 'update' | 'delete' | 'list' | 'patch';
interface CapabilitiesData {
  [key: string]: CapabilityTypes[];
}

export default class CapabilitiesService extends Service {
  @service declare readonly api: ApiService;
  @service declare readonly namespace: NamespaceService;

  SUDO_PATHS = [
    'sys/seal',
    'sys/replication/performance/primary/secondary-token',
    'sys/replication/dr/primary/secondary-token',
    'sys/replication/reindex',
    'sys/leases/lookup/',
  ];
  SUDO_PATH_PREFIXES = ['sys/leases/revoke-prefix', 'sys/leases/revoke-force'];

  /*
  Users don't always have access to the capabilities-self endpoint in the current namespace.
  This can happen when logging in to a namespace and then navigating to a child namespace.
  The "relativeNamespace" refers to the namespace the user is currently in and attempting to access capabilities for.
  Prepending "relativeNamespace" to the path while making the request to the "userRootNamespace"
  ensures we are querying capabilities-self where the user is most likely to have their policy/permissions.
  */
  relativeNamespacePaths(paths: string[]) {
    const { relativeNamespace } = this.namespace;
    // sanitizeStart ensures original path doesn't have leading slash
    return paths.map((path) => (relativeNamespace ? `${relativeNamespace}/${sanitizeStart(path)}` : path));
  }

  // map capabilities to friendly names like canRead, canUpdate, etc.
  mapCapabilities(relativeNamespacePaths: string[], capabilitiesData: CapabilitiesData) {
    const { SUDO_PATHS, SUDO_PATH_PREFIXES } = this;
    const { relativeNamespace } = this.namespace;
    // request may not return capabilities for all provided paths
    // loop provided paths and map capabilities, defaulting to true for missing paths
    return relativeNamespacePaths.reduce((mappedCapabilities: MultipleCapabilities, path) => {
      const capabilities = capabilitiesData[path];

      const getCapability = (capability: CapabilityTypes) => {
        if (!(path in capabilitiesData)) {
          return true;
        }
        if (!capabilities?.length || capabilities.includes('deny')) {
          return false;
        }
        if (capabilities.includes('root')) {
          return true;
        }
        // if the path is sudo protected, they'll need sudo + the appropriate capability
        if (SUDO_PATHS.includes(path) || SUDO_PATH_PREFIXES.find((item) => path.startsWith(item))) {
          return capabilities.includes('sudo') && capabilities.includes(capability);
        }
        return capabilities.includes(capability);
      };
      // remove relativeNamespace from the path that was added for the request
      const key = path.replace(relativeNamespace, '');
      mappedCapabilities[key] = {
        canCreate: getCapability('create'),
        canDelete: getCapability('delete'),
        canList: getCapability('list'),
        canPatch: getCapability('patch'),
        canRead: getCapability('read'),
        canSudo: getCapability('sudo'),
        canUpdate: getCapability('update'),
      };
      return mappedCapabilities;
    }, {});
  }

  async fetch(paths: string[]): Promise<MultipleCapabilities> {
    const payload = {
      paths: this.relativeNamespacePaths(paths),
      namespace: sanitizePath(this.namespace.userRootNamespace),
    };

    try {
      const { data } = await this.api.sys.queryTokenSelfCapabilities(payload);
      return this.mapCapabilities(payload.paths, data as CapabilitiesData);
    } catch (e) {
      // default to true if there is a problem fetching the model
      // we can rely on the API to gate as a fallback
      return paths.reduce((obj: MultipleCapabilities, path: string) => {
        obj[path] = {
          canCreate: true,
          canDelete: true,
          canList: true,
          canPatch: true,
          canRead: true,
          canSudo: true,
          canUpdate: true,
        };
        return obj;
      }, {});
    }
  }

  /*
  this method returns all of the capabilities for a singular path 
  */
  async fetchPathCapabilities(path: string) {
    const capabilities = await this.fetch([path]);
    return capabilities[path];
  }

  /* 
  internal method for specific capability checks below
  checks the capability model for the passed capability, ie "canRead"
  */
  async _fetchSpecificCapability(path: string, capability: keyof Capabilities) {
    const capabilities = await this.fetchPathCapabilities(path);
    return capabilities ? capabilities[capability] : true;
  }

  canRead(path: string) {
    return this._fetchSpecificCapability(path, 'canRead');
  }

  canUpdate(path: string) {
    return this._fetchSpecificCapability(path, 'canUpdate');
  }

  canPatch(path: string) {
    return this._fetchSpecificCapability(path, 'canPatch');
  }
}
