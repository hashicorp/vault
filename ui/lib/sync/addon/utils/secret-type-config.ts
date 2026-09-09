/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { SecretType } from 'vault/sync';

export interface MountOption {
  name: string;
  id: string;
  engineType: string;
  version?: number;
}

export function getSecretTypeFromMount(type: string, version?: number): SecretType | null {
  // KV requires version 2 for sync
  if (type === 'kv') {
    return version === 2 ? 'kv' : null;
  }

  if (type === 'database') {
    return 'database';
  }

  return null;
}

// Extract the engine type from the accessor field (e.g. "kv_9b39cc0f" -> "kv", "database_03ce3af1" -> "database")
export function getSecretTypeFromAccessor(accessor: string): SecretType | null {
  const prefix = accessor.split('_')[0];
  if (prefix === 'kv' || prefix === 'database') {
    return prefix;
  }
  return null;
}

export interface SecretTypeConfig {
  placeholder: string;
  noMatchesMessage: string;
  accessorType: string;
  icon: string;
  route: string;
  supportsExternalLink: boolean;
  getModels: (mount: string, secretName: string) => string[];
  getQuery?: () => Record<string, string>;
}

export const SECRET_TYPE_CONFIGS: Record<SecretType, SecretTypeConfig> = {
  kv: {
    placeholder: 'Path to secret',
    noMatchesMessage: 'No suggestions for this path',
    accessorType: 'KV v2',
    icon: 'key-values',
    route: 'kvSecretOverview',
    supportsExternalLink: true,
    getModels: (mount: string, secretName: string) => [mount, secretName],
  },
  database: {
    placeholder: 'Static role name',
    noMatchesMessage: 'No matching static roles found',
    accessorType: 'Database',
    icon: 'database',
    route: 'databaseStaticRoleOverview',
    supportsExternalLink: true,
    // secret_name includes the "static-roles/" prefix (e.g. "static-roles/my-role"); strip it for the view route
    getModels: (mount: string, secretName: string) => {
      const roleName = secretName.startsWith('static-roles/')
        ? secretName.slice('static-roles/'.length)
        : secretName;
      return [mount, `role/${roleName}`];
    },
    getQuery: () => ({ type: 'static' }),
  },
};
