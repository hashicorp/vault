/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { RecordArray } from '@ember-data/store';

export default class PaginationService extends Service {
  lazyPaginatedQuery(
    modelName: string,
    query: object,
    options?: { adapterOptions: object }
  ): Promise<RecordArray>;
  clearDataset(modelName: string);
}
