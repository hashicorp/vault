/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import Model from '@ember-data/model';
import { debug } from '@ember/debug';

/**
 * sets allByKey properties on model class. These are all the attributes on the model
 * and any belongsTo models, expanded with attribute metadata. The value returned is an
 * object where the key is the attribute name, and the value is the expanded attribute
 * metadata.
 * This decorator also exposes a helper function `_expandGroups` which, when given groups
 * as expected in field-to-attrs util, will return a similar object with the expanded
 * attributes in place of the strings in the array.
 */

export function withExpandedAttributes() {
  return function decorator(SuperClass) {
    if (!Object.prototype.isPrototypeOf.call(Model, SuperClass)) {
      // eslint-disable-next-line
      console.error(
        'withExpandedAttributes decorator must be used on instance of ember-data Model class. Decorator not applied to returned class'
      );
      return SuperClass;
    }
    return class ModelExpandedAttrs extends SuperClass {
      // Helper method for expanding dynamic groups on model
      _expandGroups(groups) {
        if (!Array.isArray(groups)) {
          throw new Error('_expandGroups expects an array of objects');
        }
        /* Expects group shape to be something like:
        [
          { default: ['ttl', 'maxTtl'] },
          { "Method Options": ['other', 'fieldNames'] },
        ]*/
        return groups.map((obj) => {
          const [key, stringArray] = Object.entries(obj)[0];
          const expanded = stringArray.map((fieldName) => this.allByKey[fieldName]).filter((f) => !!f);
          // if this fails, it might mean there are missing fields in the model or the model must be hydrated via OpenAPI
          if (expanded.length !== stringArray.length) {
            debug(`not all model fields found in allByKey for group "${key}"`);
          }
          return { [key]: expanded };
        });
      }

      _allByKey = null;
      get allByKey() {
        // Caching like this ensures allByKey only gets calculated once
        if (!this._allByKey) {
          const byKey = {};
          // First, get attr names which are on the model directly
          // By this time, OpenAPI should have populated non-explicit attrs
          const mainFields = [];
          this.eachAttribute(function (key) {
            mainFields.push(key);
          });
          const expanded = expandAttributeMeta(this, mainFields);
          expanded.forEach((attr) => {
            // Add expanded attributes from the model
            byKey[attr.name] = attr;
          });

          // Next, fetch and expand attrs for related models
          this.eachRelationship(function (name, descriptor) {
            // We don't worry about getting hasMany relationships
            if (descriptor.kind !== 'belongsTo') return;
            const rModel = this[name];
            const rAttrNames = [];
            rModel.eachAttribute(function (key) {
              rAttrNames.push(key);
            });
            const expanded = expandAttributeMeta(rModel, rAttrNames);
            expanded.forEach((attr) => {
              byKey[`${name}.${attr.name}`] = {
                ...attr,
                options: {
                  ...attr.options,
                  // This ensures the correct path is updated in FormField
                  fieldValue: `${name}.${attr.fieldValue || attr.name}`,
                },
              };
            });
          }, this);
          this._allByKey = byKey;
        }
        return this._allByKey;
      }
    };
  };
}
