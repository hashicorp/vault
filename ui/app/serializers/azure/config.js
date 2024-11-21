/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class AzureConfigSerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (!payload.data) {
      return super.normalizeResponse(...arguments);
    }
    // remove rootPasswordTtl and identityTokenT if the API's default value of 0. We don't want to display this value on configuration details if they haven't changed the default value
    if (payload.data.root_password_ttl === 0) {
      delete payload.data.root_password_ttl;
    }
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
