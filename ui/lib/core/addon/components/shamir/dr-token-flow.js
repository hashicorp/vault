/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import ShamirFlowComponent from './flow';

/**
 * @module ShamirDrTokenFlowComponent
 * ShamirDrTokenFlow is an extension of the ShamirFlow component that does the Generate Action Token workflow inside of a Modal.
 * Please note, this is not an extensive list of the required parameters -- please see ShamirFlow for others
 *
 * @example
 * ```js
 * <Shamir::DrTokenFlow @action="generate-dr-operation-token" @onCancel={{this.closeModal}} />
 * ```
 * @param {string} action - required kebab-case-string which refers to an action within the cluster adapter
 * @param {function} onCancel - if provided, function will be triggered on Cancel
 */
export default class ShamirDrTokenFlowComponent extends ShamirFlowComponent {
  @tracked generateWithPGP = false; // controls which form shows
  @tracked savedPgpKey = null;
  @tracked otp = '';
  @tracked askForPrimaryToken = false; // controls whether to show primary token input
  @tracked primaryRootToken = null; // stores the primary root token

  constructor() {
    super(...arguments);
    // Don't fetch status on init - we'll check it after the user provides the primary token
    // Fetching status here would start an unauthenticated generation attempt
  }

  reset() {
    this.generateWithPGP = false;
    this.savedPgpKey = null;
    this.otp = '';
    this.askForPrimaryToken = false;
    this.primaryRootToken = null;
    // tracked items on Shamir/Flow
    this.attemptResponse = null;
    this.errors = null;
  }

  // Values calculated from the attempt response
  get encodedToken() {
    return this.attemptResponse?.encoded_token;
  }
  get started() {
    return this.attemptResponse?.started;
  }
  get nonce() {
    return this.attemptResponse?.nonce;
  }
  get progress() {
    return this.attemptResponse?.progress;
  }
  get threshold() {
    return this.attemptResponse?.required;
  }
  get pgpText() {
    return {
      confirm: `Below is the base-64 encoded PGP Key that will be used to encrypt the generated operation token.`,
      form: `Choose a PGP Key from your computer or paste the contents of one in the form below. This key will be used to Encrypt the generated operation token.`,
    };
  }

  // Methods which override those in Shamir/Flow
  extractData(data) {
    if (this.started) {
      if (this.nonce) {
        data.nonce = this.nonce;
      }
      return data;
    }
    if (this.savedPgpKey) {
      return {
        pgp_key: this.savedPgpKey,
      };
    }
    // only if !started
    return {
      attempt: data.attempt,
    };
  }

  updateProgress(response) {
    if (response.otp) {
      // OTP is sticky -- once we get one we don't want to remove it
      // even if the current response doesn't include one.
      // See PR #5818
      this.otp = response.otp;
    }
    this.attemptResponse = response;
    return;
  }

  @action
  usePgpKey(keyfile) {
    this.savedPgpKey = keyfile;
    // Don't start generation yet - show primary token form first
    this.generateWithPGP = false;
    this.askForPrimaryToken = true;
  }

  @action
  onSubmitKey(data) {
    // Override parent to pass primaryToken
    this.attemptProgress(this.extractData(data), this.primaryRootToken);
  }

  @action
  startGenerate(evt) {
    evt.preventDefault();
    // Show the primary token input form first
    this.askForPrimaryToken = true;
  }

  @action
  updatePrimaryRootToken(evt) {
    this.primaryRootToken = evt.target.value;
  }

  @action
  async validatePrimaryRootToken() {
    if (!this.primaryRootToken) {
      this.errors = ['Primary root token is required'];
      return;
    }

    this.errors = null;

    try {
      // First, check status to validate the token without starting a new generation
      await this.attemptProgress(undefined, this.primaryRootToken);

      if (!this.started) {
        // No generation in progress, so start one
        await this.attemptProgress(this.extractData({ attempt: true }), this.primaryRootToken);

        // Check if there were errors from starting generation (e.g., invalid PGP key)
        if (this.errors) {
          return;
        }
      }

      // Only hide the primary token form if there were no errors
      this.askForPrimaryToken = false;
    } catch (e) {
      if (e.httpStatus === 403) {
        this.errors = ['Invalid primary root token. Please check the token and try again.'];
      } else {
        this.errors = [e.message || 'An error occurred while validating the token'];
      }
    }
  }

  @action
  backToPgpForm() {
    // Go back to PGP form and clear the saved PGP key
    this.askForPrimaryToken = false;
    this.generateWithPGP = true;
    this.savedPgpKey = null;
    this.errors = null;
  }

  @action
  async onCancelClose() {
    if (!this.encodedToken && this.started) {
      const adapter = this.store.adapterFor('cluster');
      await adapter.generateDrOperationToken({}, { cancel: true, token: this.primaryRootToken });
    }
    this.reset();
    if (this.args.onCancel) {
      this.args.onCancel();
    }
  }
}

/* generate-operation-token response example
{
  "started": true,
  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
  "progress": 1,
  "required": 3,
  "encoded_token": "",
  "otp": "2vPFYG8gUSW9npwzyvxXMug0",
  "otp_length": 24,
  "complete": false
}
*/
