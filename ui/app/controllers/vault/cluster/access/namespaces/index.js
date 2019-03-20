import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller from '@ember/controller';
export default Controller.extend({
  namespaceService: service('namespace'),
  accessibleNamespaces: alias('namespaceService.accessibleNamespaces'),
  currentNamespace: alias('namespaceService.path'),
  actions: {
    refreshNamespaceList() {
      // fetch new namespaces for the namespace picker
      this.get('namespaceService.findNamespacesForUser').perform();
    },
  },
});
