/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import { assert } from '@ember/debug';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiTidyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _baseUrl(backend) {
    return `${this.buildURL()}/${encodePath(backend)}`;
  }

  urlForCreateRecord(snapshot) {
    const { backend } = snapshot.record;
    const { tidyType } = snapshot.adapterOptions;

    if (!backend) {
      throw new Error('Backend missing');
    }
    switch (tidyType) {
      case 'manual-tidy':
        return `${this._baseUrl(backend)}/tidy`;
      case 'auto-tidy':
        return `${this._baseUrl(backend)}/config/auto-tidy`;
      default:
        assert('type must be one of manual-tidy, auto-tidy');
    }
  }

  createRecord(store, type, snapshot) {
    const url = this.urlForCreateRecord(snapshot);
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) });
  }

  findRecord(store, type, backend) {
    // only auto-tidy will ever be read, no need to pass the type here
    return this.ajax(`${this._baseUrl(backend)}/config/auto-tidy`, 'GET').then((resp) => {
      return resp.data;
    });
  }

  queryRecord(store, type, query) {
    const { backend, tidyType } = query;
    // only auto-tidy will ever be read, no need to pass the type here
    return this.ajax(`${this._baseUrl(backend)}/config/auto-tidy`, 'GET').then((resp) => {
      // tidyType is the primary key and sets the id for the ember data model
      return { tidyType, ...resp.data };
    });
  }
}
