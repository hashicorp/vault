/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import { camelize } from '@ember/string';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module ShamirFlowComponent
 * These components are used to manage keeping track of a shamir unseal flow.
 * This component is generic and can be overwritten for various shamir use cases.
 * The lifecycle for a Shamir flow is as follows:
 * 1. Start (optional)
 * 2. Attempt progress
 * 3. Check progress
 * 4. Check complete
 *
 * @example
 * ```js
 * <Shamir::Flow
 *  @action="unseal"
 *  @threshold={{5}}
 *  @progress={{3}}
 *  @onShamirSuccess={{transition-to "vault.cluster"}}
 * />
 * ```
 *
 * @param {string} action - adapter method name (kebab case) to call on attempt
 * @param {number} threshold - number of keys required to unlock
 * @param {number} progress - number of keys given so far for unlock
 * @param {string} inputLabel - (optional) Label for key input
 * @param {string} buttonText - (optional) CTA for the form submit button. Defaults to "Submit"
 * @param {Function} extractData - (optional) modify the payload before the action is called
 * @param {Function} updateProgress - (optional) call a side effect to check if progress has been made
 * @param {Function} checkComplete - (optional) custom logic based on adapter response. Should return boolean.
 * @param {Function} onShamirSuccess - method called when shamir unlock is complete.
 * @param {Function} onLicenseError - method called when shamir unlock fails due to licensing error
 *
 */
export default class ShamirFlowComponent extends Component {
  @service store;
  @tracked errors = null;
  @tracked attemptResponse = null;

  get action() {
    if (!this.args.action) return '';
    return camelize(this.args.action);
  }

  extractData(data) {
    if (this.args.extractData) {
      // custom data extraction
      return this.args.extractData(data);
    }

    // This method can be overwritten by extended components
    // to control what data is passed into the method action
    if (this.attemptResponse?.nonce) {
      data.nonce = this.attemptResponse.nonce;
    }
    return data;
  }

  /**
   * 2. Attempt progress. This method assumes the correct data
   * has already been extracted (use this.extractData to customize)
   * @param {object} data arbitrary data which will be passed to adapter method
   * @returns Promise which should resolve unless throwing error to parent.
   */
  async attemptProgress(data) {
    this.errors = null;
    const action = this.action;
    const adapter = this.store.adapterFor('cluster');
    const method = adapter[action];
    // Only used for DR token generate
    const checkStatus = data ? false : true;

    try {
      const resp = await method.call(adapter, data, { checkStatus });
      this.updateProgress(resp);
      this.handleComplete(resp);
      return;
    } catch (e) {
      if (e.httpStatus === 400) {
        this.errors = e.errors;
        return;
      } else {
        // if licensing error, trigger parent method to handle
        if (e.httpStatus === 500 && e.errors?.join(' ').includes('licensing is in an invalid state')) {
          this.args.onLicenseError();
        }
        throw e;
      }
    }
  }

  /**
   * 3. This method gets called after successful unseal attempt.
   * By default the response will be made available to the component,
   * but pass in @updateProgress (no params) to trigger any side effects that will
   * update passed attributes from parent.
   * @param {payload} response from the adapter method
   * @returns void
   */
  updateProgress(response) {
    if (this.args.updateProgress) {
      this.args.updateProgress();
    }
    this.attemptResponse = response;
    return;
  }

  /**
   * 4. checkComplete checks the payload for completeness, then then
   * takes calls @onShamirSuccess with no arguments if complete.
   * For custom logic, define @checkComplete which receives the
   * adapter payload.
   * @param {payload} response from the adapter method
   * @returns void
   */
  handleComplete(response) {
    const isComplete = this.checkComplete(response);
    if (isComplete) {
      if (this.args.onShamirSuccess) {
        this.args.onShamirSuccess();
      }
    }
    return;
  }

  checkComplete(response) {
    if (this.args.checkComplete) {
      return this.args.checkComplete(response);
    }
    return response.complete === true;
  }

  reset() {
    this.attemptResponse = null;
    this.errors = null;
  }

  @action
  onSubmitKey(data) {
    this.attemptProgress(this.extractData(data));
  }
}

/* example unseal response (progress)
{
  "sealed": true,
  "t": 3,
  "n": 5,
  "progress": 2,
  "version": "0.6.2"
}

example unseal response (finished)
{
  "sealed": false,
  "t": 3,
  "n": 5,
  "progress": 0,
  "version": "0.6.2",
  "cluster_name": "vault-cluster-d6ec3c7f",
  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8"
}
*/
