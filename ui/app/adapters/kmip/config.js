import BaseAdapter from './base';

export default BaseAdapter.extend({
  _url(id, modelName, snapshot) {
    let name = this.pathForType(modelName);
    // id here will be the mount path,
    // modelName will be config so we want to transpose the first two call args
    return this.buildURL(id, name, snapshot);
  },
  urlForFindRecord() {
    return this._url(...arguments);
  },
  urlForCreateRecord(modelName, snapshot) {
    return this._url(snapshot.id, modelName, snapshot);
  },
  urlForUpdateRecord() {
    return this._url(...arguments);
  },
});
