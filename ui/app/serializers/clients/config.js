/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class ClientsConfigSerializer extends ApplicationSerializer {
  // these attrs are readOnly
  attrs = {
    billingStartTimestamp: { serialize: false },
    minimumRetentionMonths: { serialize: false },
    reportingEnabled: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (!payload.data) {
      return super.normalizeResponse(...arguments);
    }
    const normalizedPayload = {
      id: payload.id,
      data: {
        ...payload.data,
        enabled: payload.data.enabled?.includes('enable') ? 'On' : 'Off',
      },
    };
    return super.normalizeResponse(store, primaryModelClass, normalizedPayload, id, requestType);
  }

  serialize() {
    const json = super.serialize(...arguments);
    if (json.enabled === 'On' || json.enabled === 'Off') {
      const oldEnabled = json.enabled;
      json.enabled = oldEnabled === 'On' ? 'enable' : 'disable';
    }
    json.retention_months = parseInt(json.retention_months, 10);
    if (isNaN(json.retention_months)) {
      throw new Error('Invalid number value');
    }
    delete json.queries_available;
    return json;
  }
}
