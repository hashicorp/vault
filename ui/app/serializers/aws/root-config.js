/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class AwsRootConfigSerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (!payload.data) {
      return super.normalizeResponse(...arguments);
    }
    // remove IdentityTokenTtl if default value of 0 is returned. We don't want to display this value on configuration details
    if (payload.data.identity_token_ttl === 0) {
      delete payload.data.identity_token_ttl;
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
