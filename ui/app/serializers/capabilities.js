/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  primaryKey: 'path',

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let response;
    // queryRecord will already have set path, and we won't have an id here
    if (id) payload.path = id;

    if (requestType === 'query') {
      // each key on the response is a path with an array of capabilities as its value
      response = Object.keys(payload.data).map((fullPath) => {
        // we use pathMap to normalize a namespace-prefixed path back to the relative path
        // this is okay because we clear capabilities when moving between namespaces
        const path = payload.pathMap ? payload.pathMap[fullPath] : fullPath;
        return { capabilities: payload.data[fullPath], path };
      });
    } else {
      response = { ...payload.data, path: payload.path };
    }
    return this._super(store, primaryModelClass, response, id, requestType);
  },

  modelNameFromPayloadKey() {
    return 'capabilities';
  },
});
