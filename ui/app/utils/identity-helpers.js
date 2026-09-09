/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  IdentityApiEntityListByIdListEnum,
  IdentityApiGroupListByIdListEnum,
} from '@hashicorp/vault-client-typescript';

/**
 * Fetches identity items (entities or groups) from the API
 * @param {Object} params - Parameters object
 * @param {string} params.identityType - The type of identity ('entity' or 'group')
 * @param {Object} params.api - The API service instance
 * @returns {Promise<Array>} Array of identity items
 */
export async function fetchIdentityItems({ identityType, api }) {
  const methodType = identityType === 'group' ? 'groupListById' : 'entityListById';
  const listEnum =
    identityType === 'group' ? IdentityApiGroupListByIdListEnum : IdentityApiEntityListByIdListEnum;
  const response = await api.identity[methodType](listEnum.TRUE);

  return api.keyInfoToArray(response);
}

/**
 * Fetches identity items (entities or groups) with their capabilities attached
 * @param {Object} params - Parameters object
 * @param {string} params.identityType - The type of identity ('entity' or 'group')
 * @param {Object} params.api - The API service instance
 * @param {Object} params.capabilities - The capabilities service instance
 * @returns {Promise<Array>} Array of items with capabilities attached (canDelete, canEdit, canAddAlias, type, alias)
 */
export async function fetchIdentityItemsWithCapabilities({ identityType, api, capabilities }) {
  const items = await fetchIdentityItems({ identityType, api });

  // Build capability paths for all items
  const capabilityPaths = items.flatMap((item) => [
    capabilities.pathFor('identityCapabilities', {
      identityType,
      id: item.id,
    }),
    capabilities.pathFor('groupAlias'),
  ]);

  // For groups, fetch detailed data (type and alias) in parallel with capabilities
  const detailsPromise =
    identityType === 'group'
      ? Promise.all(items.map((item) => api.identity.groupReadById(item.id)))
      : Promise.resolve([]);

  // Fetch capabilities and details in parallel
  const [capabilitiesMap, detailsResponses] = await Promise.all([
    capabilities.fetch(capabilityPaths),
    detailsPromise,
  ]);

  // Create a map of item details for quick lookup
  const detailsMap = {};
  if (identityType === 'group') {
    detailsResponses.forEach((response) => {
      if (response?.data) {
        detailsMap[response.data.id] = {
          type: response.data.type,
          alias: response.data.alias,
        };
      }
    });
  }

  // Attach capabilities and details to each item
  const groupAliasPath = capabilities.pathFor('groupAlias');
  const groupAliasCapabilities = capabilitiesMap[groupAliasPath];

  const itemsWithCapabilities = items.map((item) => {
    const itemCapabilityPath = capabilities.pathFor('identityCapabilities', {
      identityType,
      id: item.id,
    });
    const itemCapabilities = capabilitiesMap[itemCapabilityPath];
    const details = detailsMap[item.id] || {};

    return {
      ...item,
      ...details,
      canDelete: itemCapabilities?.canDelete || false,
      canEdit: itemCapabilities?.canUpdate || false,
      canAddAlias: groupAliasCapabilities?.canCreate || false,
    };
  });

  return itemsWithCapabilities;
}

/**
 * Build request parameters for group operations
 * @param {Object} data - Form data
 * @returns {Object} Request parameters for group API calls
 */
export function buildGroupRequestParams(data) {
  return {
    name: data.name,
    type: data.type,
    metadata: data.metadata,
    policies: data.policies,
    member_group_ids: data.member_group_ids,
    member_entity_ids: data.member_entity_ids,
  };
}

/**
 * Build request parameters for entity operations
 * @param {Object} data - Form data
 * @returns {Object} Request parameters for entity API calls
 */
export function buildEntityRequestParams(data) {
  return {
    name: data.name,
    disabled: data.disabled,
    policies: data.policies,
    metadata: data.metadata,
  };
}

/**
 * Build request parameters for entity merge operations
 * @param {Object} data - Form data
 * @returns {Object} Request parameters for entity merge API calls
 */
