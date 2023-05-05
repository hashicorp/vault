/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';
import { ModelValidations } from 'vault/app-types';

export default class PkiKeyModel extends Model {
  secretMountPath: class;
  keyId: string;
  keyName: string;
  type: string;
  keyType: string;
  keyBits: string;
  pemBundle: string;
  privateKey: string;
  isNew: boolean;
  get backend(): string;
  get canRead(): boolean;
  get canEdit(): boolean;
  get canDelete(): boolean;
  get canGenerateKey(): boolean;
  get canImportKey(): boolean;
  validate(): ModelValidations;
}
