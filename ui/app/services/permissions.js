import Service, { inject as service } from '@ember/service';
import { task, waitForProperty } from 'ember-concurrency';

const API_PATHS = {
  secrets: { engine: 'cubbyhole/' },
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

export default Service.extend({
  exactPaths: null,
  globPaths: null,
  isRootToken: null,
  store: service(),
  auth: service(),

  getPaths: task(function*() {
    if (this.paths) {
      return;
    }
    let resp = yield this.get('store')
      .adapterFor('permissions')
      .query();
    this.setPaths(resp);
    return;
  }),

  setPaths(resp) {
    this.set('exactPaths', resp.data.exact_paths);
    this.set('globPaths', resp.data.glob_paths);
    this.set('isRootToken', resp.data.root);
  },

  checkAuthToken: task(function*() {
    yield waitForProperty(this.auth, 'currentTokenName', token => !!token);
    yield this.getPaths.perform();
  }),

  hasNavPermission(navItem, routeParams) {
    if (routeParams) {
      return this.hasPermission(API_PATHS[navItem][routeParams]);
    }
    return Object.values(API_PATHS[navItem]).some(path => this.hasPermission(path));
  },

  navPathParams(navItem) {
    const path = Object.values(API_PATHS[navItem]).find(path => this.hasPermission(path));
    // if it is policies or tools, split, otherwise look through a map
    if (['policies', 'tools'].includes(navItem)) {
      return path.split('/').lastObject;
    }

    return API_PATHS_TO_ROUTE_PARAMS[path];
  },

  hasPermission(pathName) {
    if (this.isRootToken || this.hasMatchingExactPath(pathName) || this.hasMatchingGlobPath(pathName)) {
      return true;
    }
    return false;
  },

  hasMatchingExactPath(pathName) {
    const exactPaths = this.get('exactPaths');
    if (exactPaths) {
      const prefix = Object.keys(exactPaths).find(path => path.startsWith(pathName));
      return prefix && this.isNotDenied(exactPaths[prefix]);
    }
    return false;
  },

  hasMatchingGlobPath(pathName) {
    const globPaths = this.get('globPaths');
    if (globPaths) {
      const matchingPath = Object.keys(globPaths).find(k => pathName.includes(k));
      return (matchingPath && this.isNotDenied(globPaths[matchingPath])) || globPaths.hasOwnProperty('');
    }
    return false;
  },

  isNotDenied(path) {
    return !path.capabilities.includes('deny');
  },
});
