/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import BaseAdapter from './base';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    const name = snapshot.attr('name');
    return this.ajax(this._url(type.modelName, { backend: snapshot.record.backend }, name), 'POST').then(
      () => {
        return {
          id: name,
          name,
        };
      }
    );
  },

  deleteRecord(store, type, snapshot) {
    let url = this._url(type.modelName, { backend: snapshot.record.backend }, snapshot.id);
    url = `${url}?force=true`;
    return this.ajax(url, 'DELETE');
  },
});
