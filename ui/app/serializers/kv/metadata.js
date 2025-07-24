/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { assert } from '@ember/debug';
import ApplicationSerializer from '../application';
import { kvMetadataPath } from 'vault/utils/kv-path';

export default class KvMetadataSerializer extends ApplicationSerializer {
  attrs = {
    backend: { serialize: false },
    path: { serialize: false },
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
        /**
         * ENCODING CHANGE:
         * We encode the `secret` with encodeURIComponent to preserve any trailing spaces or
         * special characters as part of the path (e.g., `test-me ` → `test-me%20`).
         * This ensures paths are safe for routing and are treated as unique in Ember Data.
         */
        const encodedSecret = encodeURIComponent(secret);

        /**
         * We also encode `payload.path` if present, so any nested paths
         * containing spaces or special characters remain valid (e.g., `beep boop/` → `beep%20boop/`).
         * If there is no payload.path, we're inside a directory.
         * We add a param called full_secret_path to the model which we use to
         * navigated to the nested secret.
         */
        const fullSecretPath = payload.path
          ? `${encodeURIComponent(payload.path)}${encodedSecret}`
          : encodedSecret;

        return {
          /**
           * ID remains encoded (e.g., `secret/metadata/test-me%20`) so Ember Data can
           * use it safely as a unique key, and Vault API requests will work without further encoding.
           */
          id: kvMetadataPath(payload.backend, fullSecretPath),

          /**
           * Store the encoded path directly. This means the model will have `test-me%20` instead of
           * `test-me `. This avoids issues with trimming/normalizing raw whitespace in Ember Data.
           */
          path: encodedSecret,

          backend: payload.backend,

          /**
           * The full secret path (including parent paths) is also encoded. This is used by the UI
           * for navigation (e.g., nested secrets).
           */
          full_secret_path: fullSecretPath,
          // Adding raw_path if we need to display a decoded version of the secret
          raw_path: secret,
        };
      });
    }
    return super.normalizeItems(payload);
  }
}
