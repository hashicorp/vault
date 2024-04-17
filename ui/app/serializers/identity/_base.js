/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      if (typeof payload.data.keys[0] !== 'string') {
        // If keys is not an array of strings, it was already normalized into objects in extractLazyPaginatedData
        return payload.data.keys;
      }
      return payload.data.keys.map((key) => {
        const model = payload.data.key_info[key];
        model.id = key;
        return model;
      });
    }
    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },
});
