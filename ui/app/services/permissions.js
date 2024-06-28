/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { sanitizePath, sanitizeStart } from 'core/utils/sanitize-path';
import { task } from 'ember-concurrency';

export const PERMISSIONS_BANNER_STATES = {
  readFailed: 'read-failed',
  noAccess: 'no-ns-access',
};
const API_PATHS = {
  access: {
    methods: 'sys/auth',
    mfa: 'identity/mfa/method',
    oidc: 'identity/oidc/client',
    entities: 'identity/entity/id',
    groups: 'identity/group/id',
    leases: 'sys/leases/lookup',
    namespaces: 'sys/namespaces',
    'control-groups': 'sys/control-group/',
  },
  policies: {
    acl: 'sys/policies/acl',
    rgp: 'sys/policies/rgp',
    egp: 'sys/policies/egp',
  },
  tools: {
    wrap: 'sys/wrapping/wrap',
    lookup: 'sys/wrapping/lookup',
    unwrap: 'sys/wrapping/unwrap',
    rewrap: 'sys/wrapping/rewrap',
    random: 'sys/tools/random',
    hash: 'sys/tools/hash',
  },
  status: {
    replication: 'sys/replication',
    license: 'sys/license',
    seal: 'sys/seal',
    raft: 'sys/storage/raft/configuration',
  },
  clients: {
    activity: 'sys/internal/counters/activity',
    config: 'sys/internal/counters/config',
  },
  settings: {
    customMessages: 'sys/config/ui/custom-messages',
  },
  sync: {
    destinations: 'sys/sync/destinations',
    associations: 'sys/sync/associations',
    config: 'sys/sync/config',
    github: 'sys/sync/github-apps',
  },
};

const API_PATHS_TO_ROUTE_PARAMS = {
  'sys/auth': { route: 'vault.cluster.access.methods', models: [] },
  'identity/entity/id': { route: 'vault.cluster.access.identity', models: ['entities'] },
  'identity/group/id': { route: 'vault.cluster.access.identity', models: ['groups'] },
  'sys/leases/lookup': { route: 'vault.cluster.access.leases', models: [] },
  'sys/namespaces': { route: 'vault.cluster.access.namespaces', models: [] },
  'sys/control-group/': { route: 'vault.cluster.access.control-groups', models: [] },
  'identity/mfa/method': { route: 'vault.cluster.access.mfa', models: [] },
  'identity/oidc/client': { route: 'vault.cluster.access.oidc', models: [] },
};

/*
  The Permissions service is used to gate top navigation and sidebar items.
  It fetches a users' policy from the resultant-acl endpoint and stores their
  allowed exact and glob paths as state. It also has methods for checking whether
  a user has permission for a given path.
  The data from the resultant-acl endpoint has the following shape:
  {
    exact_paths: {
      [key: string]: {
        capabilities: string[];
      };
    };
    glob_paths: {
      [key: string]: {
        capabilities: string[];
      };
    };
    root: boolean;
    chroot_namespace?: string;
  };
  There are a couple nuances to be aware of about this response. When a
  chroot_namespace is set, all of the paths in the response will be prefixed
  with that namespace. Additionally, this endpoint is only added to the default
  policy in the user's root namespace, so we make the call to the user's root
  namespace (the namespace where the user's auth method is mounted) no matter
  what the current namespace is.
*/

export default class PermissionsService extends Service {
  @tracked exactPaths = null;
  @tracked globPaths = null;
  @tracked canViewAll = null;
  @tracked permissionsBanner = null;
  @tracked chrootNamespace = null;
  @service store;
  @service namespace;

  get fullCurrentNamespace() {
    const currentNs = this.namespace.path;
    return this.chrootNamespace
      ? `${sanitizePath(this.chrootNamespace)}/${sanitizePath(currentNs)}`
      : sanitizePath(currentNs);
  }

  @task *getPaths() {
    if (this.paths) {
      return;
    }

    try {
      const resp = yield this.store.adapterFor('permissions').query();
      this.setPaths(resp);
      return;
    } catch (err) {
      // If no policy can be found, default to showing all nav items.
      this.canViewAll = true;
      this.permissionsBanner = PERMISSIONS_BANNER_STATES.readFailed;
    }
  }

  /**
   * hasWildcardNsAccess checks if the user has a wildcard access to target namespace
   * via full glob path or any ancestors of the target namespace
   * @param {string} targetNs is the current/target namespace that we are checking access for
   * @param {object} globPaths key is path, value is object with capabilities
   * @returns {boolean} whether the user's policy includes wildcard access to NS
   */
  hasWildcardNsAccess(targetNs, globPaths = {}) {
    const nsParts = sanitizePath(targetNs).split('/');
    let matchKey = null;
    // For each section of the namespace, check if there is a matching wildcard path
    while (nsParts.length > 0) {
      // glob paths always end in a slash
      const test = `${nsParts.join('/')}/`;
      if (Object.keys(globPaths).includes(test)) {
        matchKey = test;
        break;
      }
      nsParts.pop();
    }
    // Finally, check if user has wildcard access to the root namespace
    // which is represented by an empty string
    if (!matchKey && Object.keys(globPaths).includes('')) {
      matchKey = '';
    }
    if (null === matchKey) {
      return false;
    }

    // if there is a match make sure the capabilities do not include deny
    return !this.isDenied(globPaths[matchKey]);
  }

