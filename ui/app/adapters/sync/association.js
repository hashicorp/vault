/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';

export default class SyncDestinationAdapter extends ApplicationAdapter {
  namespace = 'v1/sys/sync/destinations';

  buildURL(modelName, id, snapshot, requestType, query) {
    const { destinationType, destinationName, syncStatus } = snapshot ? snapshot.attributes() : query;
    // use sync_status to determine whether saving a record should use set or remove endpoint
    // new records will not have a status so the only scenario where remove should be used if status is synced
    const action = requestType === 'query' ? '' : syncStatus === 'SYNCED' ? '/remove' : '/set';
    return `${super.buildURL()}/${destinationType}/${destinationName}/associations${action}`;
  }

  query(store, { modelName }, query) {
    // endpoint doesn't accept the typical list query param and we don't want to pass options from lazyPaginatedQuery
    const url = this.buildURL(modelName, null, null, 'query', query);
    return this.ajax(url, 'GET');
  }

  // snapshot is needed for mount and secret_name values which are used to parse response since all associations are returned
  _setOrRemove(store, { modelName }, snapshot) {
    const url = this.buildURL(modelName, null, snapshot);
    const data = snapshot.serialize();
    return this.ajax(url, 'POST', { data }).then((resp) => {
      const id = `${data.mount}/${data.secret_name}`;
      return {
        ...resp.data.associated_secrets[id],
        id,
        destinationName: resp.data.store_name,
        destinationType: resp.data.store_type,
      };
    });
  }

  createRecord() {
    return this._setOrRemove(...arguments);
  }

  updateRecord() {
    return this._setOrRemove(...arguments);
  }
}
