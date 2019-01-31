import DS from 'ember-data';
const { attr } = DS;
import { assign } from '@ember/polyfills';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export const expandOpenApiProps = function(props) {
  debugger; //eslint-disable-line
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
      editType = details.items.type + details.type.capitalize();
    }
    attrs[prop.camelize()] = {
      editType: editType,
      type: details.type,
    };
    if (details['x-vault-displayName']) {
      attrs[prop.camelize()].label = details['x-vault-displayName'];
    }
    if (details['enum']) {
      attrs[prop.camelize()].possibleValues = details['enum'];
    }
    if (details['x-vault-displayValue']) {
      attrs[prop.camelize()].defaultValue = details['x-vault-displayValue'];
    }
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
  let allFields = [];
  for (let group in currentGroups) {
    let fieldName = Object.keys(groups[group])[0];
    allFields.concat(groups[group][fieldName]);
  }
  let otherFields = newFields.filter(field => {
    !allFields.includes(field) && !excludedFields.includes(field);
  });
  if (otherFields.length) {
    groups.default.concat(otherFields);
  }
};

export const combineFieldGroups = function(currentGroups, newFields, excludedFields) {
  let allFields = [];
  for (let group of currentGroups) {
    let fieldName = Object.keys(group)[0];
    allFields = allFields.concat(group[fieldName]);
  }
  let otherFields = newFields.filter(field => {
    !allFields.includes(field) && !excludedFields.includes(field);
  });
  if (otherFields.length) {
    currentGroups.default = currentGroups.default.concat(otherFields);
  }
  return currentGroups;
};
