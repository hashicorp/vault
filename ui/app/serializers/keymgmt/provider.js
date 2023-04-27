/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KeymgmtProviderSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  normalizeItems(payload) {
    const normalized = super.normalizeItems(payload);
    if (Array.isArray(normalized)) {
      normalized.forEach((provider) => {
        provider.id = provider.name;
        provider.backend = payload.backend;
      });
    }
    return normalized;
  }

  serialize(snapshot) {
    const json = super.serialize(...arguments);
    return {
      ...json,
      credentials: snapshot.record.credentials,
    };
  }
}
