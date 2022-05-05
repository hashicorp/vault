import Component from '@glimmer/component';
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
  @action
  redirectPreviousPage() {
    this.args.restartFlow();
  }

  @action
  async restartSetup() {
    let response = await this.postAdminDestroy();
    return response; // ARG revisit
  }

  @action
  verifyAuthenticationCode(evt) {
    evt.preventDefault();
    // let authenticationCode = this.authenticationCode;
    // ARG TODO verify the UUID;
    // if verified send confirm boolean to the parent?
    this.args.isQRCodeVerified(true);
  }

  async postAdminGenerate() {
    this.error = null;
    let adapter = this.store.adapterFor('mfa-setup');
    let response;
    try {
      response = await adapter.adminDestroy({
        entity_id: this.args.entityId,
        method_id: this.args.uuid,
      });
      // if there was a warning it won't fail but needs to be handled here
      let warnings = response.warnings || [];
      if (warnings.length > 0) {
        this.UUID = ''; // clear UUID
        const alreadyGenerated = warnings.find((w) =>
          w.includes('Entity already has a secret for MFA method')
        );
        if (alreadyGenerated) {
          // replace warning because it comes in with extra quotes: "Entity already has a secret for MFA method ""' "
          // ARG TODO confirm with Ivana on language
          this.warning = 'Entity already has a secret for MFA method';
          return 'reset_method';
        }
        this.warning = warnings; // in case other kinds of warnings comes through. Still push to third screen because it's not an error.
        return 'reset_method';
      }
    } catch (error) {
      this.UUID = ''; // clear the UUID
      this.error = error.errors;
      return 'stop_progress';
    }
    return response;
  }
}
