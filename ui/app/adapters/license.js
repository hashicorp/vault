import ClusterAdapter from './cluster';

export default ClusterAdapter.extend({
  queryRecord() {
    return this.ajax(this.buildURL() + '/license', 'GET').then(resp => {
      resp.data.id = resp.data.license_id;
      return resp.data;
    });
  },

  urlForQueryRecord() {
    return this.buildURL() + '/license';
  },
});
