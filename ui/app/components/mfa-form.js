import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
/**
 * @module MfaForm
 * The MfaForm component is used to enter a passcode when mfa is required to login
 *
 * @example
 * ```js
 * <MfaForm @clusterId={this.model.id} @authData={this.authData} />
 * ```
 * @param {string} clusterId - id of selected cluster
 * @param {object} authData - data from initial auth request -- { mfa_enforcement, backend, data }
 * @param {function} onSuccess - fired when passcode passes validation
 */

export default class MfaForm extends Component {
  @service auth;

  @tracked passcode;

  @task *validate() {
    try {
      const response = yield this.auth.totpValidate(
        { clusterId: this.args.clusterId, ...this.args.authData },
        this.passcode
      );
      this.args.onSuccess(response);
    } catch (error) {
      console.log(error);
      // do something
    }
  }

  @action submit(e) {
    e.preventDefault();
    this.validate.perform();
  }
}
