import { parseCertificate } from 'vault/helpers/parse-pki-cert';
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
      // convert issueDate to same format as other date values
      // this can be moved into the parseCertificate helper once the old pki implementation is removed
      if (parsedCert.issue_date) {
        parsedCert.issue_date = parsedCert.issue_date.valueOf();
      }
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
