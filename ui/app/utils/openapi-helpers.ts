/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { dasherize } from '@ember/string';
import { singularize } from 'ember-inflector';

// TODO: Consolidate with openapi-to-attrs once it's typescript

interface Path {
  path: string;
  itemType: string;
  itemName: string;
  operations: string[];
  action: string;
  navigation: boolean;
  param: string | false;
}
interface PathsInfo {
  apiPath: string;
  itemType: string;
  itemTypes: string[];
  paths: Path[];
}

interface OpenApiParameter {
  description?: string;
  in: string;
  name: string;
  required: boolean;
  schema: object;
}
interface DisplayAttrs {
  itemType: string;
  action: string;
  navigation?: boolean;
  description?: string;
  name?: string;
  group?: string;
  value?: string | number;
  sensitive?: boolean;
}
interface OpenApiAction {
  parameters: Array<{ name: string }>;
}
interface OpenApiPath {
  description?: string;
  parameters: OpenApiParameter[];
  'x-vault-displayAttrs': DisplayAttrs;
  get?: OpenApiAction;
  post?: OpenApiAction;
  delete?: OpenApiAction;
}

// Take object entries from the OpenAPI response and consolidate them into an object which includes itemTypes, operations, and paths
export function reducePathsByPathName(pathsInfo: PathsInfo, currentPath: [string, OpenApiPath]): PathsInfo {
  const pathName = currentPath[0];
  const pathDetails = currentPath[1];
  const displayAttrs = pathDetails['x-vault-displayAttrs'];
  if (!displayAttrs) {
    // don't include paths that don't have display attrs
    return pathsInfo;
  }

  let itemType, itemName;
  if (displayAttrs.itemType) {
    itemType = displayAttrs.itemType;
    let items = itemType.split(':');
    itemName = items[items.length - 1];
    items = items.map((item) => dasherize(singularize(item.toLowerCase())));
    itemType = items.join('~*');
  }

  if (itemType && !pathsInfo.itemTypes.includes(itemType)) {
    pathsInfo.itemTypes.push(itemType);
  }

  const operations = [];
  if (pathDetails.get) {
    operations.push('get');
  }
  if (pathDetails.post) {
    operations.push('post');
  }
  if (pathDetails.delete) {
    operations.push('delete');
  }
  if (pathDetails.get && pathDetails.get.parameters && pathDetails.get.parameters[0]?.name === 'list') {
    operations.push('list');
  }

  pathsInfo.paths.push({
    path: pathName,
    itemType: itemType || displayAttrs.itemType,
    itemName: itemName || pathsInfo.itemType || displayAttrs.itemType,
    operations,
    action: displayAttrs.action,
    navigation: displayAttrs.navigation === true,
    param: _getPathParam(pathName),
  });

  return pathsInfo;
}

const apiPathRegex = new RegExp(/\{\w+\}/, 'g');

/**
 * getPathParam takes an OpenAPI url and returns the first path param name, if it exists.
 * This is an internal method, but exported for testing.
 */
export function _getPathParam(pathName: string): string | false {
  if (!pathName) return false;
  const params = pathName.match(apiPathRegex);
  // returns array like ['{username}'] or null
  if (!params) return false;
  // strip curly brackets from param name
  // previous behavior only returned the first param, so we match that for now
  return params[0]?.replace(new RegExp('{|}', 'g'), '') || false;
}

export function pathToHelpUrlSegment(path: string): string {
  if (!path) return '';
  return path.replaceAll(apiPathRegex, 'example');
}

export function filterPathsByItemType(pathInfo: PathsInfo, itemType: string): Path[] {
  if (!itemType) {
    return pathInfo.paths;
  }
  return pathInfo.paths.filter((path) => {
    return itemType === path.itemType;
  });
}

/**
 * This object maps model names to the openAPI path that hydrates the model, given the backend path.
 */
const OPENAPI_POWERED_MODELS = {
  'role-ssh': (backend: string) => `/v1/${backend}/roles/example?help=1`,
  'auth-config/azure': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/cert': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/gcp': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/github': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/jwt': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/kubernetes': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/ldap': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/okta': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'auth-config/radius': (backend: string) => `/v1/auth/${backend}/config?help=1`,
  'kmip/config': (backend: string) => `/v1/${backend}/config?help=1`,
  'kmip/role': (backend: string) => `/v1/${backend}/scope/example/role/example?help=1`,
  'pki/role': (backend: string) => `/v1/${backend}/roles/example?help=1`,
  'pki/tidy': (backend: string) => `/v1/${backend}/config/auto-tidy?help=1`,
  'pki/sign-intermediate': (backend: string) => `/v1/${backend}/issuer/example/sign-intermediate?help=1`,
  'pki/certificate/generate': (backend: string) => `/v1/${backend}/issue/example?help=1`,
  'pki/certificate/sign': (backend: string) => `/v1/${backend}/sign/example?help=1`,
  'pki/config/acme': (backend: string) => `/v1/${backend}/config/acme?help=1`,
  'pki/config/cluster': (backend: string) => `/v1/${backend}/config/cluster?help=1`,
  'pki/config/urls': (backend: string) => `/v1/${backend}/config/urls?help=1`,
};

export function getHelpUrlForModel(modelType: string, backend: string) {
  const urlFn = OPENAPI_POWERED_MODELS[modelType as keyof typeof OPENAPI_POWERED_MODELS] as (
    backend: string
  ) => string;
  if (!urlFn) return null;
  return urlFn(backend);
}
