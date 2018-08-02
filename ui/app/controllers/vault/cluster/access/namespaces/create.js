import Ember from 'ember';

const { inject, Controller } = Ember;
export default Controller.extend({
  namespaceService: inject.service('namespace'),
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
