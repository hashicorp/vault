/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expandProperties } from '@ember/object/computed';
/*
 *
 * @param modelClass DS.Model
 * @param attributeNames Array[String]
 * @param prefixName String
 * @param map Map
 * @returns Array[Object]
 *
 * A function that takes a model and an array of attributes
 * and expands them in-place to an array of metadata about the attributes
 *
 * if passed a Model with attributes `foo` and `bar` and the array ['foo', 'bar']
 * the returned array would take the form of:
 *
 *  [
 *    {
 *      name: 'foo',
 *      type: 'string',
 *      options: {
 *        defaultValue: 'Foo'
 *      }
 *    },
 *    {
 *      name: 'bar',
 *      type: 'string',
 *      options: {
 *        defaultValue: 'Bar',
 *        editType: 'textarea',
 *        label: 'The Bar Field'
 *      }
 *    },
 *  ]
 *
 */

export const expandAttributeMeta = function (modelClass, attributeNames) {
  const fields = [];
  // expand all attributes
  // see docs for examples of brace-expansion - https://api.emberjs.com/ember/4.4/functions/@ember%2Fobject%2Fcomputed/expandProperties
  attributeNames.map((field) => expandProperties(field, (prop) => fields.push(prop)));
  // cache results of eachAttribute so we don't call it on each iteration of fields loop
  const modelAttrs = {};

  const getAttributeMeta = (klass, attrKey) => {
    // populate cache if empty
    if (!modelAttrs[klass.modelName]) {
      modelAttrs[klass.modelName] = [];
      klass.eachAttribute((name, meta) => {
        modelAttrs[klass.modelName].push(meta);
      });
    }
    // lookup attr and return meta
    return modelAttrs[klass.modelName].find((attr) => attr.name === attrKey);
  };

  return fields.map((field) => {
    let meta = {};
    // check for relationship by presence of dot nation in field name
    if (field.includes('.')) {
      const [relKey, prop] = field.split('.');
      const rel = modelClass.belongsTo(relKey);
      const relModelClass = modelClass.store.modelFor(rel.type);
      meta = getAttributeMeta(relModelClass, prop);
    } else {
      meta = getAttributeMeta(modelClass, field);
    }
    const { type, options } = meta || {};
    return {
      // using field name here because it is the full path,
      // name on the attribute meta will be relative to the relationship if applicable
      name: field,
      type,
      options,
    };
  });
};

/*
 *
 * @param modelClass DS.Model
 * @param fieldGroups Array[Object]
 * @returns Array
 *
 * A function meant for use on an Ember Data Model
 *
 * The function takes a array of groups, each group
 * being a list of attributes on the model, for example
 * `fieldGroups` could look like this
 *
 *  [
 *    { default: ['commonName', 'format'] },
 *    { Options: ['altNames', 'ipSans', 'ttl', 'excludeCnFromSans'] },
 *  ]
 *
 *  The array will get mapped over producing a new array with each attribute replaced with that attribute's metadata from the attr declaration
 */

export default function (modelClass, fieldGroups) {
  return fieldGroups.map((group) => {
    const groupKey = Object.keys(group)[0];
    const fields = expandAttributeMeta(modelClass, group[groupKey]);
    return { [groupKey]: fields };
  });
}
