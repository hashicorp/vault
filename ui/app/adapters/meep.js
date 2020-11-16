import ClusterAdapter from './cluster';

export default ClusterAdapter.extend({
  // example adapter file from the recent Key Usage metrics config page.  See url and network request http://localhost:4200/ui/vault/metrics/config
  // queryRecord() {
  //   return this.ajax(this.urlForQuery(), 'GET').then(resp => {
  //     resp.id = resp.request_id;
  //     return resp;
  //   });
  // },
  // urlForUpdateRecord() {
  //   return this.buildURL() + '/internal/counters/config';
  // },
  // urlForQuery() {
  //   return this.buildURL() + '/internal/counters/config';
  // },
});
