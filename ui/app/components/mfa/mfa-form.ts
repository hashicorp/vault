/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task, timeout } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import MfaConstraint from 'vault/resources/mfa/constraint';

import type AuthService from 'vault/vault/services/auth';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type { AuthSuccessResponse } from 'vault/vault/services/auth';
import type { Task } from 'ember-concurrency';
import type { MfaAuthData } from 'vault/vault/auth/mfa';

/**
 * @module MfaForm
 * The MfaForm component is used to enter a passcode when mfa is required to login
 *
 * @example
 * ```js
 * <Mfa::MfaForm @clusterId={this.model.id} @authData={this.authData} />
 * ```
 * @param {string} clusterId - id of selected cluster
 * @param {object} authData - data from initial auth request -- { mfaRequirement, backend, data }
 * @param {function} loginAndTransition - fired when passcode passes validation
 * @param {function} onCancel - callback that clears @authData which resets MFA and returns to the login form
 */

export const TOTP_VALIDATION_ERROR =
  'The passcode failed to validate. If you entered the correct passcode, contact your administrator.';

interface Args {
  authData: MfaAuthData;
  clusterId: string;
  loginAndTransition: Task<AuthSuccessResponse, [AuthSuccessResponse]>;
  onCancel: CallableFunction;
}

export default class MfaForm extends Component<Args> {
  @service declare readonly auth: AuthService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  @tracked constraints: MfaConstraint[] = [];
  @tracked codeDelayMessage = '';
  @tracked countdown = 0;
  @tracked error = '';
  // Self-enrollment is per MFA method, not per login enforcement (constraint)
  // Track method IDs used to fetch a QR code so we don't re-request if a user just enrolled.
  @tracked enrolledMethods = new Set<string>();

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    const { mfa_constraints = [] } = this.args.authData.mfaRequirement;
    this.constraints = mfa_constraints.map((constraint) => new MfaConstraint(constraint));

    // Trigger validation immediately if passcode or user selection is not required
    this.checkStateAndValidate();

    // Filter for constraints that have only one MFA method and it supports self-enrollment
    const filteredConstraints = this.constraints.filter((c) => c.selfEnrollMethod && c.methods.length === 1);
    // If there is only one then fetch the QR code because self-enrolling is unavoidable.
    if (filteredConstraints.length === 1 && filteredConstraints[0]) {
      const [constraint] = filteredConstraints;
      const method = constraint.selfEnrollMethod;
      if (method) this.fetchQrCode.perform(method.id, constraint);
    }
  }

  get everyConstraintSatisfied() {
    return this.constraints.every((constraint) => constraint.isSatisfied);
  }

  get currentSelfEnrollConstraint() {
    return this.constraints.find((c) => c.qrCode !== '');
  }

  get needsToChoose() {
    // If any self-enroll constraints are missing selections
    const missingSelfEnrollSelection = this.constraints
      .filter((c) => !!c.selfEnrollMethod)
      .some((c) => !c.selectedMethod);
    // If there is only one constraint but it has multiple methods
    const missingSelection = this.constraints.length === 1 && !this.constraints.some((c) => c.selectedMethod);
    return missingSelfEnrollSelection || missingSelection;
  }

  // There is only one login enforcement with only one method configured
  get singleLoginEnforcement() {
    if (this.constraints.length === 1) {
      const loginEnforcement = this.constraints[0];
      // Return a value if there is only one MFA method configured.
      return loginEnforcement?.methods.length === 1 ? loginEnforcement : null;
    }
    return null;
  }

  validate = task(async () => {
    const { authMethodType, authMountPath, mfaRequirement } = this.args.authData;
    const submitData = {
      mfa_request_id: mfaRequirement.mfa_request_id,
      mfa_constraints: this.constraints.map((c) => c.validateData),
    };
    try {
      this.error = '';
      const response = await this.auth.totpValidate({
        clusterId: this.args.clusterId,
        authMethodType,
        authMountPath,
        mfaRequirement: submitData,
      });

      // calls loginAndTransition in auth.js controller
      this.args.loginAndTransition.perform(response);
    } catch (error) {
      // Reset enrolled methods if there's an error
      this.enrolledMethods = new Set<string>();
      const errorMsg = errorMessage(error);
      const codeUsed = errorMsg.includes('code already used');
      const rateLimit = errorMsg.includes('maximum TOTP validation attempts');
      const delayMessage = codeUsed || rateLimit ? errorMsg : null;
      if (delayMessage) {
        const reason = codeUsed ? 'This code has already been used' : 'Maximum validation attempts exceeded';
        this.codeDelayMessage = `${reason}. Please wait until a new code is available.`;
        this.newCodeDelay.perform(delayMessage);
      } else if (this.singleLoginEnforcement?.selectedMethod?.uses_passcode) {
        this.error = TOTP_VALIDATION_ERROR;
      } else {
        this.error = errorMsg;
      }
    }
  });

  fetchQrCode = task(async (mfa_method_id: string, constraint: MfaConstraint) => {
    // Self-enrollment is an enterprise only feature
    if (this.version.isCommunity) return;

    const adapter = this.store.adapterFor('application');
    const { mfaRequirement } = this.args.authData;
    try {
      const { data } = await adapter.ajax('/v1/identity/mfa/method/totp/self-enroll', 'POST', {
        unauthenticated: true,
        data: { mfa_method_id, mfa_request_id: mfaRequirement.mfa_request_id },
      });
      if (data?.url) {
        // Set QR code which recomputes currentSelfEnrollConstraint and renders it for the user to scan
        constraint.qrCode = data.url;
        // Add mfa_method_id to list of already enrolled methods for client-side tracking
        this.enrolledMethods.add(mfa_method_id);
        return;
      }
      // Not sure it's realistic to get here without the endpoint throwing an error, but just in case!
      this.error = 'There was a problem generating the QR code. Please try again.';
    } catch (error) {
      this.error = errorMessage(error);
    }
  });

  newCodeDelay = task(async (errorMessage) => {
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
      await timeout(1000);
      this.countdown--;
    }
  });

  @action
  async onSelect(constraint: MfaConstraint, methodId: string) {
    // Set selectedMethod on the MfaConstraint class
    // If id is an empty string, it clears the selected method
    constraint.setSelectedMethod(methodId);

    const selectedMethod = constraint.selectedMethod;
    if (selectedMethod?.self_enrollment_enabled && !this.methodAlreadyEnrolled(selectedMethod.id)) {
      await this.fetchQrCode.perform(selectedMethod.id, constraint);
    }

    this.checkStateAndValidate();
  }

  @action
  checkStateAndValidate() {
    // Whenever all login enforcements are satisfied perform validation to save the user extra clicks
    if (this.everyConstraintSatisfied) {
      this.validate.perform();
    }
  }

  // Template helpers
  methodAlreadyEnrolled = (methodId: string) => this.enrolledMethods.has(methodId);
}
