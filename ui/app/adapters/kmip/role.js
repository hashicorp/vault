import BaseAdapter from './base';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    let name = snapshot.id || snapshot.attr('name');
    let url = this._url(
      type.modelName,
      {
        backend: snapshot.record.backend,
        scope: snapshot.record.scope,
      },
      name
    );
    return this.ajax(url, 'POST', { data: snapshot.serialize() }).then(() => {
      return {
        id: name,
        name,
      };
    });
  },

  updateRecord() {
    return this.createRecord(...arguments);
  },
});
