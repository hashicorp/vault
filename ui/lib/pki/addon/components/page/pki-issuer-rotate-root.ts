/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { waitFor } from '@ember/test-waiters';
import { task } from 'ember-concurrency';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import PkiConfigGenerateForm from 'vault/forms/secrets/pki/config/generate';
import { toLabel } from 'core/helpers/to-label';

import type RouterService from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { Breadcrumb, ValidationMap } from 'vault/vault/app-types';
import type {
  PkiGenerateRootResponse,
  PkiReadIssuerResponse,
  SecretsApiPkiRotateRootExportedEnum,
  PkiRotateRootRequest,
} from '@hashicorp/vault-client-typescript';
import type { ParsedCertificateData } from 'vault/vault/utils/parse-pki-cert';
import type ApiService from 'vault/services/api';

interface Args {
  oldRoot: PkiReadIssuerResponse;
  certData: ParsedCertificateData;
  breadcrumbs: Breadcrumb;
  parsingErrors: string;
}

const RADIO_BUTTON_KEY = {
  oldSettings: 'use-old-settings',
  customizeNew: 'customize',
};

export default class PagePkiIssuerRotateRootComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;

  @tracked displayedForm = RADIO_BUTTON_KEY.oldSettings;
  @tracked showOldSettings = false;
  @tracked modelValidations: ValidationMap | null = null;
  @tracked declare newRoot: PkiGenerateRootResponse;
  // form alerts below are only for "use old settings" option
  // validations/errors for "customize new root" are handled by <PkiGenerateRoot> component
  @tracked alertBanner = '';
  @tracked invalidFormAlert = '';

  newRootForm = new PkiConfigGenerateForm(
    'PkiGenerateRootRequest',
    { type: 'internal', ...this.args.certData },
    { isNew: true }
  );

  generateOptions = [
    {
      key: RADIO_BUTTON_KEY.oldSettings,
      icon: 'certificate',
      label: 'Use old root settings',
      description: `Provide only a new common name and issuer name, using the old root’s settings. Selecting this option generates a root with Vault-internal key material.`,
    },
    {
      key: RADIO_BUTTON_KEY.customizeNew,
      icon: 'award',
      label: 'Customize new root certificate',
      description:
        'Generates a new self-signed CA certificate and private key. This generated root will sign its own CRL.',
    },
  ];

  label = (field: string) => {
    return (
      {
        ca_chain: 'CA Chain',
        issuer_id: 'Issuer ID',
        key_id: 'Default key ID',
      }[field] || toLabel([field])
    );
  };

  // for displaying old root details, and generated root details
  get displayFields() {
    const addKeyFields = ['private_key', 'private_key_type'];
    const defaultFields = [
      'certificate',
      'ca_chain',
      'issuer_id',
      'issuer_name',
      'issuing_ca',
      'key_name',
      'key_id',
      'serial_number',
    ];
    return this.newRoot ? [...defaultFields, ...addKeyFields] : defaultFields;
  }

  get newParsedCertificate() {
    return parseCertificate(this.newRoot?.certificate || '');
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      const { isValid, state, invalidFormMessage, data } = this.newRootForm.toJSON();
      if (isValid) {
        const { type } = this.newRootForm.data;
        try {
          this.newRoot = await this.api.secrets.pkiRotateRoot(
            type as SecretsApiPkiRotateRootExportedEnum,
            this.secretMountPath.currentPath,
            data as PkiRotateRootRequest
          );
          this.flashMessages.success('Successfully generated root.');
        } catch (e) {
          const { message } = await this.api.parseError(e);
          this.alertBanner = message;
          this.invalidFormAlert = 'There was a problem generating root.';
        }
      } else {
        this.modelValidations = state;
        this.invalidFormAlert = invalidFormMessage;
      }
    })
  );

  @action
  async fetchDataForDownload(format: 'der' | 'pem') {
    try {
      const path = `/${this.secretMountPath.currentPath}/issuer/${this.newRoot.issuer_id}/${format}`;
      const response = await this.api.request.get(path);
      const body = format === 'der' ? 'blob' : 'text';
      return response[body]();
    } catch (e) {
      return null;
    }
  }
}
