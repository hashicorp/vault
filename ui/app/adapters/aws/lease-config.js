/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class AwsLeaseConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord(store, type, query) {
    const { backend } = query;
    try {
      return this.ajax(`${this.buildURL()}/${encodePath(backend)}/config/lease`, 'GET').then((resp) => {
        resp.id = backend;
        return resp;
      });
    } catch (error) {
      if (error.httpStatus !== 404) {
        // not found error occurs when no lease config is set.
        // we still want to return the model for mapping purposes, so bypass error here.
        // ARG TODO not doing anything
        return [];
      }
      return error;
    }
  }

  createOrUpdate(store, type, snapshot) {
    const { data } = snapshot.adapterOptions;
    const path = encodePath(snapshot.id);
    return this.ajax(`/v1/${path}/config/lease`, 'POST', { data });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments);
  }
}
