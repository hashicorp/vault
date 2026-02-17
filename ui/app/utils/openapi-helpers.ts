/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { debug } from '@ember/debug';
import { camelize, capitalize, dasherize } from '@ember/string';
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
interface Attribute {
  name: string;
  type: string | undefined;
  options: {
    editType?: string;
    fieldGroup?: string;
    fieldValue?: string;
    label?: string;
    readonly?: boolean;
  };
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
interface MixedAttr {
  type?: string;
  helpText?: string;
  editType?: string;
  fieldGroup: string;
  fieldValue?: string;
  label?: string;
  readonly?: boolean;
  possibleValues?: string[];
  defaultValue?: string | number | (() => string | number);
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

/**
 * This object maps model names to the openAPI path that hydrates the model, given the backend path.
 */
const OPENAPI_POWERED_MODELS = {
  'role-ssh': (backend: string) => `/v1/${backend}/roles/example?help=1`,
};

export function getHelpUrlForModel(modelType: string, backend?: string) {
  if (modelType in OPENAPI_POWERED_MODELS && backend) {
    const urlFn = OPENAPI_POWERED_MODELS[modelType as keyof typeof OPENAPI_POWERED_MODELS];
    return urlFn(backend);
  }
  return null;
}

export const expandOpenApiProps = function (
  props: OpenApiProps,
  outputFormat: 'model' | 'form' = 'model'
): Record<string, MixedAttr> {
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

    // to preserve legacy behavior set the description to helpText
    // this pattern is deprecated when used with Form class since description is better represented as subText
    const descriptionKey = outputFormat === 'model' ? 'helpText' : 'subText';
    const attrDefn: MixedAttr = {
      editType,
      [descriptionKey]: description,
      possibleValues: prop['enum'],
      fieldValue: isId ? 'mutableId' : undefined,
      fieldGroup: group || 'default',
      readOnly: isId,
      defaultValue: value || undefined,
      identifier,
    };

    if (type === 'object' && !!value && outputFormat !== 'form') {
      attrDefn.defaultValue = () => {
        return value;
      };
    }

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

    const key = outputFormat === 'model' ? camelize(propName) : propName;
    attrs[key] = attrDefn;
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
    return expandOpenApiProps(props, 'form');
  }

  return {};
};

/**
 * combineOpenApiAttrs takes attributes defined on an existing models
 * and adds in the attributes found on an OpenAPI response. The values
 * defined on the model should take precedence so we can overwrite
 * attributes from OpenAPI.
 */
export const combineOpenApiAttrs = function (
  oldAttrs: Map<string, Attribute>,
  openApiProps: Record<string, MixedAttr>
) {
  const allAttrs: Record<string, boolean> = {};
  const attrsArray: Attribute[] = [];
  const newFields: string[] = [];

  // First iterate over all the existing attrs and combine with recieved props, if they exist
  oldAttrs.forEach(function (oldAttr, name) {
    const attr: Attribute = { name, type: oldAttr.type, options: oldAttr.options };
    const openApiProp = openApiProps[name];
    if (openApiProp) {
      const { type, ...options } = openApiProp;
      // TODO: previous behavior took the openApi type no matter what
      attr.type = oldAttr.type ?? type;
      if (oldAttr.type && type && type !== oldAttr.type) {
        debug(`mismatched type for ${name} -- ${type} vs ${oldAttr.type}`);
      }
      attr.options = { ...options, ...oldAttr.options };
    }
    attrsArray.push(attr);
    // add to all attrs so we skip in the next part
    allAttrs[name] = true;
  });

  // then iterate over all the new props and add them if they haven't already been accounted for
  for (const name in openApiProps) {
    // iterate over each
    if (allAttrs[name]) {
      continue;
    } else {
      const prop = openApiProps[name];
      if (prop) {
        const { type, ...options } = prop;
        newFields.push(name);
        attrsArray.push({ name, type, options });
      }
    }
  }
  return { attrs: attrsArray, newFields };
};

// interface FieldGroups {
//   default: string[];
//   [key: string]: string[];
// }

// export const combineFieldGroups = function (
//   currentGroups: Array<Record<string, string[]>>,
//   newFields: string[],
//   excludedFields: string[]
// ) {
//   console.log({ currentGroups, newFields, excludedFields });
//   let allFields: string[] = [];
//   for (const group of currentGroups) {
//     const fields = Object.values(group)[0] || [];
//     allFields = allFields.concat(fields);
//   }
//   const otherFields = newFields.filter((field) => {
//     return !allFields.includes(field) && !excludedFields.includes(field);
//   });
//   if (otherFields.length) {
//     currentGroups[0].default = currentGroups[0].default.concat(otherFields);
//   }

//   return currentGroups;
// };
