/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { camelize, capitalize } from '@ember/string';

export const expandOpenApiProps = function (props) {
  const attrs = {};
  // expand all attributes
  for (const propName in props) {
    const prop = props[propName];
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

    const attrDefn = {
      editType,
      helpText: description,
      possibleValues: prop['enum'],
      fieldValue: isId ? 'mutableId' : null,
      fieldGroup: group || 'default',
      readOnly: isId,
      defaultValue: value || null,
    };

    if (type === 'object' && !!value) {
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
    attrs[camelize(propName)] = attrDefn;
  }
  return attrs;
};

/**
 * combineFieldGroups takes the newFields returned from OpenAPI and adds them to the default field group
 * if they are not already accounted for in other field groups
 * @param {Record<string,string[]>[]} currentGroups Field groups, as an array of objects like: [{ default: [] }, { 'TLS options': [] }]
 * @param {string[]} newFields
 * @param {string[]} excludedFields
 * @returns modified currentGroups
 */
export const combineFieldGroups = function (currentGroups, newFields, excludedFields) {
  let allFields = [];
  for (const group of currentGroups) {
    const fieldName = Object.keys(group)[0];
    allFields = allFields.concat(group[fieldName]);
  }
  const otherFields = newFields.filter((field) => {
    return !allFields.includes(field) && !excludedFields.includes(field);
  });
  if (otherFields.length) {
    currentGroups[0].default = currentGroups[0].default.concat(otherFields);
  }

  return currentGroups;
};
