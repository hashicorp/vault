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
  @tracked warning;

  async postAdminGenerate() {
    this.error = null;
    let adapter = this.store.adapterFor('mfa-setup');
    let response;
    try {
      response = await adapter.adminGenerate({
        entity_id: this.args.entityId,
        method_id: this.UUID,
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

  @action
  async verifyUUID(evt) {
    evt.preventDefault();
    this.args.saveUUID(this.UUID); // send UUID to the parent. Needs to record in case of reset method.
    let response = await this.postAdminGenerate();
    if (response === 'stop_progress') {
      this.args.isUUIDVerified(false);
    } else if (response === 'reset_method') {
      this.args.goToReset(this.warning);
    } else {
      this.args.isUUIDVerified(true);
    }
  }
  @action
  redirectPreviousPage() {
    window.history.back();
  }
}
