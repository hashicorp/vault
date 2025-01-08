/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';
import type { ModelValidations } from 'vault/app-types';

export default class SecretEngineConfig extends Model {
  backend: string;
  type: string;
  // aws
  leaseMax: any;
  lease: any;
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
  displayTitle: string | undefined;
  // wif
  identityTokenAudience: string | undefined;
  identityTokenTtl: any;

  get displayAttrs(): any;
  get isWifPluginConfigured(): boolean;
  get isAccountPluginConfigured(): boolean;
  get fieldGroupsWif(): any;
  get fieldGroupsAzure(): any;
  get fieldGroupsGcp(): any;
  get fieldGroupsIam(): any;
  formFieldGroups(accessType?: string): {
    [key: string]: string[];
  }[];
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
