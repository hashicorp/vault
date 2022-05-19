import Route from '@ember/routing/route';

export default class MfaConfigureRoute extends Route {
  beforeModel() {
    return this.store
      .query('mfa-method', {})
      .then(() => {
        // if response then they should transition to the methods page instead of staying on the configure page.
        this.transitionTo('vault.cluster.access.mfa.methods.index');
      })
      .catch(() => {
        // stay on the landing page
      });
  }
}
