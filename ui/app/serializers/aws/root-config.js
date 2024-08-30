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
    // remove identityTokenTtl and maxRetries if the API's default value of 0 or -1, respectively. We don't want to display this value on configuration details if they haven't changed the default value
    if (payload.data.identity_token_ttl === 0) {
      delete payload.data.identity_token_ttl;
    }
    if (payload.data.max_retries === -1) {
      delete payload.data.max_retries;
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
