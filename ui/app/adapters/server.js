import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  urlForFindAll() {
    return '/v1/sys/storage/raft/configuration';
  },
});
