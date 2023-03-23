/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { all } from 'rsvp';

export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';
  listPayload = { data: { list: true } };

  pathForType() {
    // backend name prepended in buildURL method
    return 'kms';
  }
  buildURL(modelName, id, snapshot, requestType, query) {
    let url = super.buildURL(...arguments);
    if (snapshot) {
      url = url.replace('kms', `${snapshot.attr('backend')}/kms`);
    } else if (query) {
      url = url.replace('kms', `${query.backend}/kms`);
    }
    return url;
  }
  buildKeysURL(query) {
    const url = this.buildURL('keymgmt/provider', null, null, 'query', query);
    return `${url}/${query.provider}/key`;
  }
  async createRecord(store, { modelName }, snapshot) {
    // create uses PUT instead of POST
    const data = store.serializerFor(modelName).serialize(snapshot);
    const url = this.buildURL(modelName, snapshot.attr('name'), snapshot, 'updateRecord');
    return this.ajax(url, 'PUT', { data }).then(() => data);
  }
  findRecord(store, type, name) {
    return super.findRecord(...arguments).then((resp) => {
      resp.data = { ...resp.data, name };
      return resp;
    });
  }
  async query(store, type, query) {
    const { backend } = query;
    const url = this.buildURL(type.modelName, null, null, 'query', query);
    return this.ajax(url, 'GET', this.listPayload).then(async (resp) => {
      // additional data is needed to fullfil the list view requirements
      // pull in full record for listed items
      const records = await all(
        resp.data.keys.map((name) => this.findRecord(store, type, name, this._mockSnapshot(query.backend)))
      );
      resp.data.keys = records.map((record) => record.data);
      resp.backend = backend;
      return resp;
    });
  }
  async queryRecord(store, type, query) {
    return this.findRecord(store, type, query.id, this._mockSnapshot(query.backend));
  }

  // when using find in query or queryRecord overrides snapshot is not available
  // ultimately buildURL requires the snapshot to pull the backend name for the dynamic segment
  // since we have the backend value from the query generate a mock snapshot
  _mockSnapshot(backend) {
    return {
      attr(prop) {
        return prop === 'backend' ? backend : null;
      },
    };
  }
}
