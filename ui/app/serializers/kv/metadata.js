/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { assert } from '@ember/debug';
import ApplicationSerializer from '../application';
import { kvMetadataPath } from 'vault/utils/kv-path';

export default class KvMetadataSerializer extends ApplicationSerializer {
  attrs = {
    oldestVersion: { serialize: false },
    createdTime: { serialize: false },
    updatedTime: { serialize: false },
    currentVersion: { serialize: false },
    versions: { serialize: false },
  };

  normalizeItems(payload) {
    if (payload.data.keys) {
      assert('payload.backend must be provided on kv/metadata list response', !!payload.backend);
      const backend = payload.backend;
      const arrayOfKeyObjects = payload.data.keys.map((path) => {
        return {
          id: kvMetadataPath(backend, path),
          backend,
          path,
        };
      });
      return arrayOfKeyObjects;
    }
    return super.normalizeItems(payload);
  }
}
