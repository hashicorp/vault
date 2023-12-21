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
  };

  serialize(snapshot) {
    const data = super.serialize(snapshot);

    if (snapshot.isNew) return data;
    // only send changed parameters for PATCH requests
    const patchData = {};
    const changedKeys = Object.keys(snapshot.changedAttributes()).map((key) => decamelize(key));
    changedKeys.forEach((key) => (patchData[key] = data[key]));
    return patchData;
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
      const connection_details = payload.data.connection_details || {};
      data.id = data.name;
      delete data.connection_details;
      return { data: { ...data, ...connection_details } };
    }
    return payload;
  }
}
