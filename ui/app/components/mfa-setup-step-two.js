import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

/**
 * @module MfaSetupStepTwo
 * MfaSetupStepTwo component is a child component used in the end user setup for MFA. It displays a qrCode or a warning and allows a user to reset the method.
 *
 * @param {string} entityId - the entityId of the user. This comes from the auth service which records it on loading of the cluster. A root user does not have an entityId.
 * @param {string} uuid - the UUID that is entered in the input on step one.
 * @param {string} qrCode - the returned url from the admin-generate post. Used to create the qrCode.
 * @param {boolean} restartFlow - a boolean that is true that is true if the user should proceed to step two or false if they should stay on step one.
 * @param {string} warning - if there is a warning returned from the admin-generate post then it's sent to the step two component in this param.
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
    } catch (error) {
      this.error = error.errors;
      return 'stop_progress';
    }
    this.args.restartFlow();
  }
}
