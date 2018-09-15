import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel(transition) {
    const type = this.modelFor('vault.cluster.settings.configure-secret-backend').get('type');
    if (type === 'pki' && transition.targetName === this.routeName) {
      return this.transitionTo('vault.cluster.settings.configure-secret-backend.section', 'cert');
    }
  },
});
