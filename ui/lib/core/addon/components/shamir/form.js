/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module ShamirFormComponent
 * These components are used to make progress against a Shamir seal.
 * Depending on the response, and external polling, the component will show
 * progress and optional
 *
 * @example
 * <Shamir::Form
 *  @onSubmit={{this.handleKeySubmit}}
 *  @threshold={{cluster.threshold}}
 *  @progress={{cluster.progress}}
 *  @fetchOnInit={{true}}
 *  @onShamirSuccess={{transition-to "vault.cluster"}}
 * />
 *
 *
 * @param {Function} onSubmit - method which handles the action for shamir. Receives { key }
 * @param {number} threshold - min number of keys required to unlock shamir seal
 * @param {number} progress - number of keys given so far for unlock
 * @param {string} buttonText - CTA for the form submit button. Defaults to "Submit"
 * @param {string} inputLabel - Label for key input. Defaults to "Shamir key portion"
 * @param {boolean} alwaysShowProgress - determines if the shamir progress should always show, or only when > 0 progress
 * @param {string} otp - if otp is present, it will show a section describing what to do with it
 *
 */
export default class ShamirFormComponent extends Component {
  @tracked key = '';
  @tracked loading = false;

  get buttonText() {
    return this.args.buttonText || 'Submit';
  }
  get showProgress() {
    return this.args.progress > 0 || this.args.alwaysShowProgress;
  }
  get inputLabel() {
    return this.args.inputLabel || 'Shamir key portion';
  }

  resetForm() {
    this.key = '';
    this.loading = false;
  }

  @action
  async onSubmit(key, evt) {
    evt.preventDefault();

    if (!key) {
      return;
    }
    // Parent handles action and passes in errors if present
    await this.args.onSubmit({ key });
    this.resetForm();
  }
}
