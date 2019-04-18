import Route from '@ember/routing/route';

export default Route.extend({
  renderTemplate() {
    let { targetName } = this.router.currentState.routerJs.activeTransition;
    let isCallback =
      targetName === 'vault.cluster.oidc-callback' || targetName === 'vault.cluster.oidc-callback-namespace';
    if (isCallback) {
      this.render('vault/cluster/oidc-callback', {
        into: 'application',
        outlet: 'main',
      });
    } else {
      this._super(...arguments);
    }
  },
});
