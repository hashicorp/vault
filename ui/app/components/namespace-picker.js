import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
  namespaceService: inject.service('namespace'),
  router: inject.service(),
  auth: inject.service(),
  //passed from the queryParam
  namespace: null,
  // internal tracking of namespace
  oldNamespace: null,

  accessibleNamespaces: computed.alias('namespaceService.accessibleNamespaces'),
  canAccessMultipleNamespaces: computed.gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.currentToken', 'namespace', function() {
    return this.get('auth.currentToken.rootNamespace') === this.get('namespace');
  }),

  namespaceDisplay: computed('namespace', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.get('namespace');
    if (namespace === '') {
      return 'Default';
    } else {
      let parts = namespace.split('/');
      return parts[parts.length - 1];
    }
  }),

  init() {
    this._super(...arguments);
    this.get('namespaceService').setNamespace(this.get('namespace'));
  },
  didReceiveAttrs() {
    let ns = this.get('namespace');
    let oldNS = this.get('oldNamespace');
    this._super(...arguments);
    if (oldNS !== null && oldNS !== ns) {
      this.get('namespaceService').setNamespace(ns);
    }
    this.set('oldNamespace', ns);
  },
});
