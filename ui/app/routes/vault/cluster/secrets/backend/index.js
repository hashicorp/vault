import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    return this.replaceWith('vault.cluster.secrets.backend.list-root');
  },
});
