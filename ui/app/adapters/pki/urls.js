import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiUrlsAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(modelName, snapshot) {
    const { backend } = snapshot.record;
    if (!backend) {
      throw new Error('Backend required on model for URL');
    }
    return `${this.buildURL()}/${encodePath(backend)}/config/urls`;
  }
}
