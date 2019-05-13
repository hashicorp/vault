import ClusterAdapter from './cluster';

export default ClusterAdapter.extend({
  queryRecord() {
    return this._super(...arguments).then(resp => {
      resp.data.id = resp.data.license_id;
      return resp.data;
    });
  },

  createRecord(store, type, snapshot) {
    let id = snapshot.attr('licenseId');
    return this._super(...arguments).then(() => {
      return {
        id,
      };
    });
  },

  updateRecord(store, type, snapshot) {
    let id = snapshot.attr('licenseId');
    return this._super(...arguments).then(() => {
      return {
        id,
      };
    });
  },

  pathForType() {
    return 'license';
  },

  urlForUpdateRecord() {
    return this.buildURL() + '/license';
  },
});
