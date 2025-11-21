/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type PkiCertificateBaseModel from './base';

type PkiCertificateGenerateModel = PkiCertificateBaseModel & {
  role: string;
};

export default PkiCertificateGenerateModel;
