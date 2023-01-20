import { helper } from '@ember/component/helper';
import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate } from 'pkijs';
import { differenceInHours, getUnixTime } from 'date-fns';

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

  const subjectParams = parseSubject(cert?.subject?.typesAndValues);
  const expiryDate = cert?.notAfter?.value;
  const issueDate = cert?.notBefore?.value;

  // TODO wrap to catch errors, only parse if cross-signing?
  // for cross-signing
  const { alt_names, uri_sans } = parseExtensions(cert?.extensions);
  const [signature_bits, use_pss] = mapSignatureBits(cert?.signatureAlgorithm);
  const exclude_cn_from_sans =
    alt_names?.length > 0 && !alt_names?.includes(subjectParams?.common_name) ? true : false;
  const ttl = `${differenceInHours(expiryDate, issueDate)}h`;

  return {
    ...subjectParams,
    can_parse: true,
    expiry_date: expiryDate, // remove along with old PKI work
    issue_date: issueDate, // remove along with old PKI work
    not_valid_after: getUnixTime(expiryDate),
    not_valid_before: getUnixTime(issueDate),
    alt_names: alt_names?.join(', '),
    uri_sans: uri_sans?.join(', '),
    signature_bits,
    use_pss,
    exclude_cn_from_sans,
    ttl,
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

//* PARSING HELPERS [lookup OIDs: http://oid-info.com/basic-search.htm]

const SUBJECT_OIDs = {
  common_name: '2.5.4.3',
  serial_number: '2.5.4.5',
  ou: '2.5.4.11',
  organization: '2.5.4.10',
  country: '2.5.4.6',
  locality: '2.5.4.7',
  province: '2.5.4.8',
  street_address: '2.5.4.9',
  postal_code: '2.5.4.17',
};
const EXTENSION_OIDs = {
  key_usage: '2.5.29.15',
  subject_alt_name: '2.5.29.17',
};
// SubjectAltName/GeneralName types (scroll up to page 38) https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.7
const SAN_TYPES = {
  alt_names: 2, // dNSName
  uri_sans: 6, // uniformResourceIdentifier
  ip_sans: 7, // iPAddress - OCTET STRING
};
const SIGNATURE_ALGORITHM_OIDs = {
  '1.2.840.113549.1.1.2': '0', // MD2-RSA
  '1.2.840.113549.1.1.4': '0', // MD5-RSA
  '1.2.840.113549.1.1.5': '0', // SHA1-RSA
  '1.2.840.113549.1.1.11': '256', // SHA256-RSA
  '1.2.840.113549.1.1.12': '384', // SHA384-RSA
  '1.2.840.113549.1.1.13': '512', // SHA512-RSA
  '1.2.840.113549.1.1.10': {
    // RSA-PSS have additional OIDs that need to be mapped
    '2.16.840.1.101.3.4.2.1': '256', // SHA-256
    '2.16.840.1.101.3.4.2.2': '384', // SHA-384
    '2.16.840.1.101.3.4.2.3': '512', // SHA-512
  },
  '1.2.840.10040.4.3': '0', // DSA-SHA1
  '2.16.840.1.101.3.4.3.2': '256', // DSA-SHA256
  '1.2.840.10045.4.1': '0', // ECDSA-SHA1
  '1.2.840.10045.4.3.2': '256', // ECDSA-SHA256
  '1.2.840.10045.4.3.3': '384', // ECDSA-SHA384
  '1.2.840.10045.4.3.4': '512', // ECDSA-SHA512
  '1.3.101.112': '0', // Ed25519
};

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
function parseSubject(subject) {
  if (!subject) return {};
  const returnValues = (OID) => {
    const values = subject.filter((rdn) => rdn?.type === OID).map((rdn) => rdn?.value?.valueBlock?.value);
    // Theoretically, there might be multiple (or no) CommonNames -- but Vault
    // presently refuses to issue certificates without CommonNames in most
    // cases. For now, return the first CommonName we find. Alternatively, we
    // might update our callers to handle multiple and return a string array
    return values ? (values?.length ? values[0] : null) : null;
  };
  const subjectValues = {};
  Object.keys(SUBJECT_OIDs).forEach((key) => (subjectValues[key] = returnValues(SUBJECT_OIDs[key])));
  return subjectValues;
}

function parseExtensions(extensions) {
  if (!extensions) return {};

  const values = {};
  for (const attrName in EXTENSION_OIDs) {
    values[attrName] = extensions.find((ext) => ext?.extnID === EXTENSION_OIDs[attrName])?.parsedValue;
  }

  if (values.subject_alt_name) {
    for (const attrName in SAN_TYPES) {
      values[attrName] = values.subject_alt_name?.altNames
        .filter((gn) => gn.type === Number(SAN_TYPES[attrName]))
        .map((gn) => gn.value);
    }
  }

  if (values.ip_sans) {
    // TODO parse octet string for IP addresses
  }

  if (values.key_usage) {
    // TODO parse key_usage
  }

  delete values.subject_alt_name;
  return values;
  /*
  values is an object with keys from EXTENSION_OIDs and SAN_TYPES
  values = {
    "key_usage": BitString
    "alt_names": string[],
    "uri_sans": string[],
    "ip_sans": OctetString[], <- currently array of OctetStrings to-be-parsed
  }
  */
}

function mapSignatureBits(sigAlgo) {
  const { algorithmId } = sigAlgo;

  // use_pss is true, additional OIDs need to be mapped
  if (algorithmId === '1.2.840.113549.1.1.10') {
    // object identifier for PSS is very nested
    const objId = sigAlgo.algorithmParams?.valueBlock?.value[0]?.valueBlock?.value[0]?.valueBlock?.value[0]
      .toString()
      .split(' : ')[1];
    return [SIGNATURE_ALGORITHM_OIDs[algorithmId][objId], true];
  }
  return [SIGNATURE_ALGORITHM_OIDs[algorithmId], false];
}

export default helper(parsePkiCert);
