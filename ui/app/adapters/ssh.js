/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { assert } from '@ember/debug';
import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  url(/*role*/) {
    assert('Override the `url` method to extend the SSH adapter', false);
  },

  createRecord(store, type, snapshot, requestType) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot, requestType);
    const role = snapshot.attr('role');

    return this.ajax(this.url(role), 'POST', { data }).then((response) => {
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },
});
