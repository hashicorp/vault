/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';
import { schedule } from '@ember/runloop';
import { resolve, Promise } from 'rsvp';
import { dasherize } from '@ember/string';
import { assert } from '@ember/debug';
import { set, get } from '@ember/object';
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

export default class StoreService extends Store {
  lazyCaches = new Map();

  setLazyCacheForModel(modelName, key, value) {
    const cacheKey = keyForCache(key);
    const cache = this.lazyCacheForModel(modelName) || new Map();
    cache.set(cacheKey, value);
    const modelKey = normalizeModelName(modelName);
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
    return this.lazyCaches.get(normalizeModelName(modelName));
  }

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
      .then((response) => {
        const serializer = this.serializerFor(modelName);
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
      newData = dataset.filter(function (item) {
        const id = item.id || item;
        return id.toLowerCase().includes(filter.toLowerCase());
      });
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
    this.unloadAll(modelName);
    return new Promise((resolve) => {
      // after the above unloadRecords are finished, push into store
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
        const model = this.peekAll(modelName).toArray();
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
    const dataSet = {
      response,
      dataset: array,
    };
    this.setLazyCacheForModel(modelName, query, dataSet);
  }

  clearDataset(modelName) {
    if (!this.lazyCaches.size) return;
    if (modelName && this.lazyCaches.has(modelName)) {
      this.lazyCaches.delete(modelName);
      return;
    }
    this.lazyCaches.clear();
  }

  clearAllDatasets() {
    this.clearDataset();
  }
  /**
   * this is designed to be a temporary workaround to an issue in the test environment after upgrading to Ember 4.12
   * when performing an unloadAll or unloadRecord for auth-method or secret-engine models within the app code an error breaks the tests
   * after the test run is finished during teardown an unloadAll happens and the error "Expected a stable identifier" is thrown
   * it seems that when the unload happens in the app, for some reason the mount-config relationship models are not unloaded
   * then when the unloadAll happens a second time during test teardown there seems to be an issue since those records should already have been unloaded
   * when logging in the teardownRecord hook, it appears that other embedded inverse: null relationships such as replication-attributes are torn down when the parent model is unloaded
   * the following fixes the issue by explicitly unloading the mount-config models associated to the parent
   * this should be looked into further to find the root cause, at which time these overrides may be removed
   */
  unloadAll(modelName) {
    const hasMountConfig = ['auth-method', 'secret-engine'];
    if (hasMountConfig.includes(modelName)) {
      this.peekAll(modelName).forEach((record) => this.unloadRecord(record));
    } else {
      super.unloadAll(modelName);
    }
  }
  unloadRecord(record) {
    const hasMountConfig = ['auth-method', 'secret-engine'];
    if (record && hasMountConfig.includes(record.constructor.modelName) && record.config) {
      super.unloadRecord(record.config);
    }
    super.unloadRecord(record);
  }
}
