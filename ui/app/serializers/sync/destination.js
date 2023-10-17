/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from 'vault/serializers/application';

export default class SyncDestinationSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  normalizeItems(payload) {
    if (payload.data?.connection_details) {
      const data = { ...payload.data, ...payload.data.connection_details };
      delete data.connection_details;
      return data;
    }
    return super.normalizeItems(...arguments);
  }
}
