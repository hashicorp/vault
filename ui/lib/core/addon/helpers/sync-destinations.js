/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

// These helper is referenced in the base sync destination model
// to return static display attributes that rely on type
const SYNC_DESTINATIONS = [
  {
    name: 'AWS Secrets Manager',
    type: 'aws-sm',
    icon: 'aws-color',
    category: 'cloud',
  },
  {
    name: 'Azure Key Vault',
    type: 'azure-kv',
    icon: 'azure-color',
    category: 'cloud',
  },
  {
    name: 'Google Secret Manager',
    type: 'gcp-sm',
    icon: 'gcp-color',
    category: 'cloud',
  },
  {
    name: 'Github Actions',
    type: 'gh',
    icon: 'github-color',
    category: 'dev-tools',
  },
  {
    name: 'Vercel Project',
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
  if (!type) return;
  assert(
    `you must pass one of the following types: ${destinationTypes().join(', ')}`,
    destinationTypes().includes(type)
  );
  return SYNC_DESTINATIONS.find((d) => d.type === type);
}

export default buildHelper(syncDestinations);
