/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const normalizedPayload = {
      id: payload.id,
      status: payload.data,
    };

    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
