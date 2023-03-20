/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { parseCertificate } from 'vault/utils/parse-pki-cert';
import ApplicationSerializer from '../../application';

export default class PkiCertificateBaseSerializer extends ApplicationSerializer {
  primaryKey = 'serial_number';

  attrs = {
    role: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.certificate) {
      // Parse certificate back from the API and add to payload
      const parsedCert = parseCertificate(payload.data.certificate);
      const json = super.normalizeResponse(
        store,
        primaryModelClass,
        { ...payload, ...parsedCert },
        id,
        requestType
      );
      return json;
    }
    return super.normalizeResponse(...arguments);
  }
}
