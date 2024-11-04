/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import type { ModelValidations } from 'vault/vault/app-types';

export default class AwsLeaseConfig extends Model {
  backend: any;
  leaseMax: any;
  lease: any;
  get attrs(): string[];
  // for some reason the following Model attrs don't exist on the Model definition
  changedAttributes(): {
    [key: string]: unknown[];
  };
  isNew: boolean;
  save(): void;
  unloadRecord(): void;
  validate(): ModelValidations;
}
