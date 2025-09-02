/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';

export default class IdentityOidcConfigModel extends Model {
  issuer: string;
  queryIssuerError: boolean;
  get attrs(): any;
  // for some reason the following Model attrs don't exist on the Model definition
  changedAttributes(): {
    [key: string]: unknown[];
  };
  rollbackAttributes(): { void };
  hasDirtyAttributes: boolean;
  isNew: boolean;
  canRead: boolean;
  save(): void;
  unloadRecord(): void;
}
