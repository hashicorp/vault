/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';
import type { ModelValidations } from 'vault/app-types';

export default class AzureConfig extends Model {
  backend: string;
  subscriptionId: string | undefined;
  tenantId: string | undefined;
  clientId: string | undefined;
  clientSecret: string | undefined;
  identityTokenAudience: string | undefined;
  identityTokenTtl: any;
  environment: string | undefined;
  rootPasswordTtl: string | undefined;

  get displayAttrs(): any;
  get isWifPluginConfigured(): boolean;
  get isAzureAccountConfigured(): boolean;
  get fieldGroupsWif(): any;
  get fieldGroupsAzure(): any;
  formFieldGroups(accessType?: string): {
    [key: string]: string[];
  }[];
  changedAttributes(): {
    [key: string]: unknown[];
  };
  isNew: boolean;
  save(): void;
  unloadRecord(): void;
}
