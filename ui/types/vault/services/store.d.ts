import Store, { RecordArray } from '@ember-data/store';

export default class StoreService extends Store {
  lazyPaginatedQuery(
    modelName: string,
    query: Object,
    options?: { adapterOptions: Object }
  ): Promise<RecordArray>;

  clearDataset(modelName: string);
}
