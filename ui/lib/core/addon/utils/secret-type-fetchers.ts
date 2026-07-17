/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  SecretsApiKvV2ListListEnum,
  SecretsApiDatabaseListStaticRolesListEnum,
} from '@hashicorp/vault-client-typescript';
import { parentKeyForKey, keyIsFolder } from 'core/utils/key-utils';

import type ApiService from 'vault/services/api';
import type { SecretType } from 'vault/sync';

export interface SecretTypeFetcher {
  fetch: (api: ApiService, mountPath: string, value: string) => Promise<string[]>;
  filter: (items: string[], value: string, isDirectory: boolean) => string[];
  onSelect: (item: string, pathToSecret: string) => string;
}

function keyWithoutParentKey(value: string): string {
  const parentKey = parentKeyForKey(value);
  return value ? value.replace(parentKey, '') : '';
}

export const SECRET_TYPE_FETCHERS: Record<SecretType, SecretTypeFetcher> = {
  kv: {
    fetch: async (api: ApiService, mountPath: string, value: string) => {
      try {
        const backend = keyIsFolder(mountPath) ? mountPath.slice(0, -1) : mountPath;
        const isDirectory = keyIsFolder(value);
        const parentDirectory = parentKeyForKey(value);
        const pathToSecret = isDirectory ? value : parentDirectory;

        const { keys } = await api.secrets.kvV2List(pathToSecret, backend, SecretsApiKvV2ListListEnum.TRUE);
        return keys || [];
      } catch (error) {
        return [];
      }
    },
    filter: (secrets: string[], value: string, isDirectory: boolean) => {
      const secretName = keyWithoutParentKey(value) || '';
      return secrets.filter((path) => {
        if (!value || isDirectory) {
          return true;
        }
        // Exclude exact matches to avoid showing current selection in suggestions
        if (secretName === path) {
          return false;
        }
        return path.toLowerCase().includes(secretName.toLowerCase());
      });
    },
    onSelect: (item: string, pathToSecret: string) => `${pathToSecret}${item}`,
  },
  database: {
    fetch: async (api: ApiService, mountPath: string) => {
      try {
        const backend = mountPath.endsWith('/') ? mountPath.slice(0, -1) : mountPath;
        const { keys } = await api.secrets.databaseListStaticRoles(
          backend,
          SecretsApiDatabaseListStaticRolesListEnum.TRUE
        );
        return (keys || []).map((key) => `static-roles/${key}`);
      } catch (error) {
        return [];
      }
    },
    filter: (roles: string[], value: string) => {
      if (!value) {
        return roles;
      }
      return roles.filter((role) => {
        // Exclude exact matches to avoid showing current selection in suggestions
        if (value === role) {
          return false;
        }
        return role.toLowerCase().includes(value.toLowerCase());
      });
    },
    onSelect: (item: string) => item,
  },
};
