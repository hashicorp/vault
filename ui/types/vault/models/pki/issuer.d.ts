/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';
import { FormField, FormFieldGroups, ModelValidations } from 'vault/app-types';
import { ParsedCertificateData } from 'vault/vault/utils/parse-pki-cert';
export default class PkiIssuerModel extends Model {
  secretMountPath: class;
  get useOpenAPI(): boolean;
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
  /** these are all instances of the capabilities model which should be converted to native class and typed
  rotateExported: any;
  rotateInternal: any;
  rotateExisting: any;
  crossSignPath: any;
  signIntermediate: any;
  -------------------- **/
  pemBundle: string;
  importedIssuers: string[];
  importedKeys: string[];
  formFields: FormField[];
  formFieldGroups: FormFieldGroups[];
  allFields: FormField[];
  get canRotateIssuer(): boolean;
  get canCrossSign(): boolean;
  get canSignIntermediate(): boolean;
  get canConfigure(): boolean;
  validate(): ModelValidations;
}
