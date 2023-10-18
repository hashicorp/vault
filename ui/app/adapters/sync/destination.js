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

  query() {
    const url = `${this._baseUrl()}/sync/destinations`;
    return this.ajax(url, 'GET', { data: { list: true } }).then((resp) => {
      return this._transformQueryResponse(resp);
    });
  }

  /*
the API response keys map to destination types (see below)
for lazyPagination we need each key to correspond to a unique destination

* original API LIST response 
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
      name: 'my-dest-1',
      type: 'aws-sm',
    },
    'gh/my-dest-1': {
      name: 'my-dest-1',
      type: 'gh',
    },
  },
  keys: ['aws-sm/destination-aws', 'gh/destination-gh'],
},
*/
  _transformQueryResponse(resp) {
    const transformedKeyInfo = {};
    // loop through each destination type (keys in key_info)
    for (const type in resp.data.key_info) {
      // iterate through each type's destination names
      resp.data.key_info[type].forEach((name) => {
        // generate id
        const id = `${type}/${name}`;
        // create object with destination's attributes for key_info
        transformedKeyInfo[id] = { name, type };
      });
    }
    // update the response here (instead of serializer) to be compatible with lazyPaginatedQuery
    resp.data.keys = Object.keys(transformedKeyInfo);
    resp.data.key_info = transformedKeyInfo;
    return resp;
  }
}
