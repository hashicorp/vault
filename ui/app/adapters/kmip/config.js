/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import BaseAdapter from './base';

export default BaseAdapter.extend({
  _url(id, modelName, snapshot) {
    const name = this.pathForType(modelName);
    // id here will be the mount path,
    // modelName will be config so we want to transpose the first two call args
    return this.buildURL(id, name, snapshot);
  },
  urlForFindRecord() {
    return this._url(...arguments);
  },
  urlForCreateRecord(modelName, snapshot) {
    const id = snapshot.record.mutableId;
    return this._url(id, modelName, snapshot);
  },
  urlForUpdateRecord() {
    return this._url(...arguments);
  },

  createRecord(store, type, snapshot) {
    return this._super(...arguments).then(() => {
      // saving returns a 204, return object with id to please ember-data...
      const id = snapshot.record.mutableId;
      return { id };
    });
  },
});
