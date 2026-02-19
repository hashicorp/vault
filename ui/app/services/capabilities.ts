/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { sanitizePath, sanitizeStart } from 'core/utils/sanitize-path';
import { PATH_MAP, SUDO_PATHS, SUDO_PATH_PREFIXES } from 'vault/utils/constants/capabilities';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type NamespaceService from 'vault/services/namespace';
import type { Capabilities, CapabilitiesMap, CapabilitiesData, CapabilityTypes } from 'vault/app-types';

type CapabilityFetchOptions = {
  routeForCache?: string;
};

export default class CapabilitiesService extends Service {
  @service declare readonly api: ApiService;
  @service declare readonly namespace: NamespaceService;

  /*
   * API path caching for <CodeGeneratorPolicyFlyout />
   * Cache API paths when capabilities are fetched during route model hooks
   * so the flyout can prefill with relevant policy paths.
   */
  @tracked routePathCache = new Map<string, Set<string>>();

  // Cache API paths requested for a particular route
  cacheRoutePaths(route: string, apiPaths: string[]) {
    this.routePathCache.set(route, new Set(apiPaths));
  }

  // Lookup the paths for a route. Returns exact match first then falls back to longest matching route prefix.
  lookupRoutePaths(route: string) {
    const exactMatch = this.routePathCache.get(route);
    if (exactMatch) {
      return exactMatch;
    }

    const routes = Array.from(this.routePathCache.keys());
    const matchingRoutes = routes.filter((r) => route.includes(r));
    if (matchingRoutes.length) {
      // Return longest matching cached route which is the most specific parent
      const bestMatch = matchingRoutes.reduce((a, b) => (a.length > b.length ? a : b), '');
      return this.routePathCache.get(bestMatch);
    }

    return null;
  }
  /* End logic for <CodeGeneratorPolicyFlyout /> API path caching */

  /*
  Add API paths to the PATH_MAP constant using a friendly key, e.g. 'syncDestinations'.
  Use the apiPath tagged template literal to build the path with dynamic segments
  Each path should include placeholders for dynamic values -> apiPath`sys/sync/destinations/${'type'}/${'name'}`
  Provide the key and an object whose keys match the dynamic segment names
  The values from the object will be inserted into the placeholders and the fully-resolved path string will be returned.
  */
  pathFor<T>(key: keyof typeof PATH_MAP, params?: T) {
    const path = PATH_MAP[key];
    if (!path) {
      throw new Error(`Path not found for key: ${key}`);
    }
    return path(params || {});
  }

  pathsForList(paths: (keyof typeof PATH_MAP)[], params: object) {
    return paths.map((path) => this.pathFor(path, params));
  }

  /*
  Users don't always have access to the capabilities-self endpoint in the current namespace.
  This can happen when logging in to a namespace and then navigating to a child namespace.
  The "relativeNamespace" refers to the namespace the user is currently in and attempting to access capabilities for.
  Prepending "relativeNamespace" to the path while making the request in the "userRootNamespace" context (meaning, "userRootNamespace" is the namespace header)
  ensures we are querying capabilities-self where the user is most likely to have their policy/permissions.
  */
  relativeNamespacePath(path: string) {
    const { relativeNamespace } = this.namespace;
    // sanitizeStart ensures original path doesn't have leading slash
    return relativeNamespace ? `${relativeNamespace}/${sanitizeStart(path)}` : path;
  }

  // map capabilities to friendly names like canRead, canUpdate, etc.
  mapCapabilities(paths: string[], capabilitiesData: CapabilitiesData) {
    // request may not return capabilities for all provided paths
    // loop provided paths and map capabilities, defaulting to true for missing paths
    return paths.reduce((mappedCapabilities: CapabilitiesMap, path) => {
      // key in capabilitiesData includes relativeNamespace if applicable
      const key = this.relativeNamespacePath(path);
      const capabilities = capabilitiesData[key];

      const getCapability = (capability: CapabilityTypes) => {
        if (!(key in capabilitiesData)) {
          return true;
        }
        if (!capabilities?.length || capabilities.includes('deny')) {
          return false;
        }
        if (capabilities.includes('root')) {
          return true;
        }
        // if the path is sudo protected, they'll need sudo + the appropriate capability
        if (SUDO_PATHS.includes(key) || SUDO_PATH_PREFIXES.find((item) => key.startsWith(item))) {
          return capabilities.includes('sudo') && capabilities.includes(capability);
        }
        return capabilities.includes(capability);
      };
      // remove relativeNamespace from the path that was added for the request
      mappedCapabilities[path] = {
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

  async fetch(paths: string[], fetchOptions: CapabilityFetchOptions = {}): Promise<CapabilitiesMap> {
    // Cache API paths if route name provided
    if (fetchOptions?.routeForCache) {
      this.cacheRoutePaths(fetchOptions.routeForCache, paths);
    }

    const payload = { paths: paths.map((path) => this.relativeNamespacePath(path)) };

    try {
      const { data } = await this.api.sys.queryTokenSelfCapabilities(
        payload,
        this.api.buildHeaders({ namespace: sanitizePath(this.namespace.userRootNamespace) })
      );
      return this.mapCapabilities(paths, data as CapabilitiesData);
    } catch (e) {
      // default to true if there is a problem fetching the model
      // we can rely on the API to gate as a fallback
      return paths.reduce((obj: CapabilitiesMap, path: string) => {
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

  // convenience method for fetching capabilities for a singular path without needing to use pathFor
  // ex: capabilities.for('syncDestinations', { type: 'github', name: 'org-sync' });
  async for<T>(key: keyof typeof PATH_MAP, params?: T, fetchOptions: CapabilityFetchOptions = {}) {
    const path = this.pathFor(key, params);
    return this.fetchPathCapabilities(path, fetchOptions);
  }

  /*
  this method returns all of the capabilities for a singular path 
  */
  async fetchPathCapabilities(path: string, fetchOptions: CapabilityFetchOptions = {}) {
    const capabilities = await this.fetch([path], fetchOptions);
    return capabilities[path] as Capabilities;
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
