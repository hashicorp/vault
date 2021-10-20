import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { run } from '@ember/runloop';
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

  @tracked
  isRenewing = false;

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  restartGuide() {
    this.wizard.restartGuide();
  }

  @action
  renewToken() {
    this.isRenewing = true;
    console.log('1');
    run.later(() => {
      this.auth.renew();
      this.isRenewing = false;
      console.log('2');
    }, 200);
  }

  @action
  revokeToken() {
    this.auth.revokeCurrentToken().then(() => {
      this.transitionToRoute('vault.cluster.logout');
    });
  }
}
