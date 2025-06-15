/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

import type { DestinationType } from 'vault/sync';
import type { SyncDestination } from 'vault/helpers/sync-destinations';

/* 
This helper is used to lookup static display properties for sync destinations
maskedParams: attributes for sensitive data, the API returns these values as '*****'
*/

const SYNC_DESTINATIONS: Array<SyncDestination> = [
  {
    name: 'AWS Secrets Manager',
    type: 'aws-sm',
    icon: 'aws-color',
    category: 'cloud',
    maskedParams: ['accessKeyId', 'secretAccessKey'],
    readonlyParams: ['name', 'region'],
    defaultValues: {
      granularity: 'secret-path',
    },
  },
  {
    name: 'Azure Key Vault',
    type: 'azure-kv',
    icon: 'azure-color',
    category: 'cloud',
    maskedParams: ['clientSecret'],
    readonlyParams: ['name', 'keyVaultUri', 'tenantId', 'cloud'],
    defaultValues: {
      granularity: 'secret-path',
    },
  },
  {
    name: 'Google Secret Manager',
    type: 'gcp-sm',
    icon: 'gcp-color',
    category: 'cloud',
    maskedParams: ['credentials'],
    readonlyParams: ['name'],
    defaultValues: {
      granularity: 'secret-path',
    },
  },
  {
    name: 'Github Actions',
    type: 'gh',
    icon: 'github-color',
    category: 'dev-tools',
    maskedParams: ['accessToken'],
    readonlyParams: ['name', 'repositoryOwner', 'repositoryName'],
    defaultValues: {
      granularity: 'secret-key',
    },
  },
  {
    name: 'Vercel Project',
    type: 'vercel-project',
    icon: 'vercel-color',
    category: 'dev-tools',
    maskedParams: ['accessToken'],
    readonlyParams: ['name', 'projectId'],
    defaultValues: {
      granularity: 'secret-key',
      deploymentEnvironments: [],
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
