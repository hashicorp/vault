/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';

export default class AwsRootConfig extends Model {
  backend: any;
  accessKey: any;
  secretKey: any;
  roleArn: any;
  identityTokenAudience: any;
  identityTokenTtl: any;
  region: any;
  iamEndpoint: any;
  stsEndpoint: any;
  maxRetries: any;
  get attrs(): any;
  get fieldGroupsWif(): any;
  get fieldGroupsIam(): any;
  formFieldGroups(accessType?: string): {
    [key: string]: string[];
  }[];
  // for some reason the following Model attrs don't exist on the Model definition
  changedAttributes(): {
    [key: string]: unknown[];
  };
  isNew: boolean;
  save(): void;
  unloadRecord(): void;
}
