/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiSignIntermediateAdapter extends ApplicationAdapter {
  namespace = 'v1';

  createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const { backend, issuerRef } = snapshot.record;
    const url = `${this.buildURL()}/${encodePath(backend)}/issuer/${encodePath(issuerRef)}/sign-intermediate`;
    const data = serializer.serialize(snapshot, type);
    return this.ajax(url, 'POST', { data }).then((result) => ({
      // sign-intermediate can happen multiple times per issuer,
      // so the ID needs to be unique from the issuer ID
      id: result.request_id,
      ...result,
    }));
  }
}
