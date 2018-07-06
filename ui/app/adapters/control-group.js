import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'control-group';
  },

  findRecord(store, type, id) {
    let baseUrl = this.buildURL(type.modelName);
    return this.ajax(`${baseUrl}/request`, 'POST', {
      data: {
        accessor: id,
      },
    }).then(response => {
      response.id = id;
      return response;
    });
  },

  urlForUpdateRecord(id, modelName) {
    let base = this.buildURL(modelName);
    return `${base}/authorize`;
  },
});
