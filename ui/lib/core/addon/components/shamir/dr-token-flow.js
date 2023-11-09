/**
 * Copyright (c) HashiCorp, Inc.
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
 * <Shamir::DrTokenFlow @action="generate-dr-operation-token" @onCancel={{this.closeModal}} />
 *
 * @param {string} action - required kebab-case-string which refers to an action within the cluster adapter
 * @param {function} onCancel - if provided, function will be triggered on Cancel
 */
export default class ShamirDrTokenFlowComponent extends ShamirFlowComponent {
  @tracked generateWithPGP = false; // controls which form shows
  @tracked savedPgpKey = null;
  @tracked otp = '';

  constructor() {
    super(...arguments);
    // Fetch status on init
    this.attemptProgress();
  }

  reset() {
    this.generateWithPGP = false;
    this.savedPgpKey = null;
    this.otp = '';
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
      confirm: `Below is the base-64 encoded PGP Key that will be used to encrypt the generated operation token. Next we'll enter portions of the root key to generate an operation token. Click the "Generate operation token" button to proceed.`,
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
    this.attemptProgress(this.extractData({ attempt: true }));
  }

  @action
  startGenerate(evt) {
    evt.preventDefault();
    this.attemptProgress(this.extractData({ attempt: true }));
  }

  @action
  async onCancelClose() {
    if (!this.encodedToken && this.started) {
      const adapter = this.store.adapterFor('cluster');
      await adapter.generateDrOperationToken({}, { cancel: true });
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
