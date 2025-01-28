/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';
import type { ModelValidations, FormFieldGroups } from 'vault/app-types';

export default class SecretEngineMountConfigModel extends Model {
  backend: string;
  type: string;
  // aws
  accessKey: any;
  secretKey: any;
  roleArn: any;
  region: any;
  iamEndpoint: any;
  stsEndpoint: any;
  maxRetries: any;
  // azure
  subscriptionId: string | undefined;
  tenantId: string | undefined;
  clientId: string | undefined;
  clientSecret: string | undefined;
  environment: string | undefined;
  rootPasswordTtl: string | undefined;
  // gcp
  credentials: string | undefined;
  ttl: any;
  maxTtl: any;
  secretAccountEmail: string | undefined;
  // wif
  identityTokenAudience: string | undefined;
  identityTokenTtl: any;

  get displayAttrs(): any;
  get isConfigured(): boolean; // used only for secret engines that return a 200 when configuration has not be set.
  get isWifPluginConfigured(): boolean;
  get isAccountPluginConfigured(): boolean;
  get fieldGroupsWif(): any;
  get fieldGroupsAccount(): any;

  formFieldGroups: FormFieldGroups[];

  changedAttributes(): {
    [key: string]: unknown[];
  };
  isNew: boolean;
  save(): void;
  unloadRecord(): void;
  destroyRecord(): void;
  rollbackAttributes(): void;
}
