/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

import type { SyncDestination, SyncDestinationType } from 'vault/vault/helpers/sync-destinations';

/* 
This helper is referenced in the base sync destination model and elsewhere to set attributes that rely on type
maskedParams: attributes for sensitive data, the API returns these values as '*****'
*/

const SYNC_DESTINATIONS: Array<SyncDestination> = [
  {
    name: 'AWS Secrets Manager',
    type: 'aws-sm',
    icon: 'aws-color',
    category: 'cloud',
    maskedParams: ['accessKeyId', 'secretAccessKey'],
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
    defaultValues: {
      granularity: 'secret-key',
    },
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
