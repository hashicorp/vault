/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store, { RecordArray } from '@ember-data/store';

export default class StoreService extends Store {
  lazyPaginatedQuery(
    modelName: string,
    query: object,
    options?: { adapterOptions: object }
  ): Promise<RecordArray>;

  clearDataset(modelName: string);
  findRecord(modelName: string, path: string);
  peekRecord(modelName: string, path: string);
  query(modelName: string, query: object);
}
