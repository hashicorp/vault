import Service, { inject as service } from '@ember/service';
import { task, waitForProperty } from 'ember-concurrency';

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

  hasPermission(pathName) {
    if (this.isRootToken || this.hasMatchingExactPath(pathName) || this.hasMatchingGlobPath(pathName)) {
      return true;
    }
    return false;
  },

  hasMatchingExactPath(pathName) {
    const exactPaths = this.get('exactPaths');
    return exactPaths && exactPaths.hasOwnProperty(pathName) && this.isNotDenied(exactPaths[pathName]);
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
