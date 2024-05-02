/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { parseCertificate } from 'vault/utils/parse-pki-cert';
import ApplicationSerializer from '../../application';
import { encodeString } from 'core/utils/b64';

export default class PkiCertificateBaseSerializer extends ApplicationSerializer {
  primaryKey = 'serial_number';

  attrs = {
    role: { serialize: false },
  };

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.certificate) {
      // Parse certificate back from the API and add to payload
      const parsedCert = parseCertificate(payload.data.certificate);
      return super.normalizeResponse(
        store,
        primaryModelClass,
        { ...payload, parsed_certificate: parsedCert, common_name: parsedCert.common_name },
        id,
        requestType
      );
    }
    return super.normalizeResponse(...arguments);
  }

  serialize(snapshot) {
    const data = super.serialize(snapshot);
    if (Object.keys(data).includes('metadata')) {
      data.metadata = encodeString(data.metadata);
    }
    return data;
  }
}
