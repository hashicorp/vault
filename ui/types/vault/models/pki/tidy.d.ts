/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import { FormField, FormFieldGroups } from 'vault/vault/app-types';

export default class PkiTidyModel extends Model {
  version: string;
  acmeAccountSafetyBuffer: string;
  tidyAcme: boolean;
  enabled: boolean;
  intervalDuration: string;
  issuerSafetyBuffer: string;
  pauseDuration: string;
  revocationQueueSafetyBuffer: string;
  safetyBuffer: string;
  tidyCertMetadata: boolean;
  tidyCertStore: boolean;
  tidyCrossClusterRevokedCerts: boolean;
  tidyExpiredIssuers: boolean;
  tidyMoveLegacyCaBundle: boolean;
  tidyRevocationQueue: boolean;
  tidyRevokedCertIssuerAssociations: boolean;
  tidyRevokedCerts: boolean;
  get useOpenAPI(): boolean;
  getHelpUrl(backend: string): string;
  allByKey: {
    intervalDuration: FormField[];
  };
  get allGroups(): FormFieldGroups[];
  get sharedFields(): FormFieldGroups[];
  get formFieldGroups(): FormFieldGroups[];
}
