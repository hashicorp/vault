/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer from '@ember-data/serializer/rest';
import { isNone, isBlank } from '@ember/utils';
import { decamelize } from '@ember/string';

export default RESTSerializer.extend({
  keyForAttribute: function (attr) {
    return decamelize(attr);
  },

  pushPayload(store, payload) {
    const transformedPayload = this.normalizeResponse(
      store,
      store.modelFor(payload.modelName),
      payload,
      payload.id,
      'findRecord'
    );
    return store.push(transformedPayload);
  },

  normalizeItems(payload) {
    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const responseJSON = this.normalizeItems(payload);
    const { modelName } = primaryModelClass;
    const transformedPayload = { [modelName]: responseJSON };
    const ret = this._super(store, primaryModelClass, transformedPayload, id, requestType);
    return ret;
  },

  serializeAttribute(snapshot, json, key, attributes) {
    const val = snapshot.attr(key);
    if (attributes.options.readOnly) {
      return;
    }
    if (
      attributes.type === 'object' &&
      val &&
      Object.keys(val).length > 0 &&
      isNone(snapshot.changedAttributes()[key])
    ) {
      return;
    }
    if (isBlank(val) && isNone(snapshot.changedAttributes()[key])) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },
});
