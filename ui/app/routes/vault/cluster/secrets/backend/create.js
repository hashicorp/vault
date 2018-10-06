import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel(transition) {
    let { secret } = this.paramsFor(this.routeName);
    return this.transitionTo('vault.cluster.secrets.backend.create-root', {
      queryParams: { initialKey: secret },
    });
  },
});
