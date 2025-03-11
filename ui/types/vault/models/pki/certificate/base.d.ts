/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';

type PkiCertificateBaseModel = WithFormFieldsAndValidationsModel & {
  secretMountPath: class;
  get backend(): string;
  altNames: string;
  commonName: string;
  caChain: string;
  certificate: string;
  excludeCnFromSans: boolean;
  expiration: number;
  ipSans: string;
  issuingCa: string;
  notValidAfter: date;
  notValidBefore: date;
  otherSans: string;
  privateKey: string;
  privateKeyType: string;
  revokePath: string;
  revocationTime: number;
  serialNumber: string;
  get canRevoke(): boolean;
};

export default PkiCertificateBaseModel;
