/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { parsedParameters } from 'vault/utils/parse-pki-cert-oids';
import ApplicationSerializer from '../application';

export default class PkiIssuerSerializer extends ApplicationSerializer {
  primaryKey = 'issuer_id';

  constructor() {
    super(...arguments);
    // remove following attrs from serialization
    const attrs = ['caChain', 'certificate', 'issuerId', 'keyId', ...parsedParameters];
    this.attrs = attrs.reduce((attrObj, attr) => {
      attrObj[attr] = { serialize: false };
      return attrObj;
    }, {});
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.certificate) {
      // Parse certificate back from the API and add to payload
      const parsedCert = parseCertificate(payload.data.certificate || payload.data.csr);
      const data = { issuer_ref: payload.issuer_id, ...payload.data, parsed_certificate: parsedCert };
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
