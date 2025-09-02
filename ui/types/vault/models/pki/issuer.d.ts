/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type { ParsedCertificateData } from 'vault/utils/parse-pki-cert';
import type CapabilitiesModel from 'vault/models/capabilities';

type PkiIssuerModel = WithFormFieldsAndValidationsModel & {
  secretMountPath: class;
  get backend(): string;
  get issuerRef(): string;
  certificate: string;
  issuerId: string;
  issuerName: string;
  keyId: string;
  uriSans: string;
  leafNotAfterBehavior: string;
  usage: string;
  manualChain: string;
  issuingCertificates: string;
  crlDistributionPoints: string;
  ocspServers: string;
  parsedCertificate: ParsedCertificateData;
  rotateExported: CapabilitiesModel;
  rotateInternal: CapabilitiesModel;
  rotateExisting: CapabilitiesModel;
  crossSignPath: CapabilitiesModel;
  signIntermediate: CapabilitiesModel;
  pemBundle: string;
  importedIssuers: string[];
  importedKeys: string[];
  get canRotateIssuer(): boolean;
  get canCrossSign(): boolean;
  get canSignIntermediate(): boolean;
  get canConfigure(): boolean;
};

export default PkiIssuerModel;
