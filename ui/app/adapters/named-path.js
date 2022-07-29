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
    return this._saveRecord(...arguments);
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
    return this.ajax(url, 'GET', { data: { list: true } });
  }
}
