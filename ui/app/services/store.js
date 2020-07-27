import { schedule } from '@ember/runloop';
import { copy } from 'ember-copy';
import { resolve, Promise } from 'rsvp';
import { dasherize } from '@ember/string';
import { assert } from '@ember/debug';
import { set, get, computed } from '@ember/object';
import DS from 'ember-data';
import clamp from 'vault/utils/clamp';
import config from 'vault/config/environment';

const { DEFAULT_PAGE_SIZE } = config.APP;

export function normalizeModelName(modelName) {
  return dasherize(modelName);
}

export function keyForCache(query) {
  /*eslint no-unused-vars: ["error", { "ignoreRestSiblings": true }]*/
  // we want to ignore size, page, responsePath, and pageFilter in the cacheKey
  const { size, page, responsePath, pageFilter, ...queryForCache } = query;
  const cacheKeyObject = Object.keys(queryForCache)
    .sort()
    .reduce((result, key) => {
      result[key] = queryForCache[key];
      return result;
    }, {});
  return JSON.stringify(cacheKeyObject);
}

export default DS.Store.extend({
  // this is a map of map that stores the caches
  lazyCaches: computed(function() {
    return new Map();
  }),

  setLazyCacheForModel(modelName, key, value) {
    const cacheKey = keyForCache(key);
    const cache = this.lazyCacheForModel(modelName) || new Map();
    cache.set(cacheKey, value);
    const lazyCaches = this.get('lazyCaches');
    const modelKey = normalizeModelName(modelName);
    lazyCaches.set(modelKey, cache);
  },

  getLazyCacheForModel(modelName, key) {
    const cacheKey = keyForCache(key);
    const modelCache = this.lazyCacheForModel(modelName);
    if (modelCache) {
      return modelCache.get(cacheKey);
    }
  },

  lazyCacheForModel(modelName) {
    return this.get('lazyCaches').get(normalizeModelName(modelName));
  },

  // This is the public interface for the store extension - to be used just
  // like `Store.query`. Special handling of the response is controlled by
  // `query.pageFilter`, `query.page`, and `query.size`.

  // Required attributes of the `query` argument are:
  //   responsePath: a string indicating the location on the response where
  //     the array of items will be found
  //   page: the page number to return
  //   size: the size of the page
  //   pageFilter: a string that will be used to do a fuzzy match against the
  //     results, this is done pre-pagination
  lazyPaginatedQuery(modelType, query /*, options*/) {
    const adapter = this.adapterFor(modelType);
    const modelName = normalizeModelName(modelType);
    const dataCache = this.getDataset(modelName, query);
    const responsePath = query.responsePath;
    assert('responsePath is required', responsePath);
    assert('page is required', typeof query.page === 'number');
    if (!query.size) {
      query.size = DEFAULT_PAGE_SIZE;
    }

    if (dataCache) {
      return resolve(this.fetchPage(modelName, query));
    }
    return adapter
      .query(this, { modelName }, query)
      .then(response => {
        const serializer = this.serializerFor(modelName);
        const datasetHelper = serializer.extractLazyPaginatedData;
        const dataset = datasetHelper
          ? datasetHelper.call(serializer, response)
          : get(response, responsePath);
        set(response, responsePath, null);
        this.storeDataset(modelName, query, response, dataset);
        return this.fetchPage(modelName, query);
      })
      .catch(function(e) {
        throw e;
      });
  },

  filterData(filter, dataset) {
    let newData = dataset || [];
    if (filter) {
      newData = dataset.filter(function(item) {
        const id = item.id || item;
        return id.toLowerCase().includes(filter.toLowerCase());
      });
    }
    return newData;
  },

  // reconstructs the original form of the response from the server
  // with an additional `meta` block
  //
  // the meta block includes:
  // currentPage, lastPage, nextPage, prevPage, total, filteredTotal
  constructResponse(modelName, query) {
    const { pageFilter, responsePath, size, page } = query;
    let { response, dataset } = this.getDataset(modelName, query);
    response = copy(response, true);
    const data = this.filterData(pageFilter, dataset);

    const lastPage = Math.ceil(data.length / size);
    const currentPage = clamp(page, 1, lastPage);
    const end = currentPage * size;
    const start = end - size;
    const slicedDataSet = data.slice(start, end);

    set(response, responsePath || '', slicedDataSet);

    response.meta = {
      currentPage,
      lastPage,
      nextPage: clamp(currentPage + 1, 1, lastPage),
      prevPage: clamp(currentPage - 1, 1, lastPage),
      total: get(dataset, 'length') || 0,
      filteredTotal: get(data, 'length') || 0,
    };

    return response;
  },

  // pushes records into the store and returns the result
  fetchPage(modelName, query) {
    const response = this.constructResponse(modelName, query);
    this.peekAll(modelName).forEach(record => {
      record.unloadRecord();
    });
    return new Promise(resolve => {
      schedule('destroy', () => {
        this.push(
          this.serializerFor(modelName).normalizeResponse(
            this,
            this.modelFor(modelName),
            response,
            null,
            'query'
          )
        );
        let model = this.peekAll(modelName).toArray();
        model.set('meta', response.meta);
        resolve(model);
      });
    });
  },

  // get cached data
  getDataset(modelName, query) {
    return this.getLazyCacheForModel(modelName, query);
  },

  // store data cache as { response, dataset}
  // also populated `lazyCaches` attribute
  storeDataset(modelName, query, response, array) {
    const dataSet = {
      response,
      dataset: array,
    };
    this.setLazyCacheForModel(modelName, query, dataSet);
  },

  clearDataset(modelName) {
    let cacheList = this.get('lazyCaches');
    if (!cacheList.size) return;
    if (modelName && cacheList.has(modelName)) {
      cacheList.delete(modelName);
      return;
    }
    cacheList.clear();
    this.set('lazyCaches', cacheList);
  },

  clearAllDatasets() {
    this.clearDataset();
  },
});
