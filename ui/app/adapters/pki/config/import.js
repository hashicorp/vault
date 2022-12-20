import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../../application';

export default class PkiConfigImportAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(modelName, snapshot) {
    // TODO: may need to check permissions to decide which path to use
    const { pemBundle, backend } = snapshot.record;
    if (!backend) {
      throw new Error('URL for create record is missing required attributes');
    }
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;
    if (pemBundle) {
      return `${baseUrl}/config/ca`;
    }
    return `${baseUrl}/intermediate/set-signed`;
  }
}
