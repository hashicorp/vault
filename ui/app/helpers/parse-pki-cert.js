import { helper } from '@ember/component/helper';
import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate } from 'pkijs';
import { getUnixTime } from 'date-fns';

/*
  We wish to get these OID_VALUES out of this certificate's subject. A
  subject is a list of RDNs, where each RDN is a (type, value) tuple
  and where a type is an OID. The OID for CN can be found here:
     
     https://datatracker.ietf.org/doc/html/rfc5280#page-112
  
  Each value is then encoded as another ASN.1 object; in the case of a
  CommonName field, this is usually a PrintableString, BMPString, or a
  UTF8String. Regardless of encoding, it should be present in the
  valueBlock's value field if it is renderable.
*/

const OID_VALUES = {
  common_name: '2.5.4.3', // http://oid-info.com/get/2.5.4.3
  serial_number: '2.5.4.5', // http://oid-info.com/get/2.5.4.5
  ou: '2.5.4.11',
  organization: '2.5.4.10',
  country: '2.5.4.6',
  locality: '2.5.4.7',
  province: '2.5.4.8',
  street_address: '2.5.4.9',
  postal_code: '2.5.4.17',
};

export function parseCertificate(certificateContent) {
  let cert;
  try {
    const cert_base64 = certificateContent.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, '');
    const cert_der = fromBase64(cert_base64);
    const cert_asn1 = asn1js.fromBER(stringToArrayBuffer(cert_der));
    cert = new Certificate({ schema: cert_asn1.result });
  } catch (error) {
    console.debug('DEBUG: Parsing Certificate', error); // eslint-disable-line
    return {
      can_parse: false,
    };
  }

  // Date instances are stored in the value field as the notAfter/notBefore
  // field themselves are Time values.
  const expiryDate = cert?.notAfter?.value;
  const issueDate = cert?.notBefore?.value;
  return {
    ...parseSubject(cert?.subject?.typesAndValues),
    can_parse: true,
    expiry_date: expiryDate, // remove along with old PKI work
    issue_date: issueDate, // remove along with old PKI work
    not_valid_after: getUnixTime(expiryDate),
    not_valid_before: getUnixTime(issueDate),
  };
}

// parses subject and returns value for each key in OID_VALUES
function parseSubject(subject) {
  const returnValues = (OID) => {
    const values = subject.filter((rdn) => rdn?.type === OID).map((rdn) => rdn?.value?.valueBlock?.value);
    // Theoretically, there might be multiple (or no) CommonNames -- but Vault
    // presently refuses to issue certificates without CommonNames in most
    // cases. For now, return the first CommonName we find. Alternatively, we
    // might update our callers to handle multiple, or join them using some
    // separator like ','.
    return values ? (values.length ? values[0] : null) : null;
  };

  const subjectValues = {};
  Object.keys(OID_VALUES).forEach((key) => (subjectValues[key] = returnValues(OID_VALUES[key])));
  return subjectValues;
}

export function parsePkiCert([model]) {
  // model has to be the responseJSON from PKI serializer
  // return if no certificate or if the "certificate" is actually a CRL
  if (!model.certificate || model.certificate.includes('BEGIN X509 CRL')) {
    return;
  }
  return parseCertificate(model.certificate);
}

export default helper(parsePkiCert);
