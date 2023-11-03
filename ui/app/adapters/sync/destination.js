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

  urlForDeleteRecord(id, modelName, snapshot) {
    const { name, type } = snapshot.attributes();
    // the modelName may be sync/destination or a child depending if it was initiated from the list or details view
    // since the id for sync/destinations is type/name it will actually generate the correct url but the slash will be encoded
    // if we normalize to use the child model name for url generation instead things will be consistent
    const normalizedModelName =
      modelName === 'sync/destination' ? `${pluralize(modelName)}/${type}` : modelName;
    return `${super.urlForDeleteRecord(name, normalizedModelName, snapshot)}`;
  }

  query(store, { modelName }) {
    return this.ajax(this.buildURL(modelName), 'GET', { data: { list: true } });
  }
}
