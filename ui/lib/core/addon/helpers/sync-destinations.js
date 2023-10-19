/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

const SYNC_DESTINATIONS = [
  {
    displayName: 'AWS Secrets Manager',
    type: 'aws-sm',
    icon: 'aws-color',
    category: 'cloud',
  },
  {
    displayName: 'Azure Key Vault',
    type: 'azure-kv',
    icon: 'azure-color',
    category: 'cloud',
  },
  {
    displayName: 'Google Secret Manager',
    type: 'gcp-sm',
    icon: 'gcp-color',
    category: 'cloud',
  },
  {
    displayName: 'Github Actions',
    type: 'gh',
    icon: 'github-color',
    category: 'dev-tools',
  },
  {
    displayName: 'Vercel Project',
    type: 'vercel-project',
    icon: 'vercel-color',
    category: 'dev-tools',
  },
];

export function syncDestinations() {
  return [...SYNC_DESTINATIONS];
}

export function destinationTypes() {
  return SYNC_DESTINATIONS.map((d) => d.type);
}

export function findDestination(type) {
  assert(
    `you must pass one of the following types: ${destinationTypes().join(', ')}`,
    destinationTypes().includes(type)
  );
  return SYNC_DESTINATIONS.find((d) => d.type === type);
}

export default buildHelper(syncDestinations);
