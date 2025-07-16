/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import type PkiActionModel from 'vault/models/pki/action';
import type PkiConfigUrlsModel from 'vault/models/pki/config/urls';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { ValidationMap } from 'vault/vault/app-types';

interface AdapterOptions {
  actionType: string;
  useIssuer: boolean | undefined;
}
interface Args {
  model: PkiActionModel;
  urls: PkiConfigUrlsModel;
  onCancel: CallableFunction;
  onComplete: CallableFunction;
  onSave?: CallableFunction;
  adapterOptions: AdapterOptions;
  hideAlertBanner: boolean;
}

/**
 * @module PkiGenerateRoot
 * PkiGenerateRoot shows only the fields valid for the generate root endpoint.
 * This component renders the form, handles the model save and rollback actions,
 * and shows the resulting data on success. onCancel is required for the cancel
 * transition, and if onSave is provided it will call that after save for any
 * side effects in the parent.
 *
 * @example
 * ```js
 * <PkiGenerateRoot @model={{this.model}} @onCancel={{transition-to "vault.cluster"}} @onSave={{fn (mut this.title) "Successful"}} @adapterOptions={{hash actionType="import" useIssuer=false}} />
 * ```
 *
 * @param {Object} model - pki/action model.
 * @callback onCancel - Callback triggered when cancel button is clicked, after model is unloaded
 * @callback onSave - Optional - Callback triggered after model is saved, as a side effect. Results are shown on the same component
 * @callback onComplete - Callback triggered when "Done" button clicked, on results view
 * @param {Object} adapterOptions - object passed as adapterOptions on the model.save method
 */
export default class PkiGenerateRootComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;

  @tracked modelValidations: ValidationMap | null = null;
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
      'commonName',
      'issuerId',
      'issuerName',
      'issuingCa',
      'keyName',
      'keyId',
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
      yield this.args.model.save({ adapterOptions: this.args.adapterOptions });
      // root generation must occur first in case templates are used for URL fields
      // this way an issuer_id exists for backend to interpolate into the template
      yield this.setUrls();
      this.flashMessages.success('Successfully generated root.');
      // This component shows the results, but call `onSave` for any side effects on parent
      if (this.args.onSave) {
        this.args.onSave();
      }
      window?.scrollTo(0, 0);
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
