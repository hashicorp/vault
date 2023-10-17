/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';

export default class SyncDestinationAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _baseUrl() {
    return `${this.buildURL()}/sys`;
  }

  // id is the destination name
  // modelName is sync/destinations/:type
  urlForFindRecord(id, modelName) {
    return `${this._baseUrl()}/${modelName}/${id}`;
  }

  // the LIST response keys map to destination types (see below)
  // we want each key to correspond to a unique destination
  // and key_info to contain the model attributes
  query() {
    const url = `${this._baseUrl()}/sync/destinations`;
    return this.ajax(url, 'GET', { data: { list: true } }).then(async (resp) => {
      const transformedKeyInfo = {};
      // loop through each destination type (keys in key_info)
      for (const type in resp.data.key_info) {
        // iterate through type's destinations and generate id
        resp.data.key_info[type].forEach((name) => {
          const id = `${type}/${name}`;
          // add object to updated key_info data
          transformedKeyInfo[id] = { id, name, type };
        });
      }
      // update the data here (instead of serializer) to be compatible with lazyPaginatedQuery
      resp.data.keys = Object.keys(transformedKeyInfo);
      resp.data.key_info = transformedKeyInfo;
      return resp;
    });
  }
}

/*
* original LIST response 
 data: {
  key_info: {
    'aws-sm': ['my-dest-1'],
    'gh': ['my-dest-1'],
  },
  keys: ['aws-sm', 'gh'],
}

* transformed LIST response 
data: {
  key_info: {
    'aws-sm/my-dest-1': {
      id: 'aws-sm/my-dest-1',
      name: 'my-dest-1',
      type: 'aws-sm',
    },
    'gh/my-dest-1': {
      id: 'gh/my-dest-1',
      name: 'my-dest-1',
      type: 'gh',
    },
  },
  keys: ['aws-sm/destination-aws', 'gh/destination-gh'],
},
*/
