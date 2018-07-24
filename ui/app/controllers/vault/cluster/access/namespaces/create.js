import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    onSave({ saveType }) {
      if (saveType === 'save') {
        return this.transitionToRoute('vault.cluster.access.namespaces.index');
      }
    },
  },
});
