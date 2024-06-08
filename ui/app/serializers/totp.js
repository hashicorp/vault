/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeItems(payload, requestType) {
    if (
      requestType !== 'queryRecord' &&
      payload.data &&
      payload.data.keys &&
      Array.isArray(payload.data.keys)
    ) {
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
        return { id: fullSecretPath, backend: payload.backend };
      });
    }

    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },

  serialize(snapshot) {
    return {
      url: snapshot.attr('url'),
    };
  },
});
