/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { WithFormFieldsAndValidationsModel } from 'vault/app-types';

export default interface SyncDestinationModel extends WithFormFieldsAndValidationsModel {
  name: string;
  type: string;
  secretNameTemplate: string;
  purgeInitiatedAt: string;
  purgeError: string;

  get icon(): string;
  get typeDisplayName(): string;
  get maskedParams(): string[];

  // aws-sm
  accessKeyId?: string;
  secretAccessKey?: string;
  region?: string;

  // azure-kv
  keyVaultUri?: string;
  clientId?: string;
  clientSecret?: string;
  tenantId?: string;
  cloud?: string;

  // gcp
  credentials?: string;

  // gh
  accessToken?: string;
  repositoryOwner?: string;
  repositoryName?: string;

  // vercel project
  accessToken?: string;
  projectId?: string;
  teamId?: string;
  deploymentEnvironments?: array;

  get canCreate(): boolean;
  get canDelete(): boolean;
  get canEdit(): boolean;
  get canRead(): boolean;
  get canSync(): boolean;
}
