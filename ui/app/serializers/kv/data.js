/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KvDataSerializer extends ApplicationSerializer {
  serialize(snapshot) {
    // Regardless of if CAS === true on the kv mount, the UI always sends the "options" object with the cas version.
    return {
      data: snapshot.attr('secretData'),
      options: {
        cas: snapshot.record.casVersion,
      },
    };
  }

  normalizeKvData(payload) {
    const { id, backend, path, data, metadata } = payload.data;
    return {
      ...payload,
      data: {
        id,
        backend,
        path,
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
