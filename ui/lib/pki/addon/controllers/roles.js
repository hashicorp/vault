import Controller from '@ember/controller';
import { inject as service } from '@ember/service';

export default class PkiRolesController extends Controller {
  @service router;

  get showSecretListHeader() {
    return this.router.currentRouteName === 'vault.cluster.secrets.backend.pki.roles.index' ? true : false;
  }
}
