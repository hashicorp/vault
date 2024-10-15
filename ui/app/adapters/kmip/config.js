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
    const id = snapshot.id || snapshot.record.mutableId;
    return this._url(id, modelName, snapshot);
  },
  urlForUpdateRecord() {
    return this._url(...arguments);
  },

  createRecord(store, type, snapshot) {
    return this._super(...arguments).then((response) => {
      // if the server does not return an id and one has not been set on the model we need to set it manually from the mutableId value
      if (!response?.id && !snapshot.record.id) {
        return { id: snapshot.record.mutableId };
      }
    });
  },
});
