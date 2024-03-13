/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import JSONSerializer from '@ember-data/serializer/json';
import { isNone, isBlank } from '@ember/utils';
import { decamelize } from '@ember/string';

export default JSONSerializer.extend({
  keyForAttribute: function (attr) {
    return decamelize(attr);
  },

  normalizeItems(payload) {
    if (payload.data && payload.data.keys && Array.isArray(payload.data.keys)) {
      const models = payload.data.keys.map((key) => {
        if (typeof key !== 'string') {
          return key;
        }
        const pk = this.primaryKey || 'id';
        let model = { [pk]: key };
        // if we've added _requestQuery in the adapter, we want
        // attach it to the individual models
        if (payload._requestQuery) {
          model = { ...model, ...payload._requestQuery };
        }
        return model;
      });
      return models;
    }
    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
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

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const responseJSON = this.normalizeItems(payload, requestType);
    delete payload._requestQuery;
    if (id && !responseJSON.id) {
      responseJSON.id = id;
    }
    let jsonAPIRepresentation = this._super(store, primaryModelClass, responseJSON, id, requestType);
    if (primaryModelClass.relatedCapabilities) {
      jsonAPIRepresentation = primaryModelClass.relatedCapabilities(jsonAPIRepresentation);
    }
    return jsonAPIRepresentation;
  },

  serializeAttribute(snapshot, json, key, attributes) {
    const val = snapshot.attr(key);
    const valHasNotChanged = isNone(snapshot.changedAttributes()[key]);
    const valIsBlank = isBlank(val);
    if (attributes.options.readOnly) {
      return;
    }
    if (valIsBlank && valHasNotChanged) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },

  serializeBelongsTo(snapshot, json) {
    return json;
  },
});
