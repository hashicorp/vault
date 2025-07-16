/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer from '@ember-data/serializer/rest';
import { decamelize } from '@ember/string';

export default RESTSerializer.extend({
  primaryKey: 'name',

  keyForAttribute: function (attr) {
    return decamelize(attr);
  },

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const secrets = payload.data.keys.map((secret) => ({ name: secret, backend: payload.backend }));
      return secrets;
    }
    Object.assign(payload, payload.data);
    delete payload.data;
    // timestamps for these two are in seconds...
    if (
      payload.type === 'aes256-gcm96' ||
      payload.type === 'chacha20-poly1305' ||
      payload.type === 'aes128-gcm96'
    ) {
      for (const version in payload.keys) {
        payload.keys[version] = payload.keys[version] * 1000;
      }
    }
    payload.id = payload.name;
    return [payload];
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const secrets = nullResponses.includes(requestType)
      ? { name: id, backend: payload.backend }
      : this.normalizeSecrets(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: secrets };
    // just return the single object because ember is picky
    if (requestType === 'queryRecord') {
      transformedPayload = { [modelName]: secrets[0] };
    }

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serialize(snapshot, requestType) {
    if (requestType === 'update') {
      const min_decryption_version = snapshot.attr('minDecryptionVersion');
      const min_encryption_version = snapshot.attr('minEncryptionVersion');
      const deletion_allowed = snapshot.attr('deletionAllowed');
      const auto_rotate_period = snapshot.attr('autoRotatePeriod');
      return {
        min_decryption_version,
        min_encryption_version,
        deletion_allowed,
        auto_rotate_period,
      };
    } else {
      snapshot.id = snapshot.attr('name');
      return this._super(snapshot, requestType);
    }
  },
});
