/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeItems(payload, requestType) {
    if (
      requestType !== 'queryRecord' &&
      payload.data &&
      payload.data.keys &&
      Array.isArray(payload.data.keys)
    ) {
      // if we have data.keys, it's a list of ids, so we map over that
      // and create objects with id's
      return payload.data.keys.map((secret) => ({
        id: secret,
        backend: payload.backend,
      }));
    }

    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },
});
