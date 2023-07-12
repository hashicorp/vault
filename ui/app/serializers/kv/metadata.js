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
      return payload.data.keys.map((secret) => {
        // If there no payload.path then we're either on a "top level" secret or the first level directory of a nested secret, e.g. "beep/". Thus, we set the path to the current secret e.g. my-secret or boop/. But we add a param called full_secret_path that has all directories if it's a nested secret. e.g. beep/boop/bop.
        const fullSecretPath = payload.path ? payload.path + secret : secret;
        return {
          id: kvMetadataPath(payload.backend, fullSecretPath),
          path: secret,
          full_secret_path: fullSecretPath,
        };
      });
    }
    return super.normalizeItems(payload);
  }
}
