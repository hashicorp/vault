import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
  namespaceService: inject.service('namespace'),
  auth: inject.service(),

  init() {
    this._super(...arguments);
    this.get('namespaceService.findNamespacesForUser').perform();
  },

  namespacePath: computed.alias('namespaceService.path'),

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
