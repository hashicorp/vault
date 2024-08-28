/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default class AwsRootConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord() {
    return this.ajax(`${this.buildURL()}/identity/oidc/config`, 'GET').then((resp) => {
      return {
        ...resp,
        id: resp.data.issuer, // id required for ember data
      };
    });
  }

  createOrUpdate(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    return this.ajax(`${this.buildURL()}/identity/oidc/config`, 'POST', { data }).then((resp) => {
      // id is returned from API so we do not need to explicitly set it here
      return {
        ...resp,
      };
    });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments);
  }
}
