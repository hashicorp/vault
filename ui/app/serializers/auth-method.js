/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from './application';
import { EmbeddedRecordsMixin } from '@ember-data/serializer/rest';

export default ApplicationSerializer.extend(EmbeddedRecordsMixin, {
  attrs: {
    config: { embedded: 'always' },
  },
  normalize(modelClass, data) {
    // embedded records need a unique value to be stored
    // use the uuid from the auth-method as the unique id for mount-config
    if (data.config && !data.config.id) {
      data.config.id = data.uuid;
    }
    return this._super(modelClass, data);
  },
  normalizeBackend(path, backend) {
    const struct = { ...backend };
    // strip the trailing slash off of the path so we
    // can navigate to it without getting `//` in the url
    struct.id = path.slice(0, -1);
    struct.path = path;
    return struct;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const isCreate = requestType === 'createRecord';
    const backends = isCreate
      ? payload.data
      : Object.keys(payload.data).map((path) => this.normalizeBackend(path, payload.data[path]));

    return this._super(store, primaryModelClass, backends, id, requestType);
  },
});
