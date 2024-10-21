/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';

export default class StoreService extends Store {
  adapterFor(modelName: string);
  createRecord(modelName: string, object);
  findRecord(modelName: string, path: string);
  peekRecord(modelName: string, path: string);
  query(modelName: string, query: object);
}
