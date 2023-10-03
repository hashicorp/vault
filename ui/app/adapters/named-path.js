/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * base adapter for resources that are saved to a path whose unique identifier is name
 * save requests are made to the same endpoint and the resource is either created if not found or updated
 * */
import ApplicationAdapter from './application';
import { assert } from '@ember/debug';
export default class NamedPathAdapter extends ApplicationAdapter {
  namespace = 'v1';
  saveMethod = 'POST'; // override when extending if PUT is used rather than POST

  _saveRecord(store, { modelName }, snapshot) {
    // since the response is empty return the serialized data rather than nothing
    const data = store.serializerFor(modelName).serialize(snapshot);
    return this.ajax(this.urlForUpdateRecord(snapshot.attr('name'), modelName, snapshot), this.saveMethod, {
      data,
    }).then(() => data);
  }

  // create does not return response similar to PUT request
  createRecord() {
    const [store, { modelName }, snapshot] = arguments;
    const name = snapshot.attr('name');
    // throw error if user attempts to create a record with same name, otherwise POST request silently overrides (updates) the existing model
    if (store.hasRecordForId(modelName, name)) {
      throw new Error(`A record already exists with the name: ${name}`);
    } else {
      return this._saveRecord(...arguments);
    }
  }

  // update uses same endpoint and method as create
  updateRecord() {
    return this._saveRecord(...arguments);
  }

  // if backend does not return name in response Ember Data will throw an error for pushing a record with no id
  // use the id (name) supplied to findRecord to set property on response data
  findRecord(store, type, name) {
    return super.findRecord(...arguments).then((resp) => {
      if (!resp.data.name) {
        resp.data.name = name;
      }
      return resp;
    });
  }

  // GET request with list=true as query param
  async query(store, type, query) {
    const url = this.urlForQuery(query, type.modelName);
    const { paramKey, filterFor, allowed_client_id } = query;
    // * 'paramKey' is a string of the param name (model attr) we're filtering for, e.g. 'client_id'
    // * 'filterFor' is an array of values to filter for (value type must match the attr type), e.g. array of ID strings
    // * 'allowed_client_id' is a valid query param to the /provider endpoint
    const queryParams = { list: true, ...(allowed_client_id && { allowed_client_id }) };
    const response = await this.ajax(url, 'GET', { data: queryParams });

    // filter LIST response only if key_info exists and query includes both 'paramKey' & 'filterFor'
    if (filterFor) assert('filterFor must be an array', Array.isArray(filterFor));
    if (response.data.key_info && filterFor && paramKey && !filterFor.includes('*')) {
      const data = this.filterListResponse(paramKey, filterFor, response.data.key_info);
      return { ...response, data };
    }
    return response;
  }

  filterListResponse(paramKey, matchValues, key_info) {
    const keyInfoAsArray = Object.entries(key_info);
    const filtered = keyInfoAsArray.filter((key) => {
      const value = key[1]; // value is an object of model attributes
      return matchValues.includes(value[paramKey]);
    });
    const filteredKeyInfo = Object.fromEntries(filtered);
    const filteredKeys = Object.keys(filteredKeyInfo);
    return { keys: filteredKeys, key_info: filteredKeyInfo };
  }
}
