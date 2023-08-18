/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KvDataSerializer extends ApplicationSerializer {
  serialize(snapshot) {
    const { secretData, casVersion } = snapshot.record;
    if (typeof casVersion === 'number') {
      /* if this is a number it is set by one of the following:
        A) user is creating initial version of a secret
         -> 0 : default value set in route
        B) user is creating a new version of a secret:
         -> metadata.current_version : has metadata read permissions (data permissions are irrelevant)
         -> secret.version : has data read permissions. without metadata read access a user is unable to navigate,
                             to older secret versions so we assume creation is from the latest version */
      return { data: secretData, options: { cas: casVersion } };
    }
    // a non-number value means no read permission for both data and metadata
    return { data: secretData };
  }

  normalizeKvData(payload) {
    const { data, metadata } = payload.data;
    return {
      ...payload,
      data: {
        ...payload.data,
        // Rename to secret_data so it doesn't get removed by normalizer
        secret_data: data,
        ...metadata,
      },
    };
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (requestType === 'queryRecord') {
      const transformed = this.normalizeKvData(payload);
      return super.normalizeResponse(store, primaryModelClass, transformed, id, requestType);
    }
    return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  }
}
