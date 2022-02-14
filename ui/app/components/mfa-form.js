import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action, set } from '@ember/object';
import { task, timeout } from 'ember-concurrency';
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

export default class MfaForm extends Component {
  @service auth;

  @tracked passcode;
  @tracked countdown;
  @tracked errors;

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
      base += ` ${this.constraints.length} methods are required for successful authentication.`;
    }
    return base;
  }

  @task *validate() {
    try {
      const response = yield this.auth.totpValidate({
        clusterId: this.args.clusterId,
        ...this.args.authData,
      });
      this.args.onSuccess(response);
    } catch (error) {
      this.errors = error.errors;
      // update if specific error can be parsed for incorrect passcode
      // this.newCodeDelay.perform();
    }
  }

  @task *newCodeDelay() {
    this.passcode = null;
    this.countdown = 30;
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
