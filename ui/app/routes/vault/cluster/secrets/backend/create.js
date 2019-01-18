import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    let { secret, initialKey } = this.paramsFor(this.routeName);
    let qp = initialKey || secret;
    return this.transitionTo('vault.cluster.secrets.backend.create-root', {
      queryParams: { initialKey: qp },
    });
  },
});
