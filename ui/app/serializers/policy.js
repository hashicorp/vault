/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  primaryKey: 'name',

  normalizePolicies(payload) {
    const data = payload.data.keys ? payload.data.keys.map((name) => ({ name })) : payload.data;
    return data;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['deleteRecord'];
    const normalizedPayload = nullResponses.includes(requestType)
      ? { name: id }
      : this.normalizePolicies(payload);
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
