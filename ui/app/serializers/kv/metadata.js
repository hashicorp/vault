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
        // If there no payload.path then we're either on a "top level" secret or the first level directory of a nested secret.
        // Thus, we set the path to the current secret or secretPrefix.
        // We add a param called full_secret_path to the model which we use to navigate to the nested secret. e.g. beep/boop/bop.
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
