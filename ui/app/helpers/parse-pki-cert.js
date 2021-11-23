import { helper } from '@ember/component/helper';
import { pki } from 'node-forge';

export function parsePkiCert([model]) {
  // model has to be the responseJSON from PKI serializer
  if (!model.certificate) {
    return;
  }
  let cert;
  // node-forge cannot parse EC (elliptical curve) certs
  // set canParse to false if unable to convert a Forge cert from PEM
  try {
    cert = pki.certificateFromPem(model.certificate);
  } catch (error) {
    return {
      can_parse: false,
    };
  }
  const commonName = cert?.subject.getField('CN') ? cert.subject.getField('CN').value : null;
  const expiryDate = cert?.validity.notAfter;
  const issueDate = cert?.validity.notBefore;
  return {
    can_parse: true,
    common_name: commonName,
    expiry_date: expiryDate,
    issue_date: issueDate,
  };
}

export default helper(parsePkiCert);
