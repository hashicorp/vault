import DS from 'ember-data';
const { attr } = DS;
import { assign } from '@ember/polyfills';
import { isEmpty } from '@ember/utils';
import { camelize, capitalize } from '@ember/string';

export const expandOpenApiProps = function(props) {
  let attrs = {};
  // expand all attributes
  for (let prop in props) {
    let details = props[prop];
    if (details.deprecated === true) {
      continue;
    }
    if (details.type === 'integer') {
      details.type = 'number';
    }
    let editType = details.type;
    if (details.format === 'seconds') {
      editType = 'ttl';
    } else if (details.items) {
      editType = details.items.type + capitalize(details.type);
    }
    let attrDefn = {
      editType: editType,
      type: details.type,
      helpText: details.description,
      sensitive: details['x-vault-displaySensitive'],
      label: details['x-vault-displayName'],
      possibleValues: details['enum'],
      defaultValue:
        details['x-vault-displayValue'] || (!isEmpty(details['default']) ? details['default'] : null),
    };
    // loop to remove empty vals
    for (let attrProp in attrDefn) {
      if (attrDefn[attrProp] == null) {
        delete attrDefn[attrProp];
      }
    }
    attrs[camelize(prop)] = attrDefn;
  }
  return attrs;
};

export const combineAttributes = function(oldAttrs, newProps) {
  let newAttrs = {};
  let newFields = [];
  oldAttrs.forEach(function(value, name) {
    if (newProps[name]) {
      newAttrs[name] = attr(newProps[name].type, assign({}, newProps[name], value.options));
    } else {
      newAttrs[name] = attr(value.type, value.options);
    }
  });
  for (let prop in newProps) {
    if (newAttrs[prop]) {
      continue;
    } else {
      newAttrs[prop] = attr(newProps[prop].type, newProps[prop]);
      newFields.push(prop);
    }
  }
  return { attrs: newAttrs, newFields };
};

export const combineFields = function(currentFields, newFields, excludedFields) {
  let otherFields = newFields.filter(field => {
    return !currentFields.includes(field) && !excludedFields.includes(field);
  });
  if (otherFields.length) {
    currentFields = currentFields.concat(otherFields);
  }
  return currentFields;
};

export const combineFieldGroups = function(currentGroups, newFields, excludedFields) {
  let allFields = [];
  for (let group of currentGroups) {
    let fieldName = Object.keys(group)[0];
    allFields = allFields.concat(group[fieldName]);
  }
  let otherFields = newFields.filter(field => {
    return !allFields.includes(field) && !excludedFields.includes(field);
  });
  if (otherFields.length) {
    currentGroups[0].default = currentGroups[0].default.concat(otherFields);
  }

  return currentGroups;
};
