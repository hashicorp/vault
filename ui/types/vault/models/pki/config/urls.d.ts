/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Model } from 'vault/app-types';

type PkiConfigUrlsModel = Model & {
  issuingCertificates: array;
  crlDistributionPoints: array;
  ocspServers: array;
  urlsPath: string;
  get canSet(): boolean;
};

export default PkiConfigUrlsModel;
