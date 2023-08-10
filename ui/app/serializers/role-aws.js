/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  primaryKey: 'name',

  extractLazyPaginatedData(payload) {
    return payload.data.keys.map((key) => {
      const model = {
        name: key,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
  },

  normalizeItems() {
    const normalized = this._super(...arguments);
    // most roles will only have one in this array,
    // we'll default to the first, and keep the array on the
    // model and show a warning if there's more than one so that
    // they don't inadvertently save
    if (normalized.credential_types) {
      normalized.credential_type = normalized.credential_types[0];
    }
    return normalized;
  },
});
