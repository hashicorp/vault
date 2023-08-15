/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../../application';

export default class PkiConfigBaseAdapter extends ApplicationAdapter {
  namespace = 'v1';

  findRecord(store, type, backend) {
    return this.ajax(this._url(backend), 'GET').then((resp) => {
      return resp.data;
    });
  }

  updateRecord(store, type, snapshot) {
    const data = snapshot.serialize();
    return this.ajax(this._url(snapshot.record.id), 'POST', { data });
  }
}
