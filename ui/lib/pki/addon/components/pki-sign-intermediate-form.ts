/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { parseCertificate } from 'vault/utils/parse-pki-cert';

import type PkiIssuersSignIntermediateForm from 'vault/forms/secrets/pki/issuers/sign-intermediate';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretMountPathService from 'vault/services/secret-mount-path';
import type { ValidationMap } from 'vault/vault/app-types';
import type { PkiIssuerSignIntermediateResponse } from '@hashicorp/vault-client-typescript';
import type { ParsedCertificateData } from 'vault/vault/utils/parse-pki-cert';

interface Args {
  onCancel: CallableFunction;
  form: PkiIssuersSignIntermediateForm;
  issuerRef: string;
}

export default class PkiSignIntermediateFormComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPathService;

  @tracked errorBanner = '';
  @tracked inlineFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;
  @tracked declare signedIntermediate: PkiIssuerSignIntermediateResponse;
  @tracked declare parsedCertificate: ParsedCertificateData;

  fields = [
    'csr',
    'use_csr_values',
    'common_name',
    'exclude_cn_from_sans',
    'customTtl',
    'not_before_duration',
    'enforce_leaf_not_after_behavior',
    'format',
    'max_path_length',
  ];

  groups = {
    'Name constraints': [
      'permitted_dns_domains',
      'permitted_email_addresses',
      'permitted_ip_ranges',
      'permitted_uri_domains',
      'excluded_dns_domains',
      'excluded_email_addresses',
      'excluded_ip_ranges',
      'excluded_uri_domains',
    ],
    'Signing options': ['use_pss', 'skid', 'signature_bits'],
    'Subject Alternative Name (SAN) Options': ['alt_names', 'ip_sans', 'uri_sans', 'other_sans'],
    'Additional subject fields': [
      'ou',
      'organization',
      'country',
      'locality',
      'province',
      'street_address',
      'postal_code',
      'serial_number',
    ],
  };

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      const { issuerRef, form } = this.args;
      const { isValid, state, invalidFormMessage, data } = form.toJSON();
      this.modelValidations = isValid ? null : state;
      this.inlineFormAlert = invalidFormMessage;
      if (isValid) {
        try {
          this.signedIntermediate = await this.api.secrets.pkiIssuerSignIntermediate(
            issuerRef,
            this.secretMountPath.currentPath,
            data
          );
          this.parsedCertificate = parseCertificate(this.signedIntermediate.certificate as string);
          this.flashMessages.success('Successfully signed CSR.');
          window?.scrollTo(0, 0);
        } catch (e) {
          const { message } = await this.api.parseError(e);
          this.errorBanner = message;
          this.inlineFormAlert = 'There was a problem signing the CSR.';
        }
      }
    })
  );
}
