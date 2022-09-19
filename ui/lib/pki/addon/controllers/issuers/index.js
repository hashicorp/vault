import Controller from '@ember/controller';
import { action } from '@ember/object';

export default class PkiRolesIssuerController extends Controller {
  // To prevent production build bug of passing D.actions to on "click": https://github.com/hashicorp/vault/pull/16983
  @action onLinkClick(D) {
    D.actions.close();
  }
}
