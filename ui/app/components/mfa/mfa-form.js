/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { service } from '@ember/service';
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
 * <Mfa::MfaForm @clusterId={this.model.id} @authData={this.authData} />
 * ```
 * @param {string} clusterId - id of selected cluster
 * @param {object} authData - data from initial auth request -- { mfa_requirement, backend, data }
 * @param {function} onSuccess - fired when passcode passes validation
 * @param {function} onError - fired for multi-method or non-passcode method validation errors
 */

export const TOTP_VALIDATION_ERROR =
  'The passcode failed to validate. If you entered the correct passcode, contact your administrator.';

export default class MfaForm extends Component {
  @service auth;

  @tracked countdown = 0;
  @tracked error;
  @tracked codeDelayMessage;

  constructor() {
    super(...arguments);
    // trigger validation immediately when passcode is not required
    const passcodeOrSelect = this.constraints.filter((constraint) => {
      return constraint.methods.length > 1 || constraint.methods.find((m) => m.uses_passcode);
    });
    if (!passcodeOrSelect.length) {
      this.validate.perform();
    }
  }

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
      // calls onMfaSuccess in auth/page.js
      this.args.onSuccess(response);
    } catch (error) {
      const errors = error.errors || [];
      const codeUsed = errors.find((e) => e.includes('code already used'));
      const rateLimit = errors.find((e) => e.includes('maximum TOTP validation attempts'));
      const delayMessage = codeUsed || rateLimit;

      if (delayMessage) {
        const reason = codeUsed ? 'This code has already been used' : 'Maximum validation attempts exceeded';
        this.codeDelayMessage = `${reason}. Please wait until a new code is available.`;
        this.newCodeDelay.perform(delayMessage);
      } else if (this.singlePasscode) {
        this.error = TOTP_VALIDATION_ERROR;
      } else {
        this.args.onError(this.auth.handleError(error));
      }
    }
  }

  @task *newCodeDelay(errorMessage) {
    let delay;

    // parse validity period from error string to initialize countdown
    const delayRegExMatches = errorMessage.match(/(\d+\w seconds)/);
    if (delayRegExMatches && delayRegExMatches.length) {
      delay = delayRegExMatches[0].split(' ')[0];
    } else {
      // default to 30 seconds if error message doesn't specify one
      delay = 30;
    }
    this.countdown = parseInt(delay);

    // skip countdown in testing environment
    if (Ember.testing) return;

    while (this.countdown > 0) {
      yield timeout(1000);
      this.countdown--;
    }
  }

  @action onSelect(constraint, id) {
    set(constraint, 'selectedId', id);
    set(
      constraint,
      'selectedMethod',
      constraint.methods.find((m) => m.id === id)
    );
  }
  @action submit(e) {
    e.preventDefault();
    this.validate.perform();
  }
}
