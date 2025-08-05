/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../../application';

export default class IdentityOidcConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord() {
    return this.ajax(`${this.buildURL()}/identity/oidc/config`, 'GET').then((resp) => {
      return {
        ...resp,
        id: 'identity-oidc-config', // id required for ember data. only one record is expected so static id is fine
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
