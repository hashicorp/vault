import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate } from 'pkijs';
import { differenceInHours, getUnixTime } from 'date-fns';
import {
  EXTENSION_OIDs,
  SUBJECT_OIDs,
  IGNORED_OIDs,
  SAN_TYPES,
  SIGNATURE_ALGORITHM_OIDs,
} from './parse-pki-cert-oids';

/* 
 It may be helpful to visualize a certificate's SEQUENCE structure alongside this parsing file.
 You can do so by decoding a certificate here: https://lapo.it/asn1js/#

 A certificate is encoded in ASN.1 data - a SEQUENCE is how you define structures in ASN.1.
 GeneralNames, Extension, AlgorithmIdentifier are all examples of SEQUENCEs 

 * Error handling: 
{ can_parse: false } -> returned if the external library cannot convert the certificate 
{ parsing_errors: [] } -> returned if the certificate was converted, but there's ANY problem parsing certificate details. 
 This means we cannot cross-sign in the UI and prompt the user to do so manually using the CLI.
 */

export function jsonToCert(jsonString) {
  const cert_base64 = jsonString.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, '');
  const cert_der = fromBase64(cert_base64);
  const cert_asn1 = asn1js.fromBER(stringToArrayBuffer(cert_der));
  return new Certificate({ schema: cert_asn1.result });
}

export function parseCertificate(certificateContent) {
  let cert;
  try {
    cert = jsonToCert(certificateContent);
  } catch (error) {
    console.debug('DEBUG: Converting Certificate', error); // eslint-disable-line
    return { can_parse: false };
  }

  let parsedCertificateValues;
  try {
    const subjectValues = parseSubject(cert?.subject?.typesAndValues);
    const extensionValues = parseExtensions(cert?.extensions);
    const [signature_bits, use_pss] = mapSignatureBits(cert?.signatureAlgorithm);
    const formattedValues = formatValues(subjectValues, extensionValues);
    parsedCertificateValues = { ...formattedValues, signature_bits, use_pss };
  } catch (error) {
    console.debug('DEBUG: Parsing Certificate', error); // eslint-disable-line
    parsedCertificateValues = { parsing_errors: [new Error('error parsing certificate values')] };
  }

  const expiryDate = cert?.notAfter?.value;
  const issueDate = cert?.notBefore?.value;
  const ttl = `${differenceInHours(expiryDate, issueDate)}h`;

  return {
    ...parsedCertificateValues,
    can_parse: true,
    expiry_date: expiryDate, // remove along with old PKI work
    issue_date: issueDate, // remove along with old PKI work
    not_valid_after: getUnixTime(expiryDate),
    not_valid_before: getUnixTime(issueDate),
    ttl,
  };
}

export function parsePkiCert(model) {
  // model has to be the responseJSON from PKI serializer
  // return if no certificate or if the "certificate" is actually a CRL
  if (!model.certificate || model.certificate.includes('BEGIN X509 CRL')) {
    return;
  }
  return parseCertificate(model.certificate);
}

export function formatValues(subject, extension) {
  if (!subject || !extension) {
    return { parsing_errors: [new Error('error formatting certificate values')] };
  }
  const { subjValues, subjErrors } = subject;
  const { extValues, extErrors } = extension;
  const parsing_errors = [...subjErrors, ...extErrors];
  const exclude_cn_from_sans =
    extValues.alt_names?.length > 0 && !extValues.alt_names?.includes(subjValues?.common_name) ? true : false;
  // now that we've finished parsing data, join all extension arrays
  for (const ext in extValues) {
    if (Array.isArray(extValues[ext])) {
      extValues[ext] = extValues[ext].length !== 0 ? extValues[ext].join(', ') : null;
    }
  }

  // TODO remove this deletion when key_usage is parsed, update test
  delete extValues.key_usage;
  return {
    ...subjValues,
    ...extValues,
    parsing_errors,
    exclude_cn_from_sans,
  };
}

//* PARSING HELPERS
/*
  We wish to get each SUBJECT_OIDs (see utils/parse-pki-cert-oids.js) out of this certificate's subject. 
  A subject is a list of RDNs, where each RDN is a (type, value) tuple
  and where a type is an OID. The OID for CN can be found here:
     
     https://datatracker.ietf.org/doc/html/rfc5280#page-112
  
  Each value is then encoded as another ASN.1 object; in the case of a
  CommonName field, this is usually a PrintableString, BMPString, or a
  UTF8String. Regardless of encoding, it should be present in the
  valueBlock's value field if it is renderable.
*/
export function parseSubject(subject) {
  if (!subject) return null;
  const values = {};
  const errors = [];
  if (subject.any((rdn) => !Object.values(SUBJECT_OIDs).includes(rdn.type))) {
    errors.push(new Error('certificate contains unsupported subject OIDs'));
  }
  const returnValues = (OID) => {
    const values = subject.filter((rdn) => rdn?.type === OID).map((rdn) => rdn?.value?.valueBlock?.value);
    // Theoretically, there might be multiple (or no) CommonNames -- but Vault
    // presently refuses to issue certificates without CommonNames in most
    // cases. For now, return the first CommonName we find. Alternatively, we
    // might update our callers to handle multiple and return a string array
    return values ? (values?.length ? values[0] : null) : null;
  };
  Object.keys(SUBJECT_OIDs).forEach((key) => (values[key] = returnValues(SUBJECT_OIDs[key])));
  return { subjValues: values, subjErrors: errors };
}

