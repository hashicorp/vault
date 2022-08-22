/**
 * base adapter for resources that are saved to a path whose unique identifier is name
 * save requests are made to the same endpoint and the resource is either created if not found or updated
 * */
import ApplicationAdapter from './application';

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
    let [store, { modelName }, snapshot] = arguments;
    let name = snapshot.attr('name');
    // throw error if user attempts to create a record with same name, otherwise POST request silently overrides (updates) the existing model
    if (store.hasRecordForId(modelName, name))
      throw new Error(`A record already exists with the name: ${name}`);
    else {
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
    const { paramKey, filterFor } = query;
    // * 'paramKey' is a string of the param name (model attr) we're filtering for, e.g. 'client_id'
    // * 'filterFor' is an array of values to filter for (value type must match the attr type), e.g. array of ID strings
    // example: the OidcProviderClientsRoute where we only want to list clients that are permitted to use the currently viewed provider
    const response = await this.ajax(url, 'GET', { data: { list: true } });

    // filter LIST response only if key_info exists and query includes both 'paramKey' & 'filterFor'
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
