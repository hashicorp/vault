import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiRoute extends Route {
  @service router;

  redirect() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.overview');
  }
}
