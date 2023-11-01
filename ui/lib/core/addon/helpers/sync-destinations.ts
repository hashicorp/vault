/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

import type { SyncDestination, SyncDestinationType } from 'vault/vault/app-types';

/* 
This helper is referenced in the base sync destination model
to return static display attributes that rely on type
maskedParams: attributes for sensitive data, the API returns these values as '*****'
*/

// TODO update maskedParams for other types: https://hashicorp.atlassian.net/browse/VAULT-21428

const SYNC_DESTINATIONS: Array<SyncDestination> = [
  {
    name: 'AWS Secrets Manager',
    type: 'aws-sm',
    icon: 'aws-color',
    category: 'cloud',
    maskedParams: ['accessKeyId', 'secretAccessKey'],
  },
  {
    name: 'Azure Key Vault',
    type: 'azure-kv',
    icon: 'azure-color',
    category: 'cloud',
    maskedParams: [],
  },
  {
    name: 'Google Secret Manager',
    type: 'gcp-sm',
    icon: 'gcp-color',
    category: 'cloud',
    maskedParams: [],
  },
  {
    name: 'Github Actions',
    type: 'gh',
    icon: 'github-color',
    category: 'dev-tools',
    maskedParams: ['accessToken'],
  },
  {
    name: 'Vercel Project',
    type: 'vercel-project',
    icon: 'vercel-color',
    category: 'dev-tools',
    maskedParams: [],
  },
];

export function syncDestinations(): Array<SyncDestination> {
  return [...SYNC_DESTINATIONS];
}

export function destinationTypes(): Array<SyncDestinationType> {
  return SYNC_DESTINATIONS.map((d) => d.type);
}

export function findDestination(type: SyncDestinationType | undefined): SyncDestination | undefined {
  return SYNC_DESTINATIONS.find((d) => d.type === type);
}

export default buildHelper(syncDestinations);
