/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class PkiTidyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _baseUrl(backend) {
    return `${this.buildURL()}/${encodePath(backend)}`;
  }

  // single tidy operations (manual) are always a new record
  createRecord(store, type, snapshot) {
    const { backend } = snapshot.record;
    const { tidyType } = snapshot.adapterOptions;
    if (tidyType === 'auto') {
      throw new Error('Auto tidy type models are never new, please use findRecord');
    }

    const url = `${this._baseUrl(backend)}/tidy`;
    return this.ajax(url, 'POST', { data: this.serialize(snapshot, tidyType) });
  }

  // saving auto-tidy config POST requests will always update
  updateRecord(store, type, snapshot) {
    const backend = snapshot.record.id;
    const { tidyType } = snapshot.adapterOptions;
    if (tidyType === 'manual') {
      throw new Error('Manual tidy type models are always new, please use createRecord');
    }

    const url = `${this._baseUrl(backend)}/config/auto-tidy`;
    return this.ajax(url, 'POST', { data: this.serialize(snapshot, tidyType) });
  }

  findRecord(store, type, backend) {
    // only auto-tidy will ever be read, no need to pass the type here
    return this.ajax(`${this._baseUrl(backend)}/config/auto-tidy`, 'GET').then((resp) => {
      return resp.data;
    });
  }

  cancelTidy(backend) {
    const url = `${this._baseUrl(backend)}`;
    return this.ajax(`${url}/tidy-cancel`, 'POST');
  }
}
