/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type PkiCertificateBaseModel from './base';

export type PkiCertificateSignModel = PkiCertificateBaseModel & {
  role: string;
  csr: string;
  removeRootsFromChain: boolean;
};

export default PkiCertificateSignModel;
