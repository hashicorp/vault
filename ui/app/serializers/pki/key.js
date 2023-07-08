/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class PkiKeySerializer extends ApplicationSerializer {
  primaryKey = 'key_id';
  attrs = {
    type: { serialize: false },
  };

  // rehydrate each keys model so all model attributes are accessible from the LIST response
  normalizeItems(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => ({ key_id: key, ...payload.data.key_info[key] }));
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }
}
