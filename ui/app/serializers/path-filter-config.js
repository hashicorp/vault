/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import RESTSerializer from '@ember-data/serializer/rest';
import { decamelize } from '@ember/string';

export default RESTSerializer.extend({
  keyForAttribute: function (attr) {
    return decamelize(attr);
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const { modelName } = primaryModelClass;
    payload.data.id = id;
    const transformedPayload = { [modelName]: payload.data };
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
