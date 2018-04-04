import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: '/v1/auth',

  pathForType(modelType) {
    // we want the last part of the path
    const type = modelType.split('/').pop();
    if (type === 'identity-whitelist' || type === 'roletag-blacklist') {
      return `tidy/${type}`;
    }
    return type;
  },

  buildURL(modelName, id, snapshot) {
    const backendId = id ? id : snapshot.belongsTo('backend').id;
    let url = `${this.get('namespace')}/${backendId}/config`;
    // aws has a lot more config endpoints
    if (modelName.includes('aws')) {
      url = `${url}/${this.pathForType(modelName)}`;
    }
    return url;
  },

  createRecord(store, type, snapshot) {
    const id = snapshot.belongsTo('backend').id;
    return this._super(...arguments).then(() => {
      return { id };
    });
  },

  updateRecord(store, type, snapshot) {
    const id = snapshot.belongsTo('backend').id;
    return this._super(...arguments).then(() => {
      return { id };
    });
  },
});
