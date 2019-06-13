import BaseAdapter from './base';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    let name = snapshot.attr('name');
    let url = this._url(type.modelName, {
      backend: snapshot.record.backend,
      scope: snapshot.record.scope,
      role: snapshot.record.role,
    });
    return this.ajax(url, 'POST', { data: snapshot.serialize() }).then(() => {
      // TODO change this to serial?
      return {
        id: name,
        name,
      };
    });
  },
});
