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
    maskedParams: ['access_key_id', 'secret_access_key'],
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
    maskedParams: ['client_secret'],
    readonlyParams: ['name', 'key_vault_uri', 'tenant_id', 'cloud'],
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
    maskedParams: ['access_token'],
    readonlyParams: ['name', 'repository_owner', 'repository_name'],
    defaultValues: {
      granularity: 'secret-key',
    },
  },
  {
    name: 'Vercel Project',
    type: 'vercel-project',
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
