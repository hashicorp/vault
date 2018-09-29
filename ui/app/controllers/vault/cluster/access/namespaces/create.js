import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
export default Controller.extend({
  namespaceService: service('namespace'),
  actions: {
    onSave({ saveType }) {
      if (saveType === 'save') {
        // fetch new namespaces for the namespace picker
        this.get('namespaceService.findNamespacesForUser').perform();
        return this.transitionToRoute('vault.cluster.access.namespaces.index');
      }
    },
  },
});
