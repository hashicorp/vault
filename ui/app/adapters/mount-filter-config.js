import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  url(id) {
    return `${this.buildURL()}/replication/performance/primary/mount-filter/${id}`;
  },

  findRecord(store, type, id) {
    return this.ajax(this.url(id), 'GET').then(resp => {
      resp.id = id;
      return resp;
    });
  },

  createRecord(store, type, snapshot) {
    return this.ajax(this.url(snapshot.id), 'PUT', {
      data: this.serialize(snapshot),
    });
  },

  updateRecord() {
    return this.createRecord(...arguments);
  },

  deleteRecord(store, type, snapshot) {
    return this.ajax(this.url(snapshot.id), 'DELETE');
  },
});
