/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import RouterService from '@ember/routing/router-service';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import PkiActionModel from 'vault/models/pki/action';
import PkiUrlsModel from 'vault/models/pki/urls';
import FlashMessageService from 'vault/services/flash-messages';
import errorMessage from 'vault/utils/error-message';

interface AdapterOptions {
  actionType: string;
  useIssuer: boolean | undefined;
}
interface Args {
  model: PkiActionModel;
  urls: PkiUrlsModel;
  onCancel: CallableFunction;
  onSuccess: CallableFunction | undefined;
  adapterOptions: AdapterOptions;
}

/**
 * @module PkiGenerateRoot
 * PkiGenerateRoot shows only the fields valid for the generate root endpoint.
 * This component renders the form, handles the model save and rollback actions,
 * and shows the resulting data on success. onCancel is required for the cancel
 * transition, and if onSuccess is provided it will call that after save for any
 * side effects in the parent.
 *
 * @example
 * ```js
 * <PkiGenerateRoot @model={{this.model}} @onCancel={{transition-to "vault.cluster"}} @onSave={{transition-to "vault.cluster.secrets"}} @adapterOptions={{hash actionType="import" useIssuer=false}} />
 * ```
 *
 * @param {Object} model - pki/action model.
 * @callback onCancel - Callback triggered when cancel button is clicked, after model is unloaded
 * @callback onSuccess - Optional - Callback triggered after model is saved. Results are shown on the same component
 * @param {Object} adapterOptions - object passed as adapterOptions on the model.save method
 */
export default class PkiGenerateRootComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked modelValidations = null;
  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  get defaultFields() {
    return [
      'type',
      'commonName',
      'issuerName',
      'customTtl',
      'notBeforeDuration',
      'format',
      'permittedDnsDomains',
      'maxPathLength',
    ];
  }

  get returnedFields() {
    return [
      'certificate',
      'expiration',
      'issuerId',
      'issuerName',
      'issuingCa',
      'keyId',
      'keyName',
      'serialNumber',
    ];
  }

  @action cancel() {
    // Generate root form will always have a new model
    this.args.model.unloadRecord();
    this.args.onCancel();
  }

  @action
  checkFormValidity() {
    if (this.args.model.validate) {
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return isValid;
    }
    return true;
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    const continueSave = this.checkFormValidity();
    if (!continueSave) return;
    try {
      yield this.setUrls();
      yield this.args.model.save({ adapterOptions: this.args.adapterOptions });
      this.flashMessages.success('Successfully generated root.');
      if (this.args.onSuccess) {
        this.args.onSuccess();
      }
    } catch (e) {
      this.errorBanner = errorMessage(e);
      this.invalidFormAlert = 'There was a problem generating the root.';
    }
  }

  async setUrls() {
    if (!this.args.urls || !this.args.urls.canSet || !this.args.urls.hasDirtyAttributes) return;
    return this.args.urls.save();
  }
}
