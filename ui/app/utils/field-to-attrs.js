import { get } from '@ember/object';
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

export const expandAttributeMeta = function(modelClass, attributeNames, namePrefix, map) {
  let fields = [];
  // expand all attributes
  attributeNames.forEach(field => expandProperties(field, x => fields.push(x)));
  let attributeMap = map || new Map();
  modelClass.eachAttribute((name, meta) => {
    let fieldName = namePrefix ? namePrefix + name : name;
    let maybeFragment = get(modelClass, fieldName);
    if (meta.isFragment && maybeFragment) {
      // pass the fragment and all fields that start with
      // the fragment name down to get extracted from the Fragment
      expandAttributeMeta(
        maybeFragment,
        fields.filter(f => f.startsWith(fieldName)),
        fieldName + '.',
        attributeMap
      );
      return;
    }
    attributeMap.set(fieldName, meta);
  });

  // we have all of the attributes in the map now,
  // so we'll replace each key in `fields` with the expanded meta
  fields = fields.map(field => {
    let meta = attributeMap.get(field);
    if (meta) {
      var { type, options } = meta;
    }
    return {
      // using field name here because it is the full path,
      // name on the attribute meta will be relative to the fragment it's on
      name: field,
      type: type,
      options: options,
    };
  });
  return fields;
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

export default function(modelClass, fieldGroups) {
  return fieldGroups.map(group => {
    const groupKey = Object.keys(group)[0];
    const fields = expandAttributeMeta(modelClass, group[groupKey]);
    return { [groupKey]: fields };
  });
}
