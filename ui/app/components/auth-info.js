import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { later } from '@ember/runloop';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module AuthInfo
 *
 * @example
 * ```js
 * <AuthInfo @activeClusterName={{cluster.name}} @onLinkClick={{action "onLinkClick"}} />
 * ```
 *
 * @param {string} activeClusterName - name of the current cluster, passed from the parent.
 * @param {Function} onLinkClick - parent action which determines the behavior on link click
 */
export default class AuthInfoComponent extends Component {
  @service auth;
  @service wizard;
  @service router;

  @tracked fakeRenew = false;

  get hasEntityId() {
    // root users will not have an entity_id because they are not associated with an entity.
    // in order to use the MFA end user setup they need an entity_id
    return !!this.auth.authData.entity_id;
  }

  get isRenewing() {
    return this.fakeRenew || this.auth.isRenewing;
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  restartGuide() {
    this.wizard.restartGuide();
  }

  @action
  renewToken() {
    this.fakeRenew = true;
    later(() => {
      this.auth.renew().then(() => {
        this.fakeRenew = this.auth.isRenewing;
      });
    }, 200);
  }

  @action
  revokeToken() {
    this.auth.revokeCurrentToken().then(() => {
      this.transitionToRoute('vault.cluster.logout');
    });
  }
}
