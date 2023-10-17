/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from 'vault/serializers/application';

export default class SyncDestinationSerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let transformedPayload = payload;

    switch (requestType) {
      case 'query':
        transformedPayload = this._transformQueryPayload(payload);
        break;
      case 'findRecord':
        transformedPayload = this._transformFindRecordPayload(payload);
        break;
      default:
        return super.normalizeResponse(...arguments);
    }
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }

  _transformQueryPayload(payload) {
    if (payload.data?.keys && Array.isArray(payload.data.keys)) {
      // return array of destination objects
      return payload.data.keys.map((id) => ({
        id,
        ...payload.data.key_info[id],
      }));
    }
    return payload;
  }

  _transformFindRecordPayload(payload) {
    if (payload?.data?.connection_details) {
      const { type, name, connection_details } = payload.data;
      const id = `${type}/${name}`;
      return { data: { id, type, name, ...connection_details } };
    }
    return payload;
  }
}
