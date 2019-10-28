import BaseAdapter from './base';

export default BaseAdapter.extend({
  urlForFindRecord(id, modelName, snapshot) {
    let name = this.pathForType(modelName);
    return this.buildURL(id, name, snapshot);
  },
});
