import Service, { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

const API_PATHS = {
  access: {
    methods: 'sys/auth',
    entities: 'identity/entities',
    groups: 'identity/groups',
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
  },
};

const API_PATHS_TO_ROUTE_PARAMS = {
  'sys/auth': ['vault.cluster.access.methods'],
  'identity/entities': ['vault.cluster.access.identity', 'entities'],
  'identity/groups': ['vault.cluster.access.identity', 'groups'],
  'sys/leases/lookup': ['vault.cluster.access.leases'],
  'sys/namespaces': ['vault.cluster.access.namespaces'],
  'sys/control-group/': ['vault.cluster.access.control-groups'],
};

/*
  The Permissions service is used to gate top navigation and sidebar items. It fetches
  a users' policy from the resultant-acl endpoint and stores their allowed exact and glob
  paths as state. It also has methods for checking whether a user has permission for a given
  path.
*/
export default Service.extend({
  exactPaths: null,
  globPaths: null,
  canViewAll: null,
  store: service(),
  auth: service(),
  namespace: service(),

  getPaths: task(function*() {
    if (this.paths) {
      return;
    }

    try {
      let resp = yield this.get('store')
        .adapterFor('permissions')
        .query();
      this.setPaths(resp);
      return;
    } catch (err) {
      // If no policy can be found, default to showing all nav items.
      this.set('canViewAll', true);
    }
  }),

  setPaths(resp) {
    this.set('exactPaths', resp.data.exact_paths);
    this.set('globPaths', resp.data.glob_paths);
    this.set('canViewAll', resp.data.root);
  },

  reset() {
    this.set('exactPaths', null);
    this.set('globPaths', null);
    this.set('canViewAll', null);
  },

  hasNavPermission(navItem, routeParams) {
    if (routeParams) {
      return this.hasPermission(API_PATHS[navItem][routeParams]);
    }
    return Object.values(API_PATHS[navItem]).some(path => this.hasPermission(path));
  },

  navPathParams(navItem) {
    const path = Object.values(API_PATHS[navItem]).find(path => this.hasPermission(path));
    if (['policies', 'tools'].includes(navItem)) {
      return path.split('/').lastObject;
    }

    return API_PATHS_TO_ROUTE_PARAMS[path];
  },

  pathNameWithNamespace(pathName) {
    const namespace = this.get('namespace').path;
    if (namespace) {
      return `${namespace}/${pathName}`;
    } else {
      return pathName;
    }
  },

  hasPermission(pathName, capabilities = [null]) {
    const path = this.pathNameWithNamespace(pathName);

    if (this.canViewAll) {
      return true;
    }

    return capabilities.every(
      capability => this.hasMatchingExactPath(path, capability) || this.hasMatchingGlobPath(path, capability)
    );
  },

  hasMatchingExactPath(pathName, capability) {
    const exactPaths = this.get('exactPaths');
    if (exactPaths) {
      const prefix = Object.keys(exactPaths).find(path => path.startsWith(pathName));
      const hasMatchingPath = prefix && !this.isDenied(exactPaths[prefix]);

      if (prefix && capability) {
        return this.hasCapability(exactPaths[prefix], capability) && hasMatchingPath;
      }

      return hasMatchingPath;
    }
    return false;
  },

  hasMatchingGlobPath(pathName, capability) {
    const globPaths = this.get('globPaths');
    if (globPaths) {
      const matchingPath = Object.keys(globPaths).find(k => {
        return pathName.includes(k) || pathName.includes(k.replace(/\/$/, ''));
      });
      const hasMatchingPath =
        (matchingPath && !this.isDenied(globPaths[matchingPath])) || globPaths.hasOwnProperty('');

      if (matchingPath && capability) {
        return this.hasCapability(globPaths[matchingPath], capability) && hasMatchingPath;
      }

      return hasMatchingPath;
    }
    return false;
  },

  hasCapability(path, capability) {
    return path.capabilities.includes(capability);
  },

  isDenied(path) {
    return path.capabilities.includes('deny');
  },
});
