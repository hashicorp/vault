import { assert } from '@ember/debug';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiActionAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(modelName, snapshot) {
    const { backend, type } = snapshot.record;
    const { actionType, useIssuer } = snapshot.adapterOptions;
    if (!backend || !actionType) {
      throw new Error('URL for create record is missing required attributes');
    }
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;
    switch (actionType) {
      case 'import':
        return useIssuer ? `${baseUrl}/issuers/import/bundle` : `${baseUrl}/config/ca`;
      case 'generate-root':
        return useIssuer ? `${baseUrl}/issuers/generate/root/${type}` : `${baseUrl}/root/generate/${type}`;
      case 'generate-csr':
        return useIssuer
          ? `${baseUrl}/issuers/generate/intermediate/${type}`
          : `${baseUrl}/intermediate/generate/${type}`;
      default:
        assert('actionType must be one of import, generate-root, or generate-csr');
    }
  }

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const url = this.urlForCreateRecord(type.modelName, snapshot);
    // Send actionType as serializer requestType so that we serialize data based on the endpoint
    const data = serializer.serialize(snapshot, snapshot.adapterOptions.actionType);
    return this.ajax(url, 'POST', { data }).then((result) => ({
      // pki/action endpoints don't correspond with a single specific entity,
      // so in ember-data we'll map it to the request ID
      id: result.request_id,
      ...result,
    }));
  }
}
