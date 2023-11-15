/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from 'vault/serializers/application';

export default class SyncAssociationSerializer extends ApplicationSerializer {
  attrs = {
    destinationName: { serialize: false },
    destinationType: { serialize: false },
    syncStatus: { serialize: false },
    updatedAt: { serialize: false },
  };

  extractLazyPaginatedData(payload) {
    if (payload) {
      const { store_name, store_type, associated_secrets } = payload.data;
      const secrets = [];
      for (const key in associated_secrets) {
        const data = associated_secrets[key];
        data.id = key;
        const association = {
          destinationName: store_name,
          destinationType: store_type,
          ...data,
        };
        secrets.push(association);
      }
      return secrets;
    }
    return payload;
  }
}
