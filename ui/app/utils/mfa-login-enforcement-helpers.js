/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { filterEnginesByMountCategory } from 'core/utils/all-engines-metadata';
import { addManyToArray, addToArray } from 'vault/helpers/add-to-array';
import { capitalize } from '@ember/string';

/**
 * Prepares targets for MFA login enforcement display
 * @param {Object} enforcement - The enforcement object containing target data
 * @param {Array} enforcement.auth_method_accessors - Array of auth method accessors
 * @param {Array} enforcement.auth_method_types - Array of auth method types
 * @param {Array} enforcement.identity_entities - Array of identity entities (or identity_entity_ids)
 * @param {Array} enforcement.identity_groups - Array of identity groups (or identity_group_ids)
 * @param {Object} api - The API service for fetching auth methods and identity data
 * @param {Object} options - Optional configuration
 * @param {boolean} options.includeFormFields - If true, includes label, key, and value fields for form usage
 * @returns {Promise<Array>} Array of prepared target objects
 */
export async function prepareTargets(enforcement, api, options = {}) {
  let authMethods;
  let targets = [];
  const { includeFormFields = false } = options;

  if (enforcement.auth_method_accessors?.length || enforcement.auth_method_types?.length) {
    // fetch all auth methods and lookup by accessor to get mount path and type
    try {
      const { data } = await api.sys.authListEnabledMethods();
      authMethods = api.responseObjectToArray(data, 'path');
    } catch (error) {
      // swallow this error
    }
  }

  if (enforcement.auth_method_accessors?.length) {
    const selectedAuthMethods = authMethods.filter((method) => {
      return enforcement.auth_method_accessors.includes(method.accessor);
    });
    targets = addManyToArray(
      targets,
      selectedAuthMethods.map((method) => ({
        ...(includeFormFields && {
          label: 'Authentication mount',
          key: 'auth_method_accessors',
          value: method.accessor,
        }),
        icon: iconForMount(method.type),
        link: 'vault.cluster.access.method',
        linkModels: [method.path.slice(0, -1)],
        title: method.path,
        subTitle: method.accessor,
      }))
    );
  }

  enforcement.auth_method_types?.forEach((type) => {
    const icon = iconForMount(type);
    const mountCount = authMethods ? authMethods.filterBy('type', type).length : 0;
    targets = addToArray(targets, {
      ...(includeFormFields && {
        label: 'Authentication method',
        key: 'auth_method_types',
        value: type,
      }),
      icon,
      title: type,
      subTitle: `All ${type} mounts (${mountCount})`,
    });
  });

  // Handle identity entities and groups using API service
  const types = [
    {
      key: 'identity_entities',
      idKey: 'identity_entity_ids',
      linkType: 'entities',
      label: 'Entity',
      apiMethod: 'entityListById',
    },
    {
      key: 'identity_groups',
      idKey: 'identity_group_ids',
      linkType: 'groups',
      label: 'Group',
      apiMethod: 'groupListById',
    },
  ];

  for (const { key, idKey, linkType, label, apiMethod } of types) {
    const itemKey = enforcement[key] ? key : idKey;
    const ids = enforcement[itemKey] || [];

    if (ids.length > 0) {
      try {
        // Fetch all entities/groups at once using list endpoint
        const response = await api.identity[apiMethod](true);
        const allItems = api.keyInfoToArray(response);

        // Filter to only the IDs we need and map to target objects
        ids.forEach((id) => {
          const item = allItems.find((i) => i.id === id);
          if (item) {
            targets = addToArray(targets, {
              ...(includeFormFields && {
                label,
                key: idKey,
                value: item,
              }),
              icon: 'user',
              link: 'vault.cluster.access.identity.show',
              linkModels: [linkType, id, 'details'],
              title: item.name,
              subTitle: id,
            });
          }
        });
      } catch (error) {
        // Skip items that can't be loaded
      }
    }
  }

  return targets;
}

/**
 * Returns the icon for a given mount type
 * @param {string} type - The mount type
 * @returns {string} The icon name
 */
export function iconForMount(type) {
  const mountableMethods = filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: true });
  const mount = mountableMethods.find((method) => method.type === type);
  return mount ? mount.glyph || mount.type : 'token';
}

/**
 * Returns the icon for a given MFA method type
 * @param {string} type - The MFA method type (e.g., 'totp', 'duo', 'pingid')
 * @returns {string} The icon name
 */
export function getMfaMethodIcon(type) {
  switch (type) {
    case 'totp':
      return 'history';
    case 'pingid':
      return 'ping-identity-color';
    case 'duo':
      return 'duo-color';
    default:
      return type;
  }
}

/**
 * Returns the display name for a given MFA method type
 * @param {string} type - The MFA method type (e.g., 'totp', 'duo', 'pingid')
 * @returns {string} The formatted display name
 */
export function getMfaMethodName(type) {
  return type === 'totp' ? type.toUpperCase() : capitalize(type);
}

/**
 * Fetches all MFA methods by listing method IDs and reading each method's details
 * @param {Object} api - The API service instance
 * @returns {Promise<Array>} Array of MFA method data objects
 */
export async function fetchMfaMethods(api) {
  const response = await api.identity.mfaListMethods(true);
  const mfaMethods = api.keyInfoToArray(response);

  mfaMethods.forEach((method) => {
    method.displayName = getMfaMethodName(method.type);
    method.icon = getMfaMethodIcon(method.type);
  });

  return mfaMethods;
}

/**
 * Fetches all MFA login enforcements by listing enforcement names and reading each enforcement's details
 * @param {Object} api - The API service instance
 * @returns {Promise<Array>} Array of MFA login enforcement data objects
 */
export async function fetchMfaLoginEnforcements(api) {
  const response = await api.identity.mfaListLoginEnforcements(true);
  return api.keyInfoToArray(response);
}
