/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// import { assert } from '@ember/debug';
import ApplicationSerializer from '../application';
// import { kvMetadataPath } from 'vault/utils/kv-path';

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
        // secrets don't have an id in the response, so we need to concat the full
        // path of the secret here - the id in the payload is added
        // in the adapter after making the request
        let fullSecretPath = payload.id ? payload.id + secret : secret;
        // if there is no path, it's a "top level" secret, so add
        // a unicode space for the id
        // https://github.com/hashicorp/vault/issues/3348
        if (!fullSecretPath) {
          fullSecretPath = '\u0020';
        }
        return {
          // id: kvMetadataPath(payload.backend, fullSecretPath),
          id: fullSecretPath,
          path: secret,
        };
      });
    }
    return super.normalizeItems(payload);
  }
}
