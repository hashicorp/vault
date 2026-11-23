/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { capitalize, dasherize } from '@ember/string';
import { singularize } from 'ember-inflector';

interface Path {
  path: string;
  itemType: string;
  itemName: string;
  operations: string[];
  action: string;
  navigation: boolean;
  param: string | false;
}
export type PathInfo = {
  apiPath: string;
  paths: Path[];
  itemTypes: string[];
  itemType?: string;
  itemID?: string;
};

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
  operationId: string;
  tags?: string[];
  responses: Record<string, { description: string }>;
  requestBody?: {
    content: Record<string, { schema: { $ref: string } }>;
    required: boolean;
  };
  parameters: Array<{ name: string }>;
}
interface OpenApiPath {
  description?: string;
  parameters?: OpenApiParameter[];
  'x-vault-displayAttrs': DisplayAttrs;
  get?: OpenApiAction;
  post?: OpenApiAction;
  delete?: OpenApiAction;
}
interface OpenApiProp {
  description: string;
  type: string;
  'x-vault-displayAttrs': {
    name?: string;
    value?: string | number;
    group?: string;
    sensitive?: boolean;
    editType?: string;
    description?: string;
    identifier?: boolean;
  };
  items?: { type: string };
  format?: string;
  isId?: boolean;
  deprecated?: boolean;
  enum?: string[];
}
export interface MixedAttr {
  type?: string;
  editType?: string;
  fieldGroup: string;
  fieldValue?: string;
  label?: string;
  readonly?: boolean;
  possibleValues?: string[];
  defaultValue?: string | number;
  sensitive?: boolean;
  readOnly?: boolean;
  name?: string;
  identifier?: boolean;
  [key: string]: unknown;
}

export type OpenApiProps = Record<string, OpenApiProp>;

export type OpenApiHelpResponse = {
  help: string;
  openapi: {
    paths: OpenApiPath[];
    components: {
      schemas: Record<string, { properties: OpenApiProps; type: string }>;
    };
    info: {
      title: string;
      version: string;
      description: string;
      license: { name: string; url: string };
    };
    openapi: string;
  };
};

// Take object entries from the OpenAPI response and consolidate them into an object which includes itemTypes, operations, and paths
export function reducePathsByPathName(pathsInfo: PathInfo, currentPath: [string, OpenApiPath]): PathInfo {
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

export function filterPathsByItemType(pathInfo: PathInfo, itemType: string): Path[] {
  if (!itemType) {
    return pathInfo.paths;
  }
  return pathInfo.paths.filter((path) => {
    return itemType === path.itemType;
  });
}

export const expandOpenApiProps = function (props: OpenApiProps): Record<string, MixedAttr> {
  const attrs: Record<string, MixedAttr> = {};
  // expand all attributes
  for (const propName in props) {
    const prop = props[propName];
    if (!prop) continue;
    let { description, items, type, format, isId, deprecated } = prop;
    if (deprecated === true) {
      continue;
    }
    let {
      name,
      value,
      group,
      sensitive,
      editType,
      description: displayDescription,
      identifier,
    } = prop['x-vault-displayAttrs'] || {};

    if (type === 'integer') {
      type = 'number';
    }

    if (displayDescription) {
      description = displayDescription;
    }

    editType = editType || type;

    if (format === 'seconds' || format === 'duration') {
      editType = 'ttl';
    } else if (items) {
      editType = items.type + capitalize(type);
    }

    const attrDefn: MixedAttr = {
      editType,
      subText: description,
      possibleValues: prop['enum'],
      fieldGroup: group || 'default',
      readOnly: isId,
      defaultValue: value || undefined,
      identifier,
    };

    if (sensitive) {
      attrDefn.sensitive = true;
    }

    // only set a label if we have one from OpenAPI
    // otherwise the propName will be humanized by the form-field component
    if (name) {
      attrDefn.label = name;
    }

    // ttls write as a string and read as a number
    // so setting type on them runs the wrong transform
    if (editType !== 'ttl' && type !== 'array') {
      attrDefn.type = type;
    }

    // loop to remove empty vals
    for (const attrProp in attrDefn) {
      if (attrDefn[attrProp] == null) {
        delete attrDefn[attrProp];
      }
    }

    attrs[propName] = attrDefn;
  }
  return attrs;
};

/*
 * extract props for post request schema from OpenAPI help response
 * returns expanded OpenAPI props to be used in forms
 */
export const propsForSchema = function (helpResponse: OpenApiHelpResponse) {
  const { openapi } = helpResponse;
  // paths is an array but it will have a single entry with the scope we're in
  const path = Object.values(openapi.paths)[0];
  const schema = path?.post?.requestBody?.content['application/json']?.schema;
  if (schema?.$ref) {
    // $ref will be shaped like `#/components/schemas/MyResponseType
    // this maps to the location of the item in openapi.components.schemas
    const schemaRef = schema.$ref.replace('#/components/schemas/', '');
    const props = openapi.components.schemas[schemaRef]?.['properties'] as OpenApiProps;
    return expandOpenApiProps(props);
  }

  return {};
};
