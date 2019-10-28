import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  router: service(),
  init() {
    this._super(...arguments);
    this.router.on('routeWillChange', transition => {
      this.set('myTargetRouteName', transition.to.name);
    });
  },
  renderTemplate() {
    let targetName = this.myTargetRouteName;
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
