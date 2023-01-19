import { assert } from '@ember/debug';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiConfigAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(modelName, snapshot) {
    const { backend, type } = snapshot.record;
    const { formType, useIssuer } = snapshot.adapterOptions;
    if (!backend || !formType) {
      throw new Error('URL for create record is missing required attributes');
    }
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;
    switch (formType) {
      case 'import':
        return useIssuer ? `${baseUrl}/issuers/import/bundle` : `${baseUrl}/config/ca`;
      case 'generate-root':
        return useIssuer ? `${baseUrl}/issuers/generate/root/${type}` : `${baseUrl}/root/generate/${type}`;
      case 'generate-csr':
        return useIssuer
          ? `${baseUrl}/issuers/generate/intermediate/${type}`
          : `${baseUrl}/intermediate/generate/${type}`;
      default:
        assert('formType must be one of import, generate-root, or generate-csr');
    }
  }

  createRecord(store, type, snapshot) {
    const { mount, actionType, issuerName } = snapshot?.adapterOptions;
    let url = `${this.buildURL()}/${encodePath(mount)}`;
    switch (actionType) {
      case 'generate-csr':
        url = `${url}/intermediate/generate/existing`;
        break;
      case 'sign-intermediate':
        url = `${url}/issuer/${issuerName}/sign-intermediate`;
        break;
      default:
        return `uhoh`;
    }
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      return resp;
    });
  }
}
