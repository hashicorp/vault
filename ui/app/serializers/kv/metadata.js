/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { assert } from '@ember/debug';
import { kvMetadataPath } from 'vault/utils/kv-path';
import ApplicationSerializer from '../application';

export default class KvMetadataSerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    const newPayload = { ...payload };
    if (payload.data.keys) {
      assert('payload.backend must be provided on kv/metadata list response', !!payload.backend);
      const backend = payload.backend;
      newPayload.data.keys = payload.data.keys.map((path) => {
        return {
          id: kvMetadataPath(backend, path),
          backend,
          path,
        };
      });
    }
    return newPayload;
  }
}
