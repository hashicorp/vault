/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { DestinationType, CredentialType } from 'sync/utils/constants';
import type { SyncDestination } from 'vault/helpers/sync-destinations';

/* 
This helper is used to lookup static display properties for sync destinations
maskedParams: attributes for sensitive data, the API returns these values as '*****'
*/

const SYNC_DESTINATIONS: Array<SyncDestination> = [
  {
    name: 'AWS Secrets Manager',
    type: DestinationType.AwsSm,
    icon: 'aws-color',
    category: 'cloud',
    maskedParams: ['access_key_id', 'secret_access_key', 'identity_token_audience', 'identity_token_key'],
    readonlyParams: ['name', 'region'],
    defaultValues: {
      granularity: 'secret-path',
      credential_type: CredentialType.ACCOUNT,
    },
    roleTypeOptions: [
      {
        title: 'IAM Credentials',
        description:
          'Use an AWS Access Key ID and Secret Access Key to allow Vault to interact directly with your AWS resources.',
        value: CredentialType.ACCOUNT,
      },
      {
        title: 'Workload Identity Federation',
        description:
          'Leverages OIDC or AWS IAM Roles for Service Accounts (IRSA) for more secure, keyless authentication.',
        value: CredentialType.WIF,
      },
    ],
  },
  {
    name: 'Azure Key Vault',
    type: DestinationType.AzureKv,
    icon: 'azure-color',
    category: 'cloud',
    maskedParams: ['client_secret', 'identity_token_audience', 'identity_token_key'],
    readonlyParams: ['name', 'key_vault_uri', 'tenant_id', 'cloud'],
    defaultValues: {
      granularity: 'secret-path',
      credential_type: CredentialType.ACCOUNT,
    },
    roleTypeOptions: [
      {
        title: 'Client Secret',
        description: 'Use client secret of an Azure app registration to authenticate.',
        value: CredentialType.ACCOUNT,
      },
      {
        title: 'Workload Identity Federation',
        description:
          'Leverages OIDC with Azure workload identity pools and providers for more secure, keyless authentication.',
        value: CredentialType.WIF,
      },
    ],
  },
  {
    name: 'Google Secret Manager',
    type: DestinationType.GcpSm,
    icon: 'gcp-color',
    category: 'cloud',
    maskedParams: ['credentials', 'identity_token_audience', 'identity_token_key'],
    readonlyParams: ['name'],
    defaultValues: {
      granularity: 'secret-path',
      credential_type: CredentialType.ACCOUNT,
    },
    roleTypeOptions: [
      {
        title: 'JSON Credentials',
        description: 'Use a JSON file from your computer to authenticate.',
        value: CredentialType.ACCOUNT,
      },
      {
        title: 'Workload Identity Federation',
        description:
          'Leverages OIDC with GCP workload identity pools and providers for more secure, keyless authentication.',
        value: CredentialType.WIF,
      },
    ],
  },
  {
    name: 'Github Actions',
    type: DestinationType.Gh,
    icon: 'github-color',
    category: 'dev-tools',
    maskedParams: ['access_token'],
    readonlyParams: ['name', 'repository_owner', 'repository_name'],
    defaultValues: {
      granularity: 'secret-key',
    },
  },
  {
    name: 'Vercel Project',
    type: DestinationType.VercelProject,
    icon: 'vercel-color',
    category: 'dev-tools',
    maskedParams: ['access_token'],
    readonlyParams: ['name', 'project_id'],
    defaultValues: {
      granularity: 'secret-key',
      deployment_environments: [],
    },
  },
];

export function syncDestinations(): Array<SyncDestination> {
  return [...SYNC_DESTINATIONS];
}

export function destinationTypes(): Array<DestinationType> {
  return SYNC_DESTINATIONS.map((d) => d.type);
}

export function findDestination(type: DestinationType) {
  const destination = SYNC_DESTINATIONS.find((d) => d.type === type);
  if (!destination) {
    throw new Error(`Destination not found for type: ${type}`);
  }
  return destination;
}

export default buildHelper(syncDestinations);
