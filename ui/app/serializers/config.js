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
    if (payload.data) {
      const data = { ...payload, ...payload.data };
      return [data];
    }
    return [payload];
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
