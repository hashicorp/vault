import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { AUTH, CLUSTER } from 'vault/lib/route-paths';

export default class VaultClusterRedirectRoute extends Route {
  @service auth;
  @service router;

  beforeModel({ to: { queryParams } }) {
    let transition;
    const isAuthed = this.auth.currentToken;
    // eslint-disable-next-line ember/no-controller-access-in-routes
    const controller = this.controllerFor('vault.cluster.auth');

    if (isAuthed && queryParams.redirect_to) {
      // if authenticated and redirect exists, redirect to that place
      transition = this.router.transitionTo(queryParams.redirect_to);
    } else if (isAuthed) {
      // if authed no redirect, go to cluster
      transition = this.router.transitionTo(CLUSTER);
    } else {
      // default go to Auth
      transition = this.router.transitionTo(AUTH);
    }
    transition.followRedirects().then(() => {
      controller.set('redirectTo', '');
    });
  }
}
