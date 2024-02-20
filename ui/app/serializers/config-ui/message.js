/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { decodeString, encodeString } from 'core/utils/b64';
import ApplicationSerializer from '../application';

export default class MessageSerializer extends ApplicationSerializer {
  attrs = {
    active: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (requestType === 'query' && !payload.meta) {
      const transformed = this.mapPayload(payload);
      return super.normalizeResponse(store, primaryModelClass, transformed, id, requestType);
    }
    if (requestType === 'queryRecord') {
      const transformed = {
        ...payload.data,
        message: decodeString(payload.data.message),
      };
      return super.normalizeResponse(store, primaryModelClass, transformed, id, requestType);
    }
    return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  }

  serialize() {
    const json = super.serialize(...arguments);
    json.message = encodeString(json.message);
    return json;
  }

  mapPayload(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => {
          const data = {
            id: key,
            ...payload.data.key_info[key],
          };
          if (data.message) data.message = decodeString(data.message);
          return data;
        });
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }

  extractLazyPaginatedData(payload) {
    return this.mapPayload(payload);
  }
}
