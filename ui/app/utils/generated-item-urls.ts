/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { sanitizePath } from 'core/utils/sanitize-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { filterPathsByItemType, pathToHelpUrlSegment } from 'vault/utils/openapi-helpers';

import type { PathInfo } from 'vault/utils/openapi-helpers';

/** Request paths for one item type of an auth mount, relative to the API base path (`/v1`). */
export interface GeneratedItemUrls {
  list: string;
  read: (id: string) => string;
  create: (id: string) => string;
  update: (id: string) => string;
  delete: (id: string) => string;
}

// "/users/{username}" -> "users": the collection segment before the path parameter.
const collectionSegment = (path = '') => path.slice(1, path.indexOf('{') - 1);

/**
 * Builds the item URLs from a mount's OpenAPI paths. Only reads encode the id; writes use it as entered,
 * matching the generated-item adapter this replaces.
 */
export function generatedItemUrls(pathInfo: PathInfo, itemType: string): GeneratedItemUrls {
  const paths = filterPathsByItemType(pathInfo, itemType);
  const getPath = paths.find((path) => path.operations.includes('get'))?.path || '';
  const createPath = paths.find((path) => path.action === 'Create' || path.operations.includes('post'))?.path;
  const deletePath = paths.find((path) => path.operations.includes('delete'))?.path;
  const base = `/${pathInfo.apiPath}`;
  const listSegment = sanitizePath(getPath);

  return {
    list: `${base}${listSegment}/`,
    read: (id) => `${base}${listSegment}/${encodePath(id)}`,
    create: (id) => `${base}${collectionSegment(createPath)}/${id}`,
    update: (id) => `${base}${collectionSegment(createPath)}/${id}`,
    delete: (id) => `${base}${collectionSegment(deletePath)}/${id}`,
  };
}

/** OpenAPI help URL for the item type's create endpoint, which describes the form fields. */
export function generatedItemHelpUrl(pathInfo: PathInfo, itemType: string): string {
  const paths = itemType ? filterPathsByItemType(pathInfo, itemType) : pathInfo.paths;
  const createPath = paths.find((path) => path.operations.includes('post') && path.action !== 'Delete');
  if (!createPath) {
    throw new Error(`No create endpoint found for ${itemType} in ${pathInfo.apiPath}`);
  }
  return `/v1/${pathInfo.apiPath}${pathToHelpUrlSegment(createPath.path).slice(1)}?help=true`;
}
