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

  async postAdminGenerate() {
    this.error = null;
    let adapter = this.store.adapterFor('mfa-setup');
    let response;
    try {
      response = await adapter.adminGenerate({
        entity_id: this.args.entityId,
        method_id: this.UUID,
      });
    } catch (error) {
      const errors = error.errors || [];
      const incorrectMethodID = errors.find((e) => e.includes('missing method ID'));
      const restartSetup = errors.find((e) => e.includes('Entity already has a secret'));
      // ARG TODO confirm that this clears the errors

      if (incorrectMethodID) {
        this.error = 'You have used an incorrect Method ID. Contact your administrator.';
        return 'stop_progress';
      } else if (restartSetup) {
        // ARG TODO switch screens to step 3 and show error message
      } else {
        // no custom error message, return adapter error
        // ARG TODO test when root and no entity_id
        this.error = error.errors;
        return 'stop_progress';
      }
    }
    // ARG TODO can check if returns already generated here
    return response;
  }

  @action
  async verifyUUID(evt) {
    evt.preventDefault();
    let response = await this.postAdminGenerate();
    if (response === 'stop_progress') {
      this.args.isUUIDVerified(false);
    } else {
      this.args.isUUIDVerified(true);
    }
  }
  @action
  redirectPreviousPage() {
    window.history.back();
  }
}
