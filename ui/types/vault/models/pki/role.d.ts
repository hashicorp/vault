/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import { ModelValidations } from 'vault/app-types';

export default class PkiRoleModel extends Model {
  name: string;
  issuerRef: string;
  validate(): ModelValidations;
  isNew: boolean;
  keyType: string;
  keyBits: string | undefined;
}
