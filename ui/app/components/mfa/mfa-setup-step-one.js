import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module MfaSetupStepOne
 * MfaSetupStepOne component is a child component used in the end user setup for MFA. It records the UUID (aka method_id) and sends a admin-generate request.
 *
 * @param {string} entityId - the entityId of the user. This comes from the auth service which records it on loading of the cluster. A root user does not have an entityId.
 * @param {function} isUUIDVerified - a function that consumes a boolean. Is true if the admin-generate is successful and false if it throws a warning or error.
 * @param {boolean} restartFlow - a boolean that is true that is true if the user should proceed to step two or false if they should stay on step one.
 * @param {function} saveUUIDandQrCode - A function that sends the inputted UUID and return qrCode from step one to the parent.
 * @param {boolean} showWarning - whether a warning is returned from the admin-generate query. Needs to be passed to step two.
 */

export default class MfaSetupStepOne extends Component {
  @service store;
  @tracked error = '';
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
      this.args.saveUUIDandQrCode(this.UUID, response.data?.url);
      // if there was a warning it won't fail but needs to be handled here and the flow needs to be interrupted
      let warnings = response.warnings || [];
      if (warnings.length > 0) {
        this.UUID = ''; // clear UUID
        const alreadyGenerated = warnings.find((w) =>
          w.includes('Entity already has a secret for MFA method')
        );
        if (alreadyGenerated) {
          this.warning =
            'A QR code has already been generated, scanned, and MFA set up for this entity. If a new code is required, contact your administrator.';
          return 'reset_method';
        }
        this.warning = warnings; // in case other kinds of warnings comes through.
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
