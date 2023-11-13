/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import IdentitySerializer from './_base';
export default IdentitySerializer.extend({
  extractLazyPaginatedData(payload) {
    return payload.data.keys.map((key) => {
      const model = {
        id: key,
        name: payload.data.key_info[key].name,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
  },
  normalizeItems(payload) {
    // due to extractLazyPaginatedData above, keys is no longer an array of strings but objects
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys;
    }
    const ret = {
      ...payload,
      ...payload.data,
    };
    delete ret.data;
    return ret;
  },
});
