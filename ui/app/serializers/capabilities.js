/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  primaryKey: 'path',

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // queryRecord will already have set this, and we won't have an id here
    if (id) {
      payload.path = id;
    }
    const response = {
      ...payload.data,
      path: payload.path,
    };
    return this._super(store, primaryModelClass, response, id, requestType);
  },

  modelNameFromPayloadKey() {
    return 'capabilities';
  },
});
