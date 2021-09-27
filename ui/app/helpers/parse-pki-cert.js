import { helper } from '@ember/component/helper';
import { pki } from 'node-forge';

export function parsePkiCert([model]) {
  // model has to be the responseJSON from PKI serializer
  if (!model.certificate) {
    return;
  }
  const cert = pki.certificateFromPem(model.certificate);
  const commonName = cert.subject.getField('CN') ? cert.subject.getField('CN').value : null;
  const issueDate = cert.validity.notBefore;
  const expiryDate = cert.validity.notAfter;
  return {
    common_name: commonName,
    issue_date: issueDate,
    expiry_date: expiryDate,
  };
}

export default helper(parsePkiCert);
