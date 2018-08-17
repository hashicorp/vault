import Ember from 'ember';

const { computed, inject, Controller } = Ember;
export default Controller.extend({
  namespaceService: inject.service('namespace'),
  accessibleNamespaces: computed.alias('namespaceService.accessibleNamespaces'),
  currentNamespace: computed.alias('namespaceService.path'),
  actions: {
    refreshNamespaceList() {
      // fetch new namespaces for the namespace picker
      this.get('namespaceService.findNamespacesForUser').perform();
    },
  },
});
