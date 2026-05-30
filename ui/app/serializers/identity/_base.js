/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeItems(payload) {
    // Extract the keys array into individual objects with the key as the id, if it exists.
    // This is for endpoints that return a list of keys and a separate key_info object with
    // the details for each key, such as the list endpoint for entities.
    if (payload.data?.keys && Array.isArray(payload.data?.keys)) {
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
