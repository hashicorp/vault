/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import FlashMessageService from 'vault/services/flash-messages';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import PkiActionModel from 'vault/models/pki/action';

/**
 * @module PkiImportPemBundle
 * PkiImportPemBundle components are used to import PKI CA certificates and keys via pem_bundle.
 * https://github.com/hashicorp/vault/blob/main/website/content/api-docs/secret/pki.mdx#import-ca-certificates-and-keys
 *
 * @example
 * ```js
 * <PkiImportPemBundle @model={{this.model}} />
 * ```
 *
 * @param {Object} model - certificate model from route
 * @callback onCancel - Callback triggered when cancel button is clicked.
 * @callback onSubmit - Callback triggered on submit success.
 */

interface AdapterOptions {
  actionType: string;
  useIssuer: boolean | undefined;
}
interface Args {
  onSave: CallableFunction | null;
  onCancel: CallableFunction;
  onComplete: CallableFunction;
  model: PkiActionModel;
  adapterOptions: AdapterOptions;
}

export default class PkiImportPemBundle extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorBanner = '';
  @tracked afterImport = false;

  get importedResponse() {
    // mapping only exists after success
    // TODO: handle issuer already exists, but key doesn't -- empty object returned here
    return this.args.model.mapping;
  }

  @task
  @waitFor
  *submitForm(event: Event) {
    event.preventDefault();
    this.errorBanner = '';
    if (!this.args.model.pemBundle) {
      this.errorBanner = 'please upload your PEM bundle';
      return;
    }
    try {
      yield this.args.model.save({ adapterOptions: this.args.adapterOptions });
      this.flashMessages.success('Successfully imported data.');
      // This component shows the results, but call `onSave` for any side effects on parent
      if (this.args.onSave) {
        this.args.onSave();
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
    }
  }

  @action
  onFileUploaded({ value }: { value: string }) {
    this.args.model.pemBundle = value;
  }

  @action
  cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }
}
