import Controller from '@ember/controller';
import { getOwner } from '@ember/application';
import { action } from '@ember/object';
import { next } from '@ember/runloop';

export default class PkiIssuerIndexController extends Controller {
  get mountPoint() {
    return getOwner(this).mountPoint;
  }
  // To prevent production build bug of passing D.actions to on "click": https://github.com/hashicorp/vault/pull/16983
  @action onLinkClick(D) {
    next(() => D.actions.close());
  }
}
