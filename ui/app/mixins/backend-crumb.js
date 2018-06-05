import Ember from 'ember';

export default Ember.Mixin.create({
  backendCrumb: Ember.computed('backend', function() {
    const backend = this.get('backend');

    if (backend === undefined) {
      throw new Error('backend-crumb mixin requires backend to be set');
    }

    return {
      label: backend,
      text: backend,
      path: 'vault.cluster.secrets.backend.list-root',
      model: backend,
    };
  }),
});
