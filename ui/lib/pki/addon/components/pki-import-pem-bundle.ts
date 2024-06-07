/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import type FlashMessageService from 'vault/services/flash-messages';
import type PkiActionModel from 'vault/models/pki/action';

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
 * @callback onSave - Callback triggered on submit success.
 * @callback onComplete - Callback triggered on "done" button click.
 */

interface AdapterOptions {
  actionType: string;
  useIssuer: boolean | undefined;
}
interface Args {
  onSave?: CallableFunction;
  onCancel: CallableFunction;
  onComplete: CallableFunction;
  model: PkiActionModel;
  adapterOptions: AdapterOptions;
}

export default class PkiImportPemBundle extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorBanner = '';

  get importedResponse() {
    const { mapping, importedIssuers, importedKeys } = this.args.model;
    // Even if there are no imported items, mapping will be an empty object from API response
    if (undefined === mapping) return null;

    const importList = (importedIssuers || []).map((issuer: string) => {
      const key = mapping[issuer];
      return { issuer, key };
    });

    // Check each imported key and make sure it's in the list
    (importedKeys || []).forEach((key) => {
      const matchIdx = importList.findIndex((item) => item.key === key);
      // If key isn't accounted for, add it without a matching issuer
      if (matchIdx === -1) {
        importList.push({ issuer: '', key });
      }
    });

    if (importList.length === 0) {
      // If no new items were imported but the import call was successful, the UI will show accordingly
      return [{ issuer: '', key: '' }];
    }
    return importList;
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
      window?.scrollTo(0, 0);
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
