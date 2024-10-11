/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { schedule } from '@ember/runloop';
import { resolve, Promise } from 'rsvp';
import { dasherize } from '@ember/string';
import { assert } from '@ember/debug';
import { set, get } from '@ember/object';
import clamp from 'vault/utils/clamp';
import config from 'vault/config/environment';
import sortObjects from 'vault/utils/sort-objects';

const { DEFAULT_PAGE_SIZE } = config.APP;

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

export default class Pagination extends Service {
  @service store;

  lazyCaches = new Map();

  setLazyCacheForModel(modelName, key, value) {
    const cacheKey = keyForCache(key);
    const cache = this.lazyCacheForModel(modelName) || new Map();
    cache.set(cacheKey, value);
    const modelKey = dasherize(modelName);
    this.lazyCaches.set(modelKey, cache);
  }

  getLazyCacheForModel(modelName, key) {
    const cacheKey = keyForCache(key);
    const modelCache = this.lazyCacheForModel(modelName);
    if (modelCache) {
      return modelCache.get(cacheKey);
    }
  }

  lazyCacheForModel(modelName) {
    return this.lazyCaches.get(dasherize(modelName));
  }

  // This is the public interface for the store extension - to be used just
  // like `Store.query`. Special handling of the response is controlled by
  // `query.pageFilter`, `query.page`, and `query.size`.

  // Required attributes of the `query` argument are:
  //   responsePath: a string indicating the location on the response where
  //     the array of items will be found
  //   page: the page number to return
  //   size: the size of the page
  //   pageFilter: a string that will be used to do a fuzzy match against the results,
  //     OR a function to be executed that will receive the dataset as the lone arg.
  //     Filter is done pre-pagination.
  lazyPaginatedQuery(modelType, query, adapterOptions) {
    const skipCache = query.skipCache;
    // We don't want skipCache to be part of the actual query key, so remove it
    delete query.skipCache;
    const adapter = this.store.adapterFor(modelType);
    const modelName = dasherize(modelType);
    const dataCache = skipCache ? this.clearDataset(modelName) : this.getDataset(modelName, query);
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
      .query(this.store, { modelName }, query, null, adapterOptions)
      .then((response) => {
        const serializer = this.store.serializerFor(modelName);
        const datasetHelper = serializer.extractLazyPaginatedData;
        const dataset = datasetHelper
          ? datasetHelper.call(serializer, response)
          : get(response, responsePath);
        set(response, responsePath, null);
        this.storeDataset(modelName, query, response, dataset);
        return this.fetchPage(modelName, query);
      })
      .catch(function (e) {
        throw e;
      });
  }

  filterData(filter, dataset) {
    let newData = dataset || [];
    if (filter) {
      if (filter instanceof Function) {
        newData = filter(dataset);
      } else {
        newData = dataset.filter((item) => {
          const id = item.id || item.name || item;
          return id.toLowerCase().includes(filter.toLowerCase());
        });
      }
    }
    return newData;
  }

  // reconstructs the original form of the response from the server
  // with an additional `meta` block
  //
  // the meta block includes:
  // currentPage, lastPage, nextPage, prevPage, total, filteredTotal
  constructResponse(modelName, query) {
    const { pageFilter, responsePath, size, page } = query;
    const { response, dataset } = this.getDataset(modelName, query);
    const resp = { ...response };
    const data = this.filterData(pageFilter, dataset);

    const lastPage = Math.ceil(data.length / size);
    const currentPage = clamp(page, 1, lastPage);
    const end = currentPage * size;
    const start = end - size;
    const slicedDataSet = data.slice(start, end);

    set(resp, responsePath || '', slicedDataSet);
    resp.meta = {
      currentPage,
      lastPage,
      nextPage: clamp(currentPage + 1, 1, lastPage),
      prevPage: clamp(currentPage - 1, 1, lastPage),
      total: dataset.length || 0,
      filteredTotal: data.length || 0,
      pageSize: size,
    };

    return resp;
  }

  // pushes records into the store and returns the result
  fetchPage(modelName, query) {
    const response = this.constructResponse(modelName, query);
    this.store.unloadAll(modelName);
    return new Promise((resolve) => {
      // push subset of records into the store
      schedule('destroy', () => {
        this.store.push(
          this.store
            .serializerFor(modelName)
            .normalizeResponse(this.store, this.store.modelFor(modelName), response, null, 'query')
        );
        const model = this.store.peekAll(modelName).slice();
        model.set('meta', response.meta);
        resolve(model);
      });
    });
  }

  // get cached data
  getDataset(modelName, query) {
    return this.getLazyCacheForModel(modelName, query);
  }

  // store data cache as { response, dataset}
  // also populated `lazyCaches` attribute
  storeDataset(modelName, query, response, array) {
    const dataset = query.sortBy ? sortObjects(array, query.sortBy) : array;
    const value = {
      response,
      dataset,
    };
    this.setLazyCacheForModel(modelName, query, value);
  }

  clearDataset(modelName) {
    if (!this.lazyCaches.size) return;
    if (modelName && this.lazyCaches.has(modelName)) {
      this.lazyCaches.delete(modelName);
      return;
    }
    this.lazyCaches.clear();
  }
}
