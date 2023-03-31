/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import FlashMessageService from 'ember-cli-flash/services/flash-messages';
import PkiActionModel from 'vault/models/pki/action';
import errorMessage from 'vault/utils/error-message';

interface Args {
  model: PkiActionModel;
  useIssuer: boolean;
  onComplete: CallableFunction;
  onCancel: CallableFunction;
  onSave?: CallableFunction;
}

/**
 * @module PkiGenerateCsrComponent
 * PkiGenerateCsr shows only the fields valid for the generate CSR endpoint.
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
export default class PkiGenerateCsrComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked modelValidations = null;
  @tracked error: string | null = null;
  @tracked alert: string | null = null;

  formFields;
  // fields rendered after CSR generation
  showFields = ['csr', 'keyId', 'privateKey', 'privateKeyType'];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.formFields = expandAttributeMeta(this.args.model, [
      'type',
      'commonName',
      'excludeCnFromSans',
      'format',
      'serialNumber',
      'addBasicConstraints',
    ]);
  }

  @action
  cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }

  async getCapability(): Promise<boolean> {
    try {
      const issuerCapabilities = await this.args.model.generateIssuerCsrPath;
      return issuerCapabilities.get('canCreate') === true;
    } catch (error) {
      return false;
    }
  }

  @task
  @waitFor
  *save(event: Event): Generator<Promise<boolean | PkiActionModel>> {
    event.preventDefault();
    try {
      const { model, onSave } = this.args;
      const { isValid, state, invalidFormMessage } = model.validate();
      if (isValid) {
        const useIssuer = yield this.getCapability();
        yield model.save({ adapterOptions: { actionType: 'generate-csr', useIssuer } });
        this.flashMessages.success('Successfully generated CSR.');
        if (onSave) {
          onSave();
        }
      } else {
        this.modelValidations = state;
        this.alert = invalidFormMessage;
      }
    } catch (e) {
      this.error = errorMessage(e);
      this.alert = 'There was a problem generating the CSR.';
    }
  }
}
