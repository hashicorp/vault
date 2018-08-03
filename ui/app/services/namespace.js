import Ember from 'ember';
import { task } from 'ember-concurrency';

const { Service, computed, inject } = Ember;
const ROOT_NAMESPACE = '';
export default Service.extend({
  store: inject.service(),
  //populated by the query param on the cluster route
  path: null,
  // list of namespaces available to the current user under the
  // current namespace
  accessibleNamespaces: null,

  inRootNamespace: computed.equal('path', ROOT_NAMESPACE),

  setNamespace(path) {
    this.set('path', path);
  },

  findNamespacesForUser: task(function*() {
    try {
      let ns = yield this.get('store').findAll('namespace');
      this.set('accessibleNamespaces', ns);
    } catch (e) {
      //do nothing here
    }
  }).drop(),
});
