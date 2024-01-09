/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';

export default class PkiConfigUrlsModel extends Model {
  get useOpenAPI(): boolean;
  getHelpUrl(backendPath: string): string;
  issuingCertificates: array;
  crlDistributionPoints: array;
  ocspServers: array;
  urlsPath: string;
  get canSet(): boolean;
}
