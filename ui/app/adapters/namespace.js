import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'namespaces';
  },
  urlForFindAll() {
    return `/${this.urlPrefix()}/namespaces?list=true`;
  },
  urlForCreateRecord(modelName, snapshot) {
    let id = snapshot.attr('path');
    return this.buildURL(modelName, id);
  },

  createRecord(store, type, snapshot) {
    let id = snapshot.attr('path');
    return this._super(...arguments).then(() => {
      return { id };
    });
  },
});
