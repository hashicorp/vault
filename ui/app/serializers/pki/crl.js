/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

// The object keys ('enabled', 'duration') in the pki/crl model map to an API parameter
// Each object is updated by a single ttl input in the pki-configuration-edit form
// In some cases 'false' means enabled (i.e. { disable: false } => ttl field is toggled ON)

const MODEL_TO_API_PARAMS = {
  // object key -> model attribute name
  // "enabledKey" -> boolean api param
  // "durationKey" -> duration api param
  // "isOppositeValue" -> whether or not to flip the ttl's enabled value
  crl_expiry_data: {
    enabledKey: 'disable',
    isOppositeValue: true,
    durationKey: 'expiry',
  },
  auto_rebuild_data: {
    enabledKey: 'auto_rebuild',
    isOppositeValue: false,
    durationKey: 'auto_rebuild_grace_period',
  },
  delta_crl_building_data: {
    enabledKey: 'enable_delta',
    isOppositeValue: false,
    durationKey: 'delta_rebuild_interval',
  },
  ocsp_expiry_data: {
    enabledKey: 'ocsp_disable',
    isOppositeValue: true,
    durationKey: 'ocsp_expiry',
  },
};
export default class PkiCrlSerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    for (const key in MODEL_TO_API_PARAMS) {
      const { enabledKey, durationKey, isOppositeValue } = MODEL_TO_API_PARAMS[key];
      const valueBlock = {
        enabled: isOppositeValue ? !payload[enabledKey] : payload[enabledKey],
        duration: payload[durationKey],
      };
      payload = { ...payload, [key]: valueBlock };
    }
    return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  }

  serialize() {
    const json = super.serialize(...arguments);
    for (const key in MODEL_TO_API_PARAMS) {
      if (key in json) {
        const { enabledKey, durationKey, isOppositeValue } = MODEL_TO_API_PARAMS[key];
        const { enabled, duration } = json[key];
        json[enabledKey] = isOppositeValue ? !enabled : enabled;
        json[durationKey] = duration;
        delete json[key];
      }
    }
    return json;
  }
}
