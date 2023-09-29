/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';

import type { ModelValidations, FormField, FormFieldGroups } from 'vault/app-types';
import type MountConfigModel from 'vault/models/mount-config';

export default class SecretEngineModel extends Model {
  path: string;
  type: string;
  description: string;
  config: MountConfigModel;
  local: boolean;
  sealWrap: boolean;
  externalEntropyAccess: boolean;
  version: number;
  privateKey: string;
  publicKey: string;
  generateSigningKey: boolean;
  lease: string;
  leaseMax: string;
  accessor: string;
  maxVersions: number;
  casRequired: boolean;
  deleteVersionAfter: string;
  get modelTypeForKV(): string;
  get isV2KV(): boolean;
  get attrs(): Array<FormField>;
  get fieldGroups(): FormFieldGroups;
  get icon(): string;
  get engineType(): string;
  get shouldIncludeInList(): boolean;
  get isSupportedBackend(): boolean;
  get backendLink(): string;
  get accessor(): string;
  get localDisplay(): string;
  get formFields(): Array<FormField>;
  get formFieldGroups(): FormFieldGroups;
  saveCA(options: object): Promise;
  saveZeroAddressConfig(): Promise;
  validate(): ModelValidations;
  // need to override isNew which is a computed prop and ts will complain since it sees it as a function
  isNew: boolean;
}
