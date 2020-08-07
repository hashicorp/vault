import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  urlForCreateRecord() {
    return '/v1/sys/storage/raft/join';
  },
});
