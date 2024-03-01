/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from 'vault/serializers/application';
import { findDestination } from 'core/helpers/sync-destinations';

export default class SyncAssociationSerializer extends ApplicationSerializer {
  attrs = {
    destinationName: { serialize: false },
    destinationType: { serialize: false },
    syncStatus: { serialize: false },
    updatedAt: { serialize: false },
    subKey: { serialize: false },
  };

  extractLazyPaginatedData(payload) {
    if (payload) {
      const { store_name, store_type, associated_secrets } = payload.data;
      const secrets = [];
      for (const key in associated_secrets) {
        const data = associated_secrets[key];
        data.id = `${data.mount}/${data.secret_name}/${data.sub_key}`;
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

  normalizeFetchByDestinations(payload) {
    const { store_name, store_type, associated_secrets } = payload.data;
    const unsynced = [];
    let lastUpdated;

    for (const key in associated_secrets) {
      const association = associated_secrets[key];
      // for display purposes, any status other than SYNCED is considered unsynced
      if (association.sync_status !== 'SYNCED') {
        unsynced.push(association.sync_status);
      }
      // use the most recent updated_at value as the last synced date
      const updated = new Date(association.updated_at);
      if (!lastUpdated || updated > lastUpdated) {
        lastUpdated = updated;
      }
    }

    const associationCount = Object.entries(associated_secrets).length;
    return {
      icon: findDestination(store_type).icon,
      name: store_name,
      type: store_type,
      associationCount,
      status: associationCount ? (unsynced.length ? `${unsynced.length} Unsynced` : 'All synced') : null,
      lastUpdated,
    };
  }
}
