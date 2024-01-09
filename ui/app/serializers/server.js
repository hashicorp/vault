/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  primaryKey: 'node_id',
  normalizeItems(payload) {
    if (payload.data && payload.data.config) {
      // rewrite the payload from data.config.servers to data.keys so we can use the application serializer
      // on it
      return payload.data.config.servers.slice(0);
    }
    return this._super(payload);
  },
});
