/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiIssuerAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _getBackend(snapshot) {
    const { record, adapterOptions } = snapshot;
    return adapterOptions?.mount || record.backend;
  }

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  }

  urlForQuery(backend, id) {
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;
    if (id) {
      return `${baseUrl}/issuer/${encodePath(id)}`;
    } else {
      return `${baseUrl}/issuers`;
    }
  }

  updateRecord(store, type, snapshot) {
    const { issuerId } = snapshot.record;
    const backend = this._getBackend(snapshot);
    const data = this.serialize(snapshot);
    const url = this.urlForQuery(backend, issuerId);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    return this.ajax(this.urlForQuery(query.backend), 'GET', this.optionsForQuery());
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;
    return this.ajax(this.urlForQuery(backend, id), 'GET', this.optionsForQuery(id));
  }

  deleteAllIssuers(backend) {
    const deleteAllIssuersAndKeysUrl = `${this.buildURL()}/${encodePath(backend)}/root`;

    return this.ajax(deleteAllIssuersAndKeysUrl, 'DELETE');
  }
}
