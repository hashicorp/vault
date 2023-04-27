/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { parseCertificate } from 'vault/utils/parse-pki-cert';
import ApplicationSerializer from '../application';

export default class PkiIssuerSerializer extends ApplicationSerializer {
  primaryKey = 'issuer_id';

  attrs = {
    caChain: { serialize: false },
    certificate: { serialize: false },
    issuerId: { serialize: false },
    keyId: { serialize: false },
    parsedCertificate: { serialize: false },
    commonName: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.certificate) {
      // Parse certificate back from the API and add to payload
      const parsedCert = parseCertificate(payload.data.certificate);
      const data = {
        issuer_ref: payload.issuer_id,
        ...payload.data,
        parsed_certificate: parsedCert,
        common_name: parsedCert.common_name,
      };
      return super.normalizeResponse(store, primaryModelClass, { ...payload, data }, id, requestType);
    }
    return super.normalizeResponse(...arguments);
  }

  // rehydrate each issuers model so all model attributes are accessible from the LIST response
  normalizeItems(payload) {
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((issuer_id) => ({
          issuer_id,
          ...payload.data.key_info[issuer_id],
        }));
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }
}
