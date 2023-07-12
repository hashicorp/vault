/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// import { assert } from '@ember/debug';
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

  // normalizeItems(payload) {
  //   if (payload.data.keys) {
  //     assert('payload.backend must be provided on kv/metadata list response', !!payload.backend);
  //     const backend = payload.backend;
  //     const arrayOfKeyObjects = payload.data.keys.map((path) => {
  //       // handling list-root and list for nested secrets. need the beep/boop/bop on the model
  //       return {
  //         id: kvMetadataPath(backend, path),
  //         backend,
  //         path,
  //       };
  //     });
  //     return arrayOfKeyObjects;
  //   }
  //   return super.normalizeItems(payload);
  // }
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      // if we have data.keys, it's a list of ids, so we map over that
      // and create objects with id's
      return payload.data.keys.map((secret) => {
        // no path? then we're either on a "top level" secret or the first of a directory of a nested secret, e.g. "beep/". Set the path to either "my-secret" or "beep/". If we do have a path then we're in a nested secret. Concat to turn "beep/" into "beep/boop/"
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
