/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type Model from '@ember-data/model';
import type { ModelValidations } from 'vault/app-types';

export default class GcpConfig extends Model {
  backend: string;
  credentials: string | undefined;
  ttl: any;
  maxTtl: any;
  secretAccountEmail: string | undefined;
  identityTokenAudience: string | undefined;
  identityTokenTtl: any;

  get displayAttrs(): any;
  get isWifPluginConfigured(): boolean;
  get fieldGroupsWif(): any;
  get fieldGroupsGcp(): any;
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
