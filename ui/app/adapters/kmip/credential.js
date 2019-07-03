import BaseAdapter from './base';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    let url = this._url(type.modelName, {
      backend: snapshot.record.backend,
      scope: snapshot.record.scope,
      role: snapshot.record.role,
    });
    url = `${url}/generate`;
    return this.ajax(url, 'POST', { data: snapshot.serialize() }).then(model => {
      model.data.id = model.data.serial_number;
      return model;
    });
  },

  deleteRecord(store, type, snapshot) {
    let url = this._url(type.modelName, {
      backend: snapshot.record.backend,
      scope: snapshot.record.scope,
      role: snapshot.record.role,
    });
    url = `${url}/revoke`;
    return this.ajax(url, 'POST', {
      data: {
        serial_number: snapshot.id,
      },
    });
  },
});
