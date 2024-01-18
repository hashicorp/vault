/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from 'vault/serializers/application';
import { decamelize } from '@ember/string';
export default class SyncDestinationSerializer extends ApplicationSerializer {
  attrs = {
    name: { serialize: false },
    type: { serialize: false },
    purgeInitiatedAt: { serialize: false },
    purgeError: { serialize: false },
  };

  serialize(snapshot) {
    const data = super.serialize(snapshot);
    if (snapshot.isNew) return data;

    // only send changed parameters for PATCH requests
    const changedKeys = Object.keys(snapshot.changedAttributes()).map((key) => decamelize(key));
    return changedKeys.reduce((payload, key) => {
      if (JSON.stringify(data[key]) === '{}') {
        // sending an empty object won't clear the previous param, set to null so PATCH removes pre-existing value
        payload[key] = null;
      } else {
        payload[key] = data[key];
      }
      return payload;
    }, {});
  }

  // interrupt application's normalizeItems, which is called in normalizeResponse by application serializer
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const transformedPayload = this._normalizePayload(payload, requestType);
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }

  extractLazyPaginatedData(payload) {
    const transformedPayload = [];
    // loop through each destination type (keys in key_info)
    for (const key in payload.data.key_info) {
      // iterate through each type's destination names
      payload.data.key_info[key].forEach((name) => {
        // remove trailing slash from key
        const type = key.replace(/\/$/, '');
        const id = `${type}/${name}`;
        // create object with destination's id and attributes, add to payload
        transformedPayload.pushObject({ id, name, type });
      });
    }
    return transformedPayload;
  }

  _normalizePayload(payload, requestType) {
    // if request is from lazyPaginatedQuery it will already have been extracted and meta will be set on object
    // for store.query it will be the raw response which will need to be extracted
    if (requestType === 'query') {
      return payload.meta ? payload : this.extractLazyPaginatedData(payload);
    } else if (payload?.data) {
      // uses name for id and spreads connection_details object into data
      const { data } = payload;
      const { connection_details, options } = data;
      data.id = data.name;
      delete data.connection_details;
      delete data.options;
      return { data: { ...data, ...connection_details, ...options } };
    }
    return payload;
  }
}
