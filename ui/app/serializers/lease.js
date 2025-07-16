/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer from '@ember-data/serializer/rest';
import { decamelize } from '@ember/string';

export default RESTSerializer.extend({
  keyForAttribute: function (attr) {
    return decamelize(attr);
  },

  normalizeAll(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const records = payload.data.keys.map((record) => {
        const fullPath = payload.prefix ? payload.prefix + record : record;
        return { id: fullPath };
      });
      return records;
    }
    return [payload.data];
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const records = this.normalizeAll(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: records };
    // just return the single object because ember is picky
    if (requestType === 'queryRecord') {
      transformedPayload = { [modelName]: records[0] };
    }

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
