/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default class MfaLoginEnforcementSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  // Return data with updated keys for hasMany relationships with ids in the name
  transformHasManyKeys(data, destination) {
    const keys = {
      model: ['mfa_methods', 'identity_entities', 'identity_groups'],
      server: ['mfa_method_ids', 'identity_entity_ids', 'identity_group_ids'],
    };
    keys[destination].forEach((newKey, index) => {
      const oldKey = destination === 'model' ? keys.server[index] : keys.model[index];
      delete Object.assign(data, { [newKey]: data[oldKey] })[oldKey];
    });
    return data;
  }
  normalize(model, data) {
    this.transformHasManyKeys(data, 'model');
    return super.normalize(model, data);
  }
  normalizeItems(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => payload.data.key_info[key]);
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }
  serialize() {
    const json = super.serialize(...arguments);
    // empty arrays are being removed from serialized json
    // ensure that they are sent to the server, otherwise removing items will not be persisted
    json.auth_method_accessors = json.auth_method_accessors || [];
    json.auth_method_types = json.auth_method_types || [];
    // TODO: create array transform which serializes an empty array if empty
    return this.transformHasManyKeys(json, 'server');
  }
}
