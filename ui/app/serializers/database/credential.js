/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import RESTSerializer from '@ember-data/serializer/rest';

export default RESTSerializer.extend({
  primaryKey: 'username',

  normalizePayload(payload) {
    if (payload.data) {
      return {
        username: payload.data.username,
        password: payload.data.password,
        leaseId: payload.lease_id,
        leaseDuration: payload.lease_duration,
        lastVaultRotation: payload.data.last_vault_rotation,
        rotationPeriod: payload.data.rotation_period,
        ttl: payload.data.ttl,
        // roleType is added on adapter
        roleType: payload.roleType,
      };
    }
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const credentials = this.normalizePayload(payload);
    const { modelName } = primaryModelClass;
    const transformedPayload = { [modelName]: credentials };

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
