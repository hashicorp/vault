import ApplicationAdapter from './application';

export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    return 'identity/mfa/login-enforcement';
  }

  _saveRecord(store, { modelName }, snapshot) {
    const data = store.serializerFor(modelName).serialize(snapshot);
    return this.ajax(this.urlForUpdateRecord(snapshot.attr('name'), modelName, snapshot), 'POST', {
      data,
    }).then(() => data);
  }
  // create does not return response similar to PUT request
  createRecord() {
    return this._saveRecord(...arguments);
  }
  // update record via POST method
  updateRecord() {
    return this._saveRecord(...arguments);
  }

  query(store, type, query) {
    const url = this.urlForQuery(query, type.modelName);
    return this.ajax(url, 'GET', { data: { list: true } });
  }
}
