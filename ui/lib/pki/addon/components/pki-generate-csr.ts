/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { toLabel } from 'core/helpers/to-label';
import PkiConfigGenerateForm from 'vault/forms/secrets/pki/config/generate';

import type FlashMessageService from 'vault/services/flash-messages';
import type CapabilitiesService from 'vault/services/capabilities';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { ValidationMap } from 'vault/app-types';
import type {
  SecretsApiPkiGenerateIntermediateExportedEnum,
  SecretsApiPkiIssuersGenerateIntermediateExportedEnum,
  PkiGenerateIntermediateRequest,
  PkiGenerateIntermediateResponse,
  PkiIssuersGenerateIntermediateRequest,
  PkiIssuersGenerateIntermediateResponse,
} from '@hashicorp/vault-client-typescript';

interface Args {
  form: PkiConfigGenerateForm;
  onComplete: CallableFunction;
  onCancel: CallableFunction;
  onSave?: CallableFunction;
}

/**
 * @module PkiGenerateCsrComponent
 * PkiGenerateCsr shows only the fields valid for the generate CSR endpoint.
 * This component renders the form, handles saving the data to the server,
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
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked error: string | null = null;
  @tracked alert: string | null = null;
  @tracked declare config: PkiGenerateIntermediateResponse | PkiIssuersGenerateIntermediateResponse;

  form = new PkiConfigGenerateForm('PkiGenerateIntermediateRequest', {}, { isNew: true });

  defaultFields = [
    'type',
    'common_name',
    'exclude_cn_from_sans',
    'format',
    'serial_number',
    'add_basic_constraints',
  ];
  // fields rendered after CSR generation
  returnedFields = ['csr', 'key_id', 'private_key', 'private_key_type'];

  detailLabel = (fieldName: string) => {
    return (
      {
        csr: 'CSR',
        key_id: 'Key ID',
      }[fieldName] || toLabel([fieldName])
    );
  };

  async fetchIssuerCapabilities() {
    try {
      const { canCreate } = await this.capabilities.for('pkiIssuersGenerateIntermediate', {
        backend: this.secretMountPath.currentPath,
        type: this.form.data.type,
      });
      return canCreate;
    } catch (e) {
      // fallback to pkiGenerateIntermediate if capabilities fetch fails
      return false;
    }
  }

  generateCsr(canUseIssuer: boolean, data: PkiConfigGenerateForm['data']) {
    if (canUseIssuer) {
      return this.api.secrets.pkiIssuersGenerateIntermediate(
        this.form.data.type as SecretsApiPkiIssuersGenerateIntermediateExportedEnum,
        this.secretMountPath.currentPath,
        data as PkiIssuersGenerateIntermediateRequest
      );
    } else {
      return this.api.secrets.pkiGenerateIntermediate(
        this.form.data.type as SecretsApiPkiGenerateIntermediateExportedEnum,
        this.secretMountPath.currentPath,
        data as PkiGenerateIntermediateRequest
      );
    }
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      try {
        const { isValid, state, invalidFormMessage, data } = this.form.toJSON();
        if (isValid) {
          const canUseIssuer = await this.fetchIssuerCapabilities();
          this.config = await this.generateCsr(canUseIssuer, data);
          this.flashMessages.success('Successfully generated CSR.');
          // This component shows the results, but call `onSave` for any side effects on parent
          this.args.onSave?.();
          window?.scrollTo(0, 0);
        } else {
          this.modelValidations = state;
          this.alert = invalidFormMessage;
        }
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.error = message;
        this.alert = 'There was a problem generating the CSR.';
      }
    })
  );
}
