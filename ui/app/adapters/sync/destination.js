/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';
import { pluralize } from 'ember-inflector';
import { decamelize } from '@ember/string';

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
    // only send changed values
    const data = {};
    for (const attr in snapshot.changedAttributes()) {
      // first array element is the old value
      const [, newValue] = snapshot.changedAttributes()[attr];
      data[decamelize(attr)] = newValue;
    }
    // TODO come back to
    // changed attributes doesn't track arrays, manually add attr
    if (Object.keys(snapshot.serialize()).includes('deployment_environments')) {
      data['deployment_environments'] = snapshot.serialize()['deployment_environments'];
    }
    return this.ajax(`${this.buildURL(modelName)}/${name}`, 'PATCH', { data });
  }

  urlForDeleteRecord(id, modelName, snapshot) {
    const { name, type } = snapshot.attributes();
    // the modelName may be sync/destination or a child depending if it was initiated from the list or details view
    // since the id for sync/destinations is type/name it will actually generate the correct url but the slash will be encoded
    // if we normalize to use the child model name for url generation instead things will be consistent
    const normalizedModelName =
      modelName === 'sync/destination' ? `${pluralize(modelName)}/${type}` : modelName;
    return `${super.urlForDeleteRecord(name, normalizedModelName, snapshot)}`;
  }

  query(store, { modelName }) {
    return this.ajax(this.buildURL(modelName), 'GET', { data: { list: true } });
  }

  // return normalized query response
  // useful for fetching data directly without loading models into store
  async normalizedQuery() {
    const queryResponse = await this.query(this.store, { modelName: 'sync/destination' });
    const serializer = this.store.serializerFor('sync/destination');
    return serializer.extractLazyPaginatedData(queryResponse);
  }
}
