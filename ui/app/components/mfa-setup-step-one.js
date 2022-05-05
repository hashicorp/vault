import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module MfaSetupStepOne
 * MfaSetupStepOne components are used to...
 *
 * @example
 * ```js
 * <MfaSetupStepOne @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} isUUIDVerified - action from parent that returns a boolean from the child regarding whether or not the UUID was verified. If true then proceed to step 2.
 */

export default class MfaSetupStepOne extends Component {
  @service store;
  @tracked error;
  @tracked warning = '';
  @tracked qrCode = '';

  @action
  redirectPreviousPage() {
    this.args.restartFlow();
    window.history.back();
  }

  @action
  async verifyUUID(evt) {
    evt.preventDefault();
    let response = await this.postAdminGenerate();
    if (response === 'stop_progress') {
      this.args.isUUIDVerified(false);
    } else if (response === 'reset_method') {
      this.args.showWarning(this.warning);
    } else {
      this.args.isUUIDVerified(true);
    }
  }

  async postAdminGenerate() {
    this.error = '';
    this.warning = '';
    let adapter = this.store.adapterFor('mfa-setup');
    let response;
    try {
      response = await adapter.adminGenerate({
        entity_id: this.args.entityId,
        method_id: this.UUID, // comes from value on the input
      });
      this.args.saveUUIDandQrCode(this.UUID, response.data?.url); // parent needs to keep track of UUID and qrCode.
      // if there was a warning it won't fail but needs to be handled here and the flow needs to be interrupted
      let warnings = response.warnings || [];
      if (warnings.length > 0) {
        this.UUID = ''; // clear UUID
        const alreadyGenerated = warnings.find((w) =>
          w.includes('Entity already has a secret for MFA method')
        );
        if (alreadyGenerated) {
          // replace warning because it comes in with extra quotes: "Entity already has a secret for MFA method ""' "
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
