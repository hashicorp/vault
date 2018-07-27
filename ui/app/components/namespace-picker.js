import Ember from 'ember';
import { task } from 'ember-concurrency';

const { Component, computed, inject } = Ember;

export default Component.extend({
  namespaceService: inject.service('namespace'),
  router: inject.service(),
  store: inject.service(),
  //passed from the queryParam
  namespace: null,

  init() {
    this._super(...arguments);
    this.get('findNamespacesForUser').perform();
  },

  namespaceDisplay: computed('namespace', 'accessibleNamespaces', function() {
    let namespace = this.get('namespace');
    if (namespace === '') {
      return 'Default';
    } else {
      let parts = namespace.split('/');
      return parts[parts.length - 1];
    }
  }),

  findNamespacesForUser: task(function*() {
    try {
      let ns = yield this.get('store').findAll('namespace');
      this.set('accessibleNamespaces', ns);
    } catch (e) {
      //do nothing here
    }
  }).drop(),

  didRecieveAttrs() {
    this._super(...arguments);
    this.get('namespaceService').setNamespace(this.get('namespace'));
  },
});