export function buildEntityMergeRequestParams(data) {
  return {
    from_entity_ids: data.from_entity_ids,
    to_entity_id: data.to_entity_id,
    force: data.force,
  };
}

/**
 * Build request parameters for alias operations
 * @param {Object} data - Form data
 * @returns {Object} Request parameters for alias API calls
 */
export function buildAliasRequestParams(data) {
  return {
    canonical_id: data.canonical_id,
    mount_accessor: data.mount_accessor,
    name: data.name,
  };
}

/**
 * Handle entity merge operation
 * @param {Object} params - Parameters object
 * @param {Object} params.api - The API service instance
 * @param {Object} params.data - Form data
 * @returns {Promise<Object>} API response
 */
export async function handleEntityMerge({ api, data }) {
  const params = buildEntityMergeRequestParams(data);
  return await api.identity.entityMerge(params);
}

/**
 * Handle alias creation
 * @param {Object} params - Parameters object
 * @param {Object} params.api - The API service instance
 * @param {Object} params.model - The model object
 * @param {Object} params.data - Form data
 * @returns {Promise<Object>} API response
 */
export async function handleAliasCreate({ api, model, data }) {
  const params = buildAliasRequestParams(data);
  const method = model.identityType === 'group' ? 'groupCreateAlias' : 'entityCreateAlias';
  return await api.identity[method](params);
}

/**
 * Handle alias update
 * @param {Object} params - Parameters object
 * @param {Object} params.api - The API service instance
 * @param {Object} params.model - The model object
 * @param {Object} params.data - Form data
 * @returns {Promise<Object>} API response
 */
export async function handleAliasUpdate({ api, model, data }) {
  const params = buildAliasRequestParams(data);
  const method = model.identityType === 'group' ? 'groupUpdateAliasById' : 'entityUpdateAliasById';
  return await api.identity[method](params);
}

/**
 * Handle entity/group update
 * @param {Object} params - Parameters object
 * @param {Object} params.api - The API service instance
 * @param {Object} params.model - The model object
 * @param {Object} params.data - Form data
 * @returns {Promise<Object>} API response
 */
export async function handleUpdate({ api, model, data }) {
  const params = buildGroupRequestParams(data);
  const method = model.identityType === 'group' ? 'groupUpdateById' : 'entityUpdateById';
  return await api.identity[method](model.itemId, params);
}

/**
 * Handle entity/group creation
 * @param {Object} params - Parameters object
 * @param {Object} params.api - The API service instance
 * @param {Object} params.model - The model object
 * @param {Object} params.data - Form data
 * @returns {Promise<Object>} API response
 */
export async function handleCreate({ api, model, data }) {
  const isGroup = model.identityType === 'group';
  const method = isGroup ? 'groupCreate' : 'entityCreate';
  const params = isGroup ? buildGroupRequestParams(data) : buildEntityRequestParams(data);
  return await api.identity[method](params);
}

/**
 * Determine which API operation to perform based on mode and model type
 * @param {Object} params - Parameters object
 * @param {Object} params.api - The API service instance
 * @param {Object} params.model - The model object
 * @param {string} params.mode - The operation mode ('create', 'edit', 'merge')
 * @param {Object} params.data - Form data
 * @returns {Promise<Object>} API response
 */
export async function performSaveOperation({ api, model, mode, data }) {
  const isAlias = model.form.identityFormType === 'alias';

  if (mode === 'merge') {
    return await handleEntityMerge({ api, data });
  }

  if (isAlias) {
    return mode === 'create'
      ? await handleAliasCreate({ api, model, data })
      : await handleAliasUpdate({ api, model, data });
  }

  return mode === 'edit'
    ? await handleUpdate({ api, model, data })
    : await handleCreate({ api, model, data });
}

/**
 * Extract the ID from the save response
 * @param {Object} params - Parameters object
 * @param {string} params.mode - The operation mode
 * @param {Object} params.data - Form data
 * @param {Object} params.response - API response
 * @param {Object} params.model - The model object
 * @returns {string} The extracted ID
 */
export function extractSavedId({ mode, data, response, model }) {
  if (mode === 'merge') {
    return data?.to_entity_id;
  }
  return response?.data?.id || model.itemId;
}
