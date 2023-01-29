import { helper } from '@ember/component/helper';
import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Convert } from 'pvtsutils';
import { Certificate } from 'pkijs';

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
  // We wish to get the CN element out of this certificate's subject. A
  // subject is a list of RDNs, where each RDN is a (type, value) tuple
  // and where a type is an OID. The OID for CN can be found here:
  //
  //    http://oid-info.com/get/2.5.4.3
  //    https://datatracker.ietf.org/doc/html/rfc5280#page-112
  //
  // Each value is then encoded as another ASN.1 object; in the case of a
  // CommonName field, this is usually a PrintableString, BMPString, or a
  // UTF8String. Regardless of encoding, it should be present in the
  // valueBlock's value field if it is renderable.
  const commonNameOID = '2.5.4.3';
  const commonNames = cert?.subject?.typesAndValues
    .filter((rdn) => rdn?.type === commonNameOID)
    .map((rdn) => rdn?.value?.valueBlock?.value);

  // Theoretically, there might be multiple (or no) CommonNames -- but Vault
  // presently refuses to issue certificates without CommonNames in most
  // cases. For now, return the first CommonName we find. Alternatively, we
  // might update our callers to handle multiple, or join them using some
  // separator like ','.
  const commonName = commonNames ? (commonNames.length ? commonNames[0] : null) : null;

  // Date instances are stored in the value field as the notAfter/notBefore
  // field themselves are Time values.
  const expiryDate = cert?.notAfter?.value;
  const issueDate = cert?.notBefore?.value;
  const serialNumber = Convert.ToHex(cert.serialNumber.valueBlock.valueHex)
    .match(/.{1,2}/g)
    .join(':');
  return {
    can_parse: true,
    common_name: commonName,
    serial_number: serialNumber,
    expiry_date: expiryDate,
    issue_date: issueDate,
    not_valid_after: expiryDate.valueOf(),
    not_valid_before: issueDate.valueOf(),
  };
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
