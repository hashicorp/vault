/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';
import { assert } from '@ember/debug';
import { all } from 'rsvp';

export default class SyncAssociationAdapter extends ApplicationAdapter {
  namespace = 'v1/sys/sync';

  buildURL(modelName, id, snapshot, requestType, query = {}) {
    const { destinationType, destinationName } = snapshot ? snapshot.attributes() : query;
    if (!destinationType || !destinationName) {
      return `${super.buildURL()}/associations`;
    }
    const { action } = snapshot?.adapterOptions || {};
    const uri = action ? `/${action}` : '';
    return `${super.buildURL()}/destinations/${destinationType}/${destinationName}/associations${uri}`;
  }

  query(store, { modelName }, query) {
    // endpoint doesn't accept the typical list query param and we don't want to pass options from lazyPaginatedQuery
    const url = this.buildURL(modelName, null, null, 'query', query);
    return this.ajax(url, 'GET');
  }

  // typically associations are queried for a specific destination which is what the standard query method does
  // in specific cases we can query all associations to access total_associations and total_secrets values
  queryAll() {
    const url = `${this.buildURL('sync/association')}`;
    return this.ajax(url, 'GET', { data: { list: true } }).then((response) => {
      const { total_associations, total_secrets } = response.data;
      return { total_associations, total_secrets };
    });
  }

  // fetch associations for many destinations
  // returns aggregated association information for each destination
  // information includes total associations, total unsynced and most recent updated datetime
  async fetchByDestinations(destinations) {
    const promises = destinations.map(({ name: destinationName, type: destinationType }) => {
      return this.query(this.store, { modelName: 'sync/association' }, { destinationName, destinationType });
    });
    const queryResponses = await all(promises);
    const serializer = this.store.serializerFor('sync/association');
    return queryResponses.map((response) => serializer.normalizeFetchByDestinations(response));
  }

  // array of association data for each destination a secret is synced to
  fetchSyncStatus({ mount, secretName }) {
    const url = `${this.buildURL()}/destinations`;
    return this.ajax(url, 'GET', { data: { mount, secret_name: secretName } }).then((resp) => {
      const { associated_destinations } = resp.data;
      const syncData = [];
      for (const key in associated_destinations) {
        const data = associated_destinations[key];
        // renaming keys to match query() response
        syncData.push({
          destinationType: data.type,
          destinationName: data.name,
          syncStatus: data.sync_status,
          updatedAt: data.updated_at,
        });
      }
      return syncData;
    });
  }

  // snapshot is needed for mount and secret_name values which are used to parse response since all associations are returned
  _setOrRemove(store, { modelName }, snapshot) {
    assert(
      "action type of set or remove required when saving association => association.save({ adapterOptions: { action: 'set' }})",
      ['set', 'remove'].includes(snapshot?.adapterOptions?.action)
    );
    const url = this.buildURL(modelName, null, snapshot);
    const data = snapshot.serialize();
    const serializer = store.serializerFor('sync/association');

    return this.ajax(url, 'POST', { data }).then((resp) => {
      const association = Object.values(resp.data.associated_secrets).find((association) => {
        return association.mount === data.mount && association.secret_name === data.secret_name;
      });

      // generate an id if an association is found
      // (an association may not be found if the secret is being unsynced)
      const id = association ? serializer.generateId(association) : undefined;

      return {
        ...association,
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
