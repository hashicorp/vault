/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';

export default class IdentityTokenModel extends Model {
  issuer: string;
  get attrs(): any;
  // for some reason the following Model attrs don't exist on the Model definition
  changedAttributes(): {
    [key: string]: unknown[];
  };
  isNew: boolean;
  canRead: boolean;
  save(): void;
  unloadRecord(): void;
}
