import { helper } from '@ember/component/helper';
import { pki } from 'node-forge';

export function parsePkiCert([model]) {
  // model has to be the responseJSON from PKI serializer
  if (!model.certificate) {
    throw new Error();
  }
  const cert = pki.certificateFromPem(model.certificate);
  const commonName = cert.subject.getField('CN').value;
  const issueDate = cert.validity.notBefore;
  const expiryDate = cert.validity.notAfter;
  const certMetadata = {
    common_name: commonName,
    issue_date: issueDate,
    expiry_date: expiryDate,
  };
  return certMetadata;
}

export default helper(parsePkiCert);
