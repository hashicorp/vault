import { helper } from '@ember/component/helper';
import { pki } from 'node-forge';

export function parsePkiCert([model]) {
  // model has to be the responseJSON from PKI serializer
  if (!model.certificate) {
    return;
  }
  let cert;
  // node-forge cannot parse EC (elliptical curve) certs
  // return original response if unable to convert a Forge cert from PEM
  try {
    cert = pki.certificateFromPem(model.certificate);
  } catch (error) {
    return model;
  }
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
