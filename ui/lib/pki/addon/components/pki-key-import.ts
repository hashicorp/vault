/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import { service } from '@ember/service';
import trimRight from 'vault/utils/trim-right';
import FormField from 'vault/utils/forms/field';

import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

/**
 * @module PkiKeyImport
 * PkiKeyImport components are used to import PKI keys.
 *
 *
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSubmit - Callback triggered on submit success.
 */
interface Args {
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

export default class PkiKeyImport extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked declare keyName: string;
  @tracked declare pemBundle: string;

  keyNameField = new FormField('keyName', 'string', {
    subText: `Optional, human-readable name for this key. The name must be unique across all keys and cannot be 'default'.`,
  });

  @task
  @waitFor
  *submitForm(event: Event) {
    event.preventDefault();
    try {
      yield this.api.secrets.pkiImportKey(this.secretMountPath.currentPath, {
        key_name: this.keyName,
        pem_bundle: this.pemBundle,
      });
      this.flashMessages.success(`Successfully imported key${this.keyName ? ` ${this.keyName}.` : '.'}`);
      this.args.onSave();
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.errorBanner = message;
      this.invalidFormAlert = 'There was a problem importing key.';
    }
  }

  @action
  onFileUploaded({ value, filename }: { value: string; filename: string }) {
    this.pemBundle = value;
    if (!this.keyName) {
      const trimmedFileName = trimRight(filename, ['.json', '.pem']);
      this.keyName = trimmedFileName;
    }
  }
}
