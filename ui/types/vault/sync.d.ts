/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type AwsSmForm from 'vault/forms/sync/aws-sm';
import type AzureKvForm from 'vault/forms/sync/azure-kv';
import type GcpSmForm from 'vault/forms/sync/gcp-sm';
import type GhForm from 'vault/forms/sync/gh';
import type VercelProjectForm from 'vault/forms/sync/vercel-project';

export type ListDestination = {
  id: string;
  name: string;
  type: DestinationType;
  icon?: string;
  typeDisplayName?: string;
};

export type AssociatedSecret = {
  mount: string;
  secretName: string;
  syncStatus: string;
  updatedAt: Date;
  destinationType: DestinationType;
  destinationName: string;
};

export type AssociatedDestination = {
  type: string;
  name: string;
  syncStatus: string;
  updatedAt: Date;
};

export type SyncStatus = {
  destinationType: string;
  destinationName: string;
  syncStatus: string;
  updatedAt: string;
};

export type DestinationMetrics = {
  icon?: string;
  name?: string;
  type?: DestinationType;
  associationCount: number;
  status: string | null;
  lastUpdated?: Date;
};

export type AssociationMetrics = {
  totalAssociations: number;
  totalSecrets: number;
};

export type DestinationType = 'aws-sm' | 'azure-kv' | 'gcp-sm' | 'gh' | 'vercel-project';

export type DestinationName =
  | 'AWS Secrets Manager'
  | 'Azure Key Vault'
  | 'Google Secret Manager'
  | 'Github Actions'
  | 'Vercel Project';

export type Destination = {
  name: string;
  type: DestinationType;
  connectionDetails: DestinationConnectionDetails;
  options: DestinationOptions;
  // only present if delete action has been initiated
  purgeInitiatedAt?: string;
  purgeError?: string;
};

export type DestinationConnectionDetails = {
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
};

export type DestinationOptions = {
  granularity?: string; // expected as granularity in request
  granularityLevel?: string; // returned as granularityLevel from response
  secretNameTemplate: string;
  customTags?: Record<string, string>;
};

export type DestinationForm = AwsSmForm | AzureKvForm | GcpSmForm | GhForm | VercelProjectForm;
