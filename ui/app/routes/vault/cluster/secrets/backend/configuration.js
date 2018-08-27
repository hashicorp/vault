import Ember from 'ember';

export default Ember.Route.extend({
  wizard: Ember.inject.service(),
  model() {
    let backend = this.modelFor('vault.cluster.secrets.backend');
    if (this.get('wizard.featureState') === 'list') {
      this.get('wizard').transitionFeatureMachine(
        this.get('wizard.featureState'),
        'CONTINUE',
        backend.get('type')
      );
    }
    return backend;
  },
});
