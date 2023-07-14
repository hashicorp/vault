/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * @module ShamirModalFlow
 * ShamirModalFlow is an extension of the ShamirFlow component that does the Generate Action Token workflow inside of a Modal.
 * Please note, this is not an extensive list of the required parameters -- please see ShamirFlow for others
 *
 * @example
 * ```js
 * <ShamirModalFlow @onClose={action 'onClose'}>This copy is the main paragraph when the token flow has not started</ShamirModalFlow>
 * ```
 * @param {function} onClose - This function will be triggered when the modal intends to be closed
 */
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ShamirFlow from './shamir/flow';

const pgpKeyFileDefault = () => ({ value: '' });
export default class ShamirModalFlow extends ShamirFlow {
  @tracked started = false;
  @tracked generateAction = '';
  @tracked generateWithPGP = false;
  @tracked haveSavedPGPKey = false;
  @tracked pgpKeyFile = pgpKeyFileDefault();

  constructor() {
    super(...arguments);
    this.startGenerate();
  }

  async startGenerate(data, evt) {
    console.log({ data, evt });
    const action = this.action;
    const adapter = this.store.adapterFor('cluster');
    const method = adapter[action];
    try {
      const resp = await method.call(adapter, {}, { checkStatus: true });
      this.updateProgress(resp);
      this.checkComplete(resp);
      return;
    } catch (e) {
      if (e.httpStatus === 400) {
        this.errors = e.errors;
        return;
      } else {
        // if licensing error, trigger parent method to handle
        if (e.httpStatus === 500 && e.errors?.join(' ').includes('licensing is in an invalid state')) {
          this.onLicenseError();
        }
        throw e;
      }
    }
  }

  get generateAction() {
    // TODO: this is redundant, this component is specific to this action
    return this.args.action === 'generate-dr-operation-token';
  }

  get generateStep() {
    const { generateWithPGP, attemptResponse, haveSavedPGPKey } = this;
    if (!generateWithPGP && !attemptResponse?.pgp_key) {
      return 'chooseMethod';
    }
    if (generateWithPGP) {
      if (attemptResponse?.pgp_key && haveSavedPGPKey) {
        return 'beginGenerationWithPGP';
      } else {
        return 'providePGPKey';
      }
    }
    return '';
  }
  get encoded_token() {
    return this.attemptResponse?.encoded_token;
  }
  get started() {
    return this.attemptResponse?.started;
  }

  @action setKey(_, keyFile) {
    this.pgpKey = keyFile.value;
    this.pgpKeyFile = keyFile;
  }

  @action
  onCancelClose() {
    if (this.attemptResponse.encoded_token) {
      this.send('reset');
    } else if (this.generateAction && !this.started) {
      if (this.generateStep !== 'chooseMethod') {
        this.send('reset');
      }
    } else {
      const adapter = this.store.adapterFor('cluster');
      adapter.generateDrOperationToken({}, { cancel: true });
      this.send('reset');
    }
    this.args.onClose();
  }
}
