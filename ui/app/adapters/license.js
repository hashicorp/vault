import ClusterAdapter from './cluster';

export default ClusterAdapter.extend({
  queryRecord() {
    return this.ajax(this.buildURL() + '/license', 'GET').then(resp => {
      resp.data.id = resp.data.license_id;
      return resp.data;
    });
  },

  // urlForUpdateRecord(){
  //   return this.buildURL() + '/license';
  // },

  urlForQueryRecord() {
    return this.buildURL() + '/license';
  },
});
