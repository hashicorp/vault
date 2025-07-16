/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { StringMap, WithFormFieldsAndValidationsModel } from 'vault/app-types';
import type CapabilitiesModel from '../capabilities';

type PkiActionModel = WithFormFieldsAndValidationsModel & {
  secretMountPath: unknown;
  actionType: string | null;
  pemBundle: string;
  importedIssuers: string[];
  importedKeys: string[];
  mapping: StringMap;
  type: string;
  issuerName: string;
  keyName: string;
  keyRef: string;
  commonName: string;
  altNames: string[];
  ipSans: string[];
  uriSans: string[];
  otherSans: string[];
  format: string;
  privateKeyFormat: string;
  keyType: string;
  keyBits: string;
  maxPathLength: number;
  excludeCnFromSans: boolean;
  permittedDnsDomains: string;
  ou: string[];
  serialNumber: string;
  addBasicConstraints: boolean;
  notBeforeDuration: string;
  managedKeyName: string;
  managedKeyId: string;
  customTtl: string;
  ttl: string;
  notAfter: string;
  issuerId: string;
  csr: string;
  caChain: string;
  keyId: string;
  privateKey: string;
  privateKeyType: string;
  get backend(): string;
  // apiPaths for capabilities
  importBundlePath: Promise<CapabilitiesModel>;
  generateIssuerRootPath: Promise<CapabilitiesModel>;
  generateIssuerCsrPath: Promise<CapabilitiesModel>;
  crossSignPath: string;
  // Capabilities
  get canImportBundle(): boolean;
  get canGenerateIssuerRoot(): boolean;
  get canGenerateIssuerIntermediate(): boolean;
  get canCrossSign(): boolean;
};

export default PkiActionModel;
