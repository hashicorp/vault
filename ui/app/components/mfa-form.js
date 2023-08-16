import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action, set } from '@ember/object';
import { task, timeout } from 'ember-concurrency';
import { numberToWord } from 'vault/helpers/number-to-word';
/**
 * @module MfaForm
 * The MfaForm component is used to enter a passcode when mfa is required to login
 *
 * @example
 * ```js
 * <MfaForm @clusterId={this.model.id} @authData={this.authData} />
 * ```
 * @param {string} clusterId - id of selected cluster
 * @param {object} authData - data from initial auth request -- { mfa_requirement, backend, data }
 * @param {function} onSuccess - fired when passcode passes validation
 */

export const VALIDATION_ERROR =
  'The passcode failed to validate. If you entered the correct passcode, contact your administrator.';

export default class MfaForm extends Component {
  @service auth;

  @tracked countdown;
  @tracked error;

  get constraints() {
    return this.args.authData.mfa_requirement.mfa_constraints;
  }
  get multiConstraint() {
    return this.constraints.length > 1;
  }
  get singleConstraintMultiMethod() {
    return !this.isMultiConstraint && this.constraints[0].methods.length > 1;
  }
  get singlePasscode() {
    return (
      !this.isMultiConstraint &&
      this.constraints[0].methods.length === 1 &&
      this.constraints[0].methods[0].uses_passcode
    );
  }
  get description() {
    let base = 'Multi-factor authentication is enabled for your account.';
    if (this.singlePasscode) {
      base += ' Enter your authentication code to log in.';
    }
    if (this.singleConstraintMultiMethod) {
      base += ' Select the MFA method you wish to use.';
    }
    if (this.multiConstraint) {
      const num = this.constraints.length;
      base += ` ${numberToWord(num, true)} methods are required for successful authentication.`;
    }
    return base;
  }

  @task *validate() {
    try {
      this.error = null;
      const response = yield this.auth.totpValidate({
        clusterId: this.args.clusterId,
        ...this.args.authData,
      });
      this.args.onSuccess(response);
    } catch (error) {
      const codeUsed = (error.errors || []).find((e) => e.includes('code already used;'));
      if (codeUsed) {
        // parse validity period from error string to initialize countdown
        const seconds = parseInt(codeUsed.split('in ')[1].split(' seconds')[0]);
        this.newCodeDelay.perform(seconds);
      } else {
        this.error = VALIDATION_ERROR;
      }
    }
  }

  @task *newCodeDelay(timePeriod) {
    this.countdown = timePeriod;
    while (this.countdown) {
      yield timeout(1000);
      this.countdown--;
    }
  }

  @action onSelect(constraint, id) {
    set(constraint, 'selectedId', id);
    set(constraint, 'selectedMethod', constraint.methods.findBy('id', id));
  }
  @action submit(e) {
    e.preventDefault();
    this.validate.perform();
  }
}