export function parseExtensions(extensions) {
  if (!extensions) return null;
  const values = {};
  const errors = [];
  const allowedOids = Object.values({ ...EXTENSION_OIDs, ...IGNORED_OIDs });
  if (extensions.any((ext) => !allowedOids.includes(ext.extnID))) {
    errors.push(new Error('certificate contains unsupported extension OIDs'));
  }

  // make each extension its own key/value pair
  for (const attrName in EXTENSION_OIDs) {
    values[attrName] = extensions.find((ext) => ext.extnID === EXTENSION_OIDs[attrName])?.parsedValue;
  }

  if (values.subject_alt_name) {
    // we only support SANs of type 2 (altNames), 6 (uri) and 7 (ipAddress)
    const supportedTypes = Object.values(SAN_TYPES);
    const supportedNames = Object.keys(SAN_TYPES);
    const sans = values.subject_alt_name?.altNames;
    if (!sans) {
      errors.push(new Error('certificate contains unsupported subjectAltName values'));
    } else if (sans.any((san) => !supportedTypes.includes(san.type))) {
      // pass along error that unsupported values exist
      errors.push(new Error('subjectAltName contains unsupported types'));
      // still check and parse any supported values
      if (sans.any((san) => supportedTypes.includes(san.type))) {
        supportedNames.forEach((attrName) => {
          values[attrName] = sans
            .filter((gn) => gn.type === Number(SAN_TYPES[attrName]))
            .map((gn) => gn.value);
        });
      }
    } else if (sans.every((san) => supportedTypes.includes(san.type))) {
      supportedNames.forEach((attrName) => {
        values[attrName] = sans.filter((gn) => gn.type === Number(SAN_TYPES[attrName])).map((gn) => gn.value);
      });
    } else {
      errors.push(new Error('unsupported subjectAltName values'));
    }
  }

  // permitted_dns_domains
  if (values.name_constraints) {
    // we only support Name Constraints of dnsName (type 2), this value lives in the permittedSubtree of the Name Constraints sequence
    // permittedSubtrees contain an array of subtree objects, each object has a 'base' key and EITHER a 'minimum' or 'maximum' key
    // GeneralSubtree { "base": {   "type": 2,  "value": "dnsname1.com" }, minimum: 0 }
    const nameConstraints = values.name_constraints;
    if (Object.keys(nameConstraints).includes('excludedSubtrees')) {
      errors.push(new Error('nameConstraints contains excludedSubtrees'));
    } else if (nameConstraints.permittedSubtrees.any((subtree) => subtree.minimum !== 0)) {
      errors.push(new Error('nameConstraints permittedSubtree contains non-zero minimums'));
    } else if (nameConstraints.permittedSubtrees.any((subtree) => subtree.maximum)) {
      errors.push(new Error('nameConstraints permittedSubtree contains maximum'));
    } else if (nameConstraints.permittedSubtrees.any((subtree) => subtree.base.type !== 2)) {
      errors.push(new Error('nameConstraints permittedSubtree can only contain dnsName (type 2)'));
      // still check and parse any supported values
      if (nameConstraints.permittedSubtrees.any((subtree) => subtree.base.type === 2)) {
        values.permitted_dns_domains = nameConstraints.permittedSubtrees
          .filter((gn) => gn.base.type === 2)
          .map((gn) => gn.base.value);
      }
    } else if (nameConstraints.permittedSubtrees.every((subtree) => subtree.base.type === 2)) {
      values.permitted_dns_domains = nameConstraints.permittedSubtrees.map((gn) => gn.base.value);
    } else {
      errors.push(new Error('unsupported nameConstraints values'));
    }
  }

  if (values.basic_constraints) {
    values.max_path_length = values.basic_constraints?.pathLenConstraint;
  }

  if (values.ip_sans) {
    // TODO parse octet string for IP addresses
  }

  if (values.key_usage) {
    // TODO parse key_usage
  }

  delete values.subject_alt_name;
  delete values.basic_constraints;
  delete values.name_constraints;
  return { extValues: values, extErrors: errors };
  /*
  values is an object with keys from EXTENSION_OIDs and SAN_TYPES
  values = {
    "alt_names": string[],
    "uri_sans": string[],
    "permitted_dns_domains": string[],
    "max_path_length": int,
    "key_usage": BitString, <- to-be-parsed
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
