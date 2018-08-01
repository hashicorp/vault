import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
  namespaceService: inject.service('namespace'),
  router: inject.service(),
  auth: inject.service(),

  //passed from the queryParam
  namespace: null,

  namespacePath: computed('namespace', function() {
    let namespace = this.get('namespace');
    // the queryParam default is 'default',
    // but the default ns path for vault is ''
    if (namespace === 'default') {
      return '';
    }
    return namespace;
  }),

  init() {
    this._super(...arguments);
    this.get('namespaceService').setNamespace(this.get('namespacePath'));
  },
  didReceiveAttrs() {
    let ns = this.get('namespacePath');
    let oldNS = this.get('oldNamespace');
    this._super(...arguments);
    if (oldNS !== null && oldNS !== ns) {
      this.get('namespaceService').setNamespace(ns);
    }
    this.set('oldNamespace', ns);
  },

  // internal tracking of namespace
  oldNamespace: null,

  accessibleNamespaces: computed.alias('namespaceService.accessibleNamespaces'),
  canAccessMultipleNamespaces: computed.gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.authData.userRootNamespace', 'namespacePath', function() {
    return this.get('auth.authData.userRootNamespace') === this.get('namespacePath');
  }),

  namespaceDisplay: computed('namespacePath', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.get('namespacePath');
    if (namespace === '') {
      return '';
    } else {
      let parts = namespace.split('/');
      return parts[parts.length - 1];
    }
  }),
});
