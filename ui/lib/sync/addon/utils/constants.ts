/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export enum CredentialType {
  ACCOUNT = 'account',
  WIF = 'wif',
}

export enum DestinationType {
  AwsSm = 'aws-sm',
  AzureKv = 'azure-kv',
  GcpSm = 'gcp-sm',
  Gh = 'gh',
  VercelProject = 'vercel-project',
}

export const CLOUD_DESTINATION_TYPES = [
  DestinationType.AwsSm,
  DestinationType.AzureKv,
  DestinationType.GcpSm,
] as const;

export type CloudDestinationType = (typeof CLOUD_DESTINATION_TYPES)[number];

const COMMON_WIF_FIELDS = ['identity_token_audience', 'identity_token_ttl', 'identity_token_key'];

export const ACCOUNT_CREDENTIAL_FIELDS: Record<CloudDestinationType, string[]> = {
  [DestinationType.AwsSm]: ['access_key_id', 'secret_access_key'],
  [DestinationType.AzureKv]: ['client_secret'],
  [DestinationType.GcpSm]: ['credentials'],
};

export const WIF_CREDENTIAL_FIELDS: Record<CloudDestinationType, string[]> = {
  [DestinationType.AwsSm]: [...COMMON_WIF_FIELDS],
  [DestinationType.AzureKv]: [...COMMON_WIF_FIELDS],
  [DestinationType.GcpSm]: [...COMMON_WIF_FIELDS, 'service_account_email'],
};
