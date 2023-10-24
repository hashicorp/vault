/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from 'vault/serializers/application';

export default class SyncDestinationSerializer extends ApplicationSerializer {
  attrs = {
    name: { serialize: false },
    type: { serialize: false },
  };

  // interrupt application's normalizeItems, which is called in normalizeResponse by application serializer
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let transformedPayload = payload;

    if (requestType === 'findRecord') {
      transformedPayload = this._normalizeFindRecord(payload);
    }
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }

  extractLazyPaginatedData(payload) {
    const transformedPayload = [];
    // loop through each destination type (keys in key_info)
    for (const type in payload.data.key_info) {
      // iterate through each type's destination names
      payload.data.key_info[type].forEach((name) => {
        const id = `${type}/${name}`;
        // create object with destination's id and attributes, add to payload
        transformedPayload.pushObject({ id, name, type });
      });
    }
    return transformedPayload;
  }

  // generates id and spreads connection_details object into data
  _normalizeFindRecord(payload) {
    if (payload?.data?.connection_details) {
      const { type, name, connection_details } = payload.data;
      const id = `${type}/${name}`;
      return { data: { id, type, name, ...connection_details } };
    }
    return payload;
  }
}
