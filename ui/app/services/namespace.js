import Ember from 'ember';
import { task } from 'ember-concurrency';

const { Service, computed, inject } = Ember;
const ROOT_NAMESPACE = '';
export default Service.extend({
  store: inject.service(),
  auth: inject.service(),
  userRootNamespace: computed.alias('auth.authData.userRootNamespace'),
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
    // uses the adapter and the raw response here since
    // models get wiped when switching namespaces and we
    // want to keep track of these separately
    let store = this.get('store');
    let adapter = store.adapterFor('namespace');
    try {
      let ns = yield adapter.findAll(store, 'namespace', null, {
        adapterOptions: {
          forUser: true,
          namespace: this.get('userRootNamespace'),
        },
      });
      this.set('accessibleNamespaces', ns.data.keys.map(n => n.replace(/\/$/, '')));
    } catch (e) {
      //do nothing here
    }
  }).drop(),
});
