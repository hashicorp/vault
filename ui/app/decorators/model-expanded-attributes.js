/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import Model from '@ember-data/model';

/**
 * sets allByKey properties on model class. These are all the attributes on the model
 * and any belongsTo models, expanded for attribute metadata. The value returned is an
 * object where the key is the attribute name, and the value is the expanded attribute
 * metadata.
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
              byKey[`${name}.${attr.name}`] = attr;
            });
          }, this);
          this._allByKey = byKey;
        }
        return this._allByKey;
      }
    };
  };
}