  // This method is called to recalculate whether to show the permissionsBanner when the namespace changes
  calcNsAccess() {
    if (this.canViewAll) {
      this.permissionsBanner = null;
      return;
    }
    const namespace = this.fullCurrentNamespace;
    const allowed =
      // check if the user has wildcard access to the relative root namespace
      this.hasWildcardNsAccess(namespace, this.globPaths) ||
      // or if any of their glob paths start with the namespace
      Object.keys(this.globPaths).any((k) => k.startsWith(namespace) && !this.isDenied(this.globPaths[k])) ||
      // or if any of their exact paths start with the namespace
      Object.keys(this.exactPaths).any((k) => k.startsWith(namespace) && !this.isDenied(this.exactPaths[k]));
    this.permissionsBanner = allowed ? null : PERMISSIONS_BANNER_STATES.noAccess;
  }

  setPaths(resp) {
    this.exactPaths = resp.data.exact_paths;
    this.globPaths = resp.data.glob_paths;
    this.canViewAll = resp.data.root;
    this.chrootNamespace = resp.data.chroot_namespace;
    this.calcNsAccess();
  }

  reset() {
    this.exactPaths = null;
    this.globPaths = null;
    this.canViewAll = null;
    this.chrootNamespace = null;
    this.permissionsBanner = null;
  }

  hasNavPermission(navItem, routeParams, requireAll) {
    if (routeParams) {
      // check that the user has permission to access all (requireAll = true) or any of the routes when array is passed
      // useful for hiding nav headings when user does not have access to any of the links
      const params = Array.isArray(routeParams) ? routeParams : [routeParams];
      const evalMethod = !Array.isArray(routeParams) || requireAll ? 'every' : 'some';
      return params[evalMethod]((param) => {
        // viewing the entity and groups pages require the list capability, while the others require the default, which is anything other than deny
        const capability = param === 'entities' || param === 'groups' ? ['list'] : [null];
        return this.hasPermission(API_PATHS[navItem][param], capability);
      });
    }
    return Object.values(API_PATHS[navItem]).some((path) => this.hasPermission(path));
  }

  navPathParams(navItem) {
    const path = Object.values(API_PATHS[navItem]).find((path) => this.hasPermission(path));
    if (['policies', 'tools'].includes(navItem)) {
      return { models: [path.split('/').lastObject] };
    }

    return API_PATHS_TO_ROUTE_PARAMS[path];
  }

  pathNameWithNamespace(pathName) {
    const namespace = this.fullCurrentNamespace;
    if (namespace) {
      return `${sanitizePath(namespace)}/${sanitizeStart(pathName)}`;
    } else {
      return pathName;
    }
  }

  hasPermission(pathName, capabilities = [null]) {
    if (this.canViewAll) {
      return true;
    }
    const path = this.pathNameWithNamespace(pathName);

    return capabilities.every(
      (capability) =>
        this.hasMatchingExactPath(path, capability) || this.hasMatchingGlobPath(path, capability)
    );
  }

  hasMatchingExactPath(pathName, capability) {
    const exactPaths = this.exactPaths;
    if (exactPaths) {
      const prefix = Object.keys(exactPaths).find((path) => path.startsWith(pathName));
      const hasMatchingPath = prefix && !this.isDenied(exactPaths[prefix]);

      if (prefix && capability) {
        return this.hasCapability(exactPaths[prefix], capability) && hasMatchingPath;
      }

      return hasMatchingPath;
    }
    return false;
  }

  hasMatchingGlobPath(pathName, capability) {
    const globPaths = this.globPaths;
    if (globPaths) {
      const matchingPath = Object.keys(globPaths).find((k) => {
        return pathName.includes(k) || pathName.includes(k.replace(/\/$/, ''));
      });
      const hasMatchingPath =
        (matchingPath && !this.isDenied(globPaths[matchingPath])) ||
        Object.prototype.hasOwnProperty.call(globPaths, '');

      if (matchingPath && capability) {
        return this.hasCapability(globPaths[matchingPath], capability) && hasMatchingPath;
      }

      return hasMatchingPath;
    }
    return false;
  }

  hasCapability(path, capability) {
    return path.capabilities.includes(capability);
  }

  isDenied(path) {
    return path.capabilities.includes('deny');
  }
}
