import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  urlForFindAll() {
    return '/v1/sys/storage/raft/configuration';
  },
  urlForDeleteRecord() {
    return '/v1/sys/storage/raft/remove-peer';
  },
  deleteRecord(store, type, snapshot) {
    let server_id = snapshot.attr('nodeId');
    let url = '/v1/sys/storage/raft/remove-peer';
    return this.ajax(url, 'POST', { data: { server_id } });
  },
});
