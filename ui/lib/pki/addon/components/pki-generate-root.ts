/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { toLabel } from 'core/helpers/to-label';
import { parseCertificate, type ParsedCertificateData } from 'vault/utils/parse-pki-cert';
import PkiConfigGenerateForm from 'vault/forms/secrets/pki/config/generate';
import PkiUrlsForm from 'vault/forms/secrets/pki/config/urls';

import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { ValidationMap } from 'vault/app-types';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type {
  SecretsApiPkiGenerateRootExportedEnum,
  SecretsApiPkiIssuersGenerateRootExportedEnum,
  PkiGenerateRootRequest,
  PkiGenerateRootResponse,
  PkiIssuersGenerateRootRequest,
  PkiIssuersGenerateRootResponse,
  SecretsApiPkiRotateRootExportedEnum,
  PkiRotateRootRequest,
} from '@hashicorp/vault-client-typescript';

interface Args {
  withUrls?: boolean;
  canSetUrls?: boolean;
  rotateCertData?: ParsedCertificateData;
  onCancel: CallableFunction;
  onComplete: CallableFunction;
  onSave?: CallableFunction;
}

/**
 * @module PkiGenerateRoot
 * PkiGenerateRoot shows only the fields valid for the generate root endpoint.
 * This component renders the form, handles saving the data to the server,
 * and shows the resulting data on success. onCancel is required for the cancel
 * transition, and if onSave is provided it will call that after save for any
 * side effects in the parent.
 *
 * @param {boolean} withUrls - whether or not to show the urls fields.
 * @param {boolean} canSetUrls - whether or not the user has capability to set urls.
 * @param {boolean} rotateCertData - cert data to be populated in form and rotated.
 * @callback onCancel - Callback triggered when cancel button is clicked, after model is unloaded
 * @callback onSave - Optional - Callback triggered after model is saved, as a side effect. Results are shown on the same component
 * @callback onComplete - Callback triggered when "Done" button clicked, on results view
 */
export default class PkiGenerateRootComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked declare config: PkiGenerateRootResponse | PkiIssuersGenerateRootResponse;
  @tracked declare parsedCertificate: ParsedCertificateData;

  form = new PkiConfigGenerateForm('PkiGenerateRootRequest', this.args.rotateCertData, {
    isNew: true,
  });
  urlsForm = new PkiUrlsForm({}, { isNew: true });

  get defaultFields() {
    return [
      'type',
      'common_name',
      'issuer_name',
      'customTtl', // UI only and used to yield PkiNotValidAfterForm
      'not_before_duration',
      'format',
      'permitted_dns_domains',
      'max_path_length',
    ];
  }

  get returnedFields() {
    return [
      'certificate',
      'common_name',
      'issuer_id',
      'issuer_name',
      'issuing_ca',
      'key_name',
      'key_id',
      'serial_number',
    ];
  }

  detailLabel = (fieldName: string) => {
    const label = toLabel([fieldName]);
    return (
      {
        issuer_id: 'Issuer ID',
        issuing_ca: 'Issuing CA',
        key_id: 'Key ID',
      }[fieldName] || label
    );
  };

  linkForField = (fieldName: string) => {
    return {
      issuer_id: 'issuers.issuer.details',
      key_id: 'keys.key.details',
    }[fieldName];
  };

  valueForField = (fieldName: string) => {
    if (fieldName === 'common_name') {
      return this.parsedCertificate.common_name;
    }
    return this.config[fieldName as keyof typeof this.config];
  };

  isCertificateField = (fieldName: string) => {
    return ['certificate', 'issuing_ca', 'csr', 'private_key'].includes(fieldName);
  };

  async fetchIssuerCapabilities() {
    try {
      const { canCreate } = await this.capabilities.for('pkiIssuersGenerateRoot', {
        backend: this.secretMountPath.currentPath,
        type: this.form.data.type,
      });
      return canCreate;
    } catch (e) {
      // fallback to pkiGenerateRoot if capabilities fetch fails
      return false;
    }
  }

  async generateRoot(data: PkiConfigGenerateForm['data']) {
    const { type } = this.form.data;
    const { currentPath } = this.secretMountPath;

    if (this.args.rotateCertData) {
      return this.api.secrets.pkiRotateRoot(
        type as SecretsApiPkiRotateRootExportedEnum,
        currentPath,
        data as PkiRotateRootRequest
      );
    } else {
      const canUseIssuer = await this.fetchIssuerCapabilities();
      if (canUseIssuer) {
        return this.api.secrets.pkiIssuersGenerateRoot(
          type as SecretsApiPkiIssuersGenerateRootExportedEnum,
          this.secretMountPath.currentPath,
          data as PkiIssuersGenerateRootRequest
        );
      } else {
        return this.api.secrets.pkiGenerateRoot(
          type as SecretsApiPkiGenerateRootExportedEnum,
          this.secretMountPath.currentPath,
          data as PkiGenerateRootRequest
        );
      }
    }
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      const { isValid, state, invalidFormMessage, data } = this.form.toJSON();

      if (isValid) {
        try {
          this.config = await this.generateRoot(data);
          this.parsedCertificate = parseCertificate(this.config.certificate || '');
          // root generation must occur first in case templates are used for URL fields
          // this way an issuer_id exists for backend to interpolate into the template
          await this.setUrls();
          this.flashMessages.success('Successfully generated root.');
          // This component shows the results, but call `onSave` for any side effects on parent
          this.args.onSave?.(this.config);
          window?.scrollTo(0, 0);
        } catch (e) {
          const { message } = await this.api.parseError(e);
          this.errorBanner = message;
          this.invalidFormAlert = 'There was a problem generating the root.';
        }
      } else {
        this.modelValidations = state;
        this.invalidFormAlert = invalidFormMessage;
      }
    })
  );

  async setUrls() {
    const { withUrls, canSetUrls } = this.args;
    if (withUrls && canSetUrls) {
      const { data } = this.urlsForm.toJSON();
      await this.api.secrets.pkiConfigureUrls(this.secretMountPath.currentPath, data);
    }
  }
}
