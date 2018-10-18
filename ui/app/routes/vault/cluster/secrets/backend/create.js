import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    let { secret } = this.paramsFor(this.routeName);
    return this.transitionTo('vault.cluster.secrets.backend.create-root', {
      queryParams: { initialKey: secret },
    });
  },
});
