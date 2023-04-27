/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { assert } from '@ember/debug';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiActionAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForCreateRecord(modelName, snapshot) {
    const { type } = snapshot.record;
    const { actionType, useIssuer, issuerRef, mount } = snapshot.adapterOptions;
    // if the backend mount is passed, we want that to override the URL's mount path
    const backend = mount || snapshot.record.backend;
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
      case 'sign-intermediate':
        return `${baseUrl}/issuer/${encodePath(issuerRef)}/sign-intermediate`;
      case 'rotate-root':
        return `${baseUrl}/root/rotate/${type}`;
      default:
        assert('actionType must be one of import, generate-root, generate-csr or sign-intermediate');
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
