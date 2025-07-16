/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import { service } from '@ember/service';
import trimRight from 'vault/utils/trim-right';
import errorMessage from 'vault/utils/error-message';
import type PkiKeyModel from 'vault/models/pki/key';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module PkiKeyImport
 * PkiKeyImport components are used to import PKI keys.
 *
 * @example
 * ```js
 * <PkiKeyImport @model={{this.model}} />
 * ```
 *
 * @param {Object} model - pki/key model.
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSubmit - Callback triggered on submit success.
 */
interface Args {
  model: PkiKeyModel;
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

export default class PkiKeyImport extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  @task
  @waitFor
  *submitForm(event: Event) {
    event.preventDefault();
    try {
      const { keyName } = this.args.model;
      yield this.args.model.save({ adapterOptions: { import: true } });
      this.flashMessages.success(`Successfully imported key${keyName ? ` ${keyName}.` : '.'}`);
      this.args.onSave();
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was a problem importing key.';
    }
  }

  @action
  onFileUploaded({ value, filename }: { value: string; filename: string }) {
    this.args.model.pemBundle = value;
    if (!this.args.model.keyName) {
      const trimmedFileName = trimRight(filename, ['.json', '.pem']);
      this.args.model.keyName = trimmedFileName;
    }
  }

  @action
  cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }
}
