/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class GcpConfigSerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (!payload.data) {
      return super.normalizeResponse(...arguments);
    }

    const normalizedPayload = {
      id: payload.id,
      backend: payload.backend,
      data: {
        ...payload.data,
      },
    };
    return super.normalizeResponse(store, primaryModelClass, normalizedPayload, id, requestType);
  }
}
