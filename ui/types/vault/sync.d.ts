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
  type_display_name?: string;
};

export type AssociatedSecret = {
  mount: string;
  secret_name: string;
  sync_status: string;
  updated_at: Date;
  destination_type: DestinationType;
  destination_name: string;
};

export type AssociatedDestination = {
  type: string;
  name: string;
  sync_status: string;
  updated_at: Date;
};

export type SyncStatus = {
  destination_type: string;
  destination_name: string;
  sync_status: string;
  updated_at: string;
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
  total_associations: number;
  total_secrets: number;
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
  connection_details: DestinationConnectionDetails;
  options: DestinationOptions;
  // only present if delete action has been initiated
  purge_initiated_at?: string;
  purge_error?: string;
};

export type DestinationConnectionDetails = {
  // aws-sm
  access_key_id?: string;
  secret_access_key?: string;
  region?: string;
  // azure-kv
  key_vault_uri?: string;
  client_id?: string;
  client_secret?: string;
  tenant_id?: string;
  cloud?: string;
  // gcp
  credentials?: string;
  // gh
  access_token?: string;
  repository_owner?: string;
  repository_name?: string;
  // vercel project
  access_token?: string;
  project_id?: string;
  team_id?: string;
  deployment_environments?: array;
};

export type DestinationOptions = {
  granularity?: string; // expected as granularity in request
  granularity_level?: string; // returned as granularity_level from response
  secret_name_template: string;
  custom_tags?: Record<string, string>;
};

export type DestinationForm = AwsSmForm | AzureKvForm | GcpSmForm | GhForm | VercelProjectForm;
