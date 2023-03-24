/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { belongsTo } from '@ember-data/model';
import { assert, debug } from '@ember/debug';
import { typeOf } from '@ember/utils';
import { isArray } from '@ember/array';

/*
 *
 * attachCapabilities
 *
 * @param modelClass = An Ember Data model class
 * @param capabilities - an Object whose keys will added to the model class as related 'capabilities' models
 * and whose values should be functions that return the id of the related capabilities model
 *
 * definition of capabilities be done shorthand with the apiPath tagged template function
 *
 *
 * @usage
 *
 * let Model = DS.Model.extend({
 *   backend: attr(),
 *   scope: attr(),
 * });
 *
 * export default attachCapabilities(Model, {
 *   updatePath: apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`,
 * });
 *
 */
export default function attachCapabilities(modelClass, capabilities) {
  const capabilityKeys = Object.keys(capabilities);
  const newRelationships = capabilityKeys.reduce((ret, key) => {
    ret[key] = belongsTo('capabilities');
    return ret;
  }, {});

  //TODO: move this to the application serializer and do it JIT instead of on app boot
  debug(`adding new relationships: ${capabilityKeys.join(', ')} to ${modelClass.toString()}`);
  modelClass.reopen(newRelationships);
  modelClass.reopenClass({
    // relatedCapabilities is called in the application serializer's
    // normalizeResponse hook to add the capabilities relationships to the
    // JSON-API document used by Ember Data
    relatedCapabilities(jsonAPIDoc) {
      let { data, included } = jsonAPIDoc;
      if (!data) {
        data = jsonAPIDoc;
      }
      if (isArray(data)) {
        const newData = data.map(this.relatedCapabilities);
        return {
          data: newData,
          included,
        };
      }
      const context = {
        id: data.id,
        ...data.attributes,
      };
      for (const newCapability of capabilityKeys) {
        const templateFn = capabilities[newCapability];
        const type = typeOf(templateFn);
        assert(`expected value of ${newCapability} to be a function but found ${type}.`, type === 'function');
        data.relationships[newCapability] = {
          data: {
            type: 'capabilities',
            id: templateFn(context),
          },
        };
      }

      if (included) {
        return {
          data,
          included,
        };
      } else {
        return data;
      }
    },
  });
  return modelClass;
}
