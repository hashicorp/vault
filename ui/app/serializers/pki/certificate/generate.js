import { parseCertificate } from 'vault/helpers/parse-pki-cert';
import ApplicationSerializer from '../../application';

export default class PkiCertificateGenerateSerializer extends ApplicationSerializer {
  primaryKey = 'serial_number';

  serialize() {
    const json = super.serialize(...arguments);
    delete json.name;
    return json;
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
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
}
