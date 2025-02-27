/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { FormField, FormFieldGroups, WithFormFieldsModel } from 'vault/app-types';

type PkiTidyModel = WithFormFieldsModel & {
  version: string;
  acmeAccountSafetyBuffer: string;
  tidyAcme: boolean;
  enabled: boolean;
  intervalDuration: string;
  minStartupBackoffDuration: string;
  maxStartupBackoffDuration: string;
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
  allByKey: {
    intervalDuration: FormField[];
  };
  get sharedFields(): FormFieldGroups[];
};

export default PkiTidyModel;
