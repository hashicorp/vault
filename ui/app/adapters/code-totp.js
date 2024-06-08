/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isEmpty } from '@ember/utils';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  // TODO: Validation is not yet supported
  //createRecord(store, type, snapshot) {
  //  const serializer = store.serializerFor(type.modelName);
  //  const data = serializer.serialize(snapshot);
  //  const { id } = snapshot;
  //  const path = snapshot.record.path;
  //  return this.ajax(this.urlForCode(snapshot.attr('backend'), path || id), 'POST', { data }).then(() => {
  //    data.id = path || id;
  //    return data;
  //  });
  //},

  urlForCode(backend, id) {
    let url = `${this.buildURL()}/${encodePath(backend)}/${this.pathForType()}/`;

    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }

    return url;
  },

  pathForType() {
    return 'code';
  },

  queryRecord(store, type, query) {
    const { id, backend } = query;
    return this.ajax(this.urlForCode(backend, id), 'GET').then((resp) => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },
});
