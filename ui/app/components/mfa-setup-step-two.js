import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

/**
 * @module MfaSetupStepTwo
 * MfaSetupStepTwo components are used to...
 *
 * @example
 * ```js
 * <MfaSetupStepTwo @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class MfaSetupStepTwo extends Component {
  @service store;

  @action
  redirectPreviousPage() {
    this.args.restartFlow();
    window.history.back();
  }

  @action
  async restartSetup() {
    this.error = null;
    let adapter = this.store.adapterFor('mfa-setup');
    try {
      await adapter.adminDestroy({
        entity_id: this.args.entityId,
        method_id: this.args.uuid,
      });
      // if there was a warning it won't fail but needs to be handled here
    } catch (error) {
      this.error = error.errors;
      return 'stop_progress';
    }
    // restart to step one.
    this.args.restartFlow();
  }

  @action
  verifyAuthenticationCode(evt) {
    evt.preventDefault();
    // let authenticationCode = this.authenticationCode;
    // ARG TODO verify the UUID;
    // if verified send confirm boolean to the parent?
    this.args.isQRCodeVerified(true);
  }
}
