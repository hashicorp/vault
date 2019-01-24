import DS from 'ember-data';
const { attr } = DS;
import { assign } from '@ember/polyfills';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export const expandOpenApiProps = function(props) {
  let attrs = {};
  // expand all attributes
  for (let prop in props) {
    let details = props[prop];
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
