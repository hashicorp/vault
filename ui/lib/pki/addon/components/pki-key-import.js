/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import trimRight from 'vault/utils/trim-right';
import { waitFor } from '@ember/test-waiters';

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

export default class PkiKeyImport extends Component {
  @service flashMessages;

  @tracked errorBanner;
  @tracked invalidFormAlert;

  @task
  @waitFor
  *submitForm(event) {
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
  onFileUploaded({ value, filename }) {
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
