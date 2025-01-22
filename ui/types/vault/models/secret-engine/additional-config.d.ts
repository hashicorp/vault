/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';
import type { ModelValidations, FormField } from 'vault/app-types';

export default class SecretEngineAdditionalConfigModel extends Model {
  backend: string;
  type: string;
  // aws lease
  leaseMax: any;
  lease: any;

  get displayAttrs(): any;

  formFields: FormField[];
  changedAttributes(): {
    [key: string]: unknown[];
  };
  isNew: boolean;
  save(): void;
  unloadRecord(): void;
  destroyRecord(): void;
  rollbackAttributes(): void;
  validate(): ModelValidations;
}
