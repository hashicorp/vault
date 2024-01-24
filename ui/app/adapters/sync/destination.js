/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';
import { pluralize } from 'ember-inflector';

export default class SyncDestinationAdapter extends ApplicationAdapter {
  namespace = 'v1/sys';

  pathForType(modelName) {
    return modelName === 'sync/destination' ? pluralize(modelName) : modelName;
  }

  urlForCreateRecord(modelName, snapshot) {
    const { name } = snapshot.attributes();
    return `${super.urlForCreateRecord(modelName, snapshot)}/${name}`;
  }

  updateRecord(store, { modelName }, snapshot) {
    const { name } = snapshot.attributes();
    return this.ajax(`${this.buildURL(modelName)}/${name}`, 'PATCH', { data: snapshot.serialize() });
  }

  urlForDeleteRecord(id, modelName, snapshot) {
    const { name, type } = snapshot.attributes();
    // the only delete option in the UI is to purge which unsyncs all secrets prior to deleting
    return `${this.buildURL('sync/destinations')}/${type}/${name}?purge=true`;
  }

  query(store, { modelName }) {
    return this.ajax(this.buildURL(modelName), 'GET', { data: { list: true } });
  }

  createRecord(store, type, snapshot) {
    const id = `${snapshot.record.type}:${snapshot.record.name}`;
    return super.createRecord(...arguments).then((resp) => {
      resp.id = id;
      return resp;
    });
  }

  // return normalized query response
  // useful for fetching data directly without loading models into store
  async normalizedQuery() {
    const queryResponse = await this.query(this.store, { modelName: 'sync/destination' });
    const serializer = this.store.serializerFor('sync/destination');
    return serializer.extractLazyPaginatedData(queryResponse);
  }
}
