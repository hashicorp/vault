/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type PkiCertificateBaseModel from './base';

type PkiCertificateSignIntermediateModel = PkiCertificateBaseModel & {
  role: string;
  csr: string;
  issuerRef: string;
  maxPathLength: string;
  notBeforeDuration: string;
  permittedDnsDomains: string;
  useCsrValues: boolean;
  usePss: boolean;
  skid: string;
  signatureBits: string;
};

export default PkiCertificateSignIntermediateModel;
