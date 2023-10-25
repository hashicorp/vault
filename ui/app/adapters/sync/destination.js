/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';

export default class SyncDestinationAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _baseUrl() {
    return `${this.buildURL()}/sys`;
  }

  createRecord(store, type, snapshot) {
    const { name, type: destinationType } = snapshot.attributes();
    const url = `${this._baseUrl()}/sync/destinations/${destinationType}/${name}`;

    return this.ajax(url, 'POST', { data: snapshot.serialize() }).then((resp) => ({
      id: `${destinationType}/${name}`,
      ...resp,
    }));
  }

  // modelName is sync/destinations/:type
  // id is the destination name
  urlForFindRecord(id, modelName) {
    return `${this._baseUrl()}/${modelName}/${id}`;
  }

  query() {
    const url = `${this._baseUrl()}/sync/destinations`;
    return this.ajax(url, 'GET', { data: { list: true } });
  }
}
