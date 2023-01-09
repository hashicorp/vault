import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiConfigAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(modelName, snapshot) {
    const { backend, formType, type } = snapshot.record;
    if (!backend || !formType) {
      throw new Error('URL for create record is missing required attributes');
    }
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}/issuers`;
    switch (formType) {
      case 'generate-root':
        return `${baseUrl}/generate/root/${type}`;
      case 'generate-csr':
        return `${baseUrl}/generate/intermediate/${type}`;
      default:
        return `${baseUrl}/import/bundle`;
    }
  }
}
