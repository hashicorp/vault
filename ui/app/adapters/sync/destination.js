/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';
import { pluralize } from 'ember-inflector';

export default class SyncDestinationAdapter extends ApplicationAdapter {
  namespace = 'v1/sys';

  pathForType(modelName) {
    return modelName === 'sync/destination' ? pluralize(modelName) : modelName;
  }

  urlForCreateRecord(modelName, snapshot) {
    const { name } = snapshot.attributes();
    return `${super.urlForCreateRecord(modelName, snapshot)}/${name}`;
  }

  query(store, { modelName }) {
    return this.ajax(this.buildURL(modelName), 'GET', { data: { list: true } });
  }
}
