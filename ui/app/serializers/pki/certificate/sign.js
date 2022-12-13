import { parseCertificate } from 'vault/helpers/parse-pki-cert';
import ApplicationSerializer from '../../application';

export default class PkiCertificateSignSerializer extends ApplicationSerializer {
  primaryKey = 'serial_number';

  serialize() {
    const json = super.serialize(...arguments);
    // role name is part of the URL, remove from payload
    delete json.name;
    return json;
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (requestType === 'createRecord' && payload.data.certificate) {
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
