/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

// The object keys ('enabled', 'duration') in the pki/crl model map to an API parameter
// Each object is updated by a single ttl input in the pki-configuration-edit form
// In some cases 'false' means enabled (i.e. { disable: false } => ttl field is toggled ON)

const MODEL_TO_API_PARAMS = ['crl_expiry_data', 'delta_crl_building_data', 'ocsp_expiry_data'];
export default class PkiCrlSerializer extends ApplicationSerializer {
  // normalizeResponse(store, primaryModelClass, payload, id, requestType) {
  //   for (const key in MODEL_TO_API_PARAMS) {
  //     const { enabledKey, durationKey, isOppositeValue } = MODEL_TO_API_PARAMS[key];
  //     const valueBlock = {
  //       enabled: isOppositeValue ? !payload[enabledKey] : payload[enabledKey],
  //       duration: payload[durationKey],
  //     };
  //     payload = { ...payload, [key]: valueBlock };
  //   }
  //   return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  // }
  serialize() {
    let json = super.serialize(...arguments);
    MODEL_TO_API_PARAMS.forEach((key) => {
      if (key in json) {
        json = { ...json, ...json[key] };
        delete json[key];
      }
    });
    return json;
  }
}
