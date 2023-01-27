import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate, CertificateRevocationList, CertificateChainValidationEngine } from 'pkijs';
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

export function parseCertificate(certificateContent) {
  let cert;
  try {
    const cert_base64 = certificateContent.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, '');
    const cert_der = fromBase64(cert_base64);
    const cert_asn1 = asn1js.fromBER(stringToArrayBuffer(cert_der));
    cert = new Certificate({ schema: cert_asn1.result });
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

const rootCertRaw = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUF21YbWKDM9VE4DWQkjZgO44J6zkwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgxMB4XDTIzMDEyNzE3NDU0\nOVoXDTIzMDIyODE3NDYxOVowHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgx\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzvM3XbISerXJkDDBKVNZ\nze57AKqcYWX0JJG4bAaDuYrRHAYGC0tl0K+RxEm/QVRe3wveADjBl0cewxl8CdJJ\nBgiTWe4UmrPTJkYdstpPZIRDm9NQ+gsZHsDqTH2ffEWUUyhZz/B6wE8uwF+qD/wA\nkyJkfx0wCOl4NtWjxjCVjygNbUJ2xwRCILwvLZdEPWfh+1eaBjkblWdohWTh7DjH\nVZxf6wCs4qBb45vfT35im5UAzwWBPDrZW/4MgyI+6G6HO8wF7tOcLKHdupIUN55n\n05jYiOcRoAJrVdbkj/s+r1pCyPoMnd1CvKxGlUUiVghHf4VH20/aVGe6J2pjUIg9\nMQIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQU+++fiUOshjrT/O7uN9d46h2pYWkwHwYDVR0jBBgwFoAU+++fiUOshjrT\n/O7uN9d46h2pYWkwDQYJKoZIhvcNAQELBQADggEBAKGT+gKEcK9eAwfEdbJTc0UW\nNq/8Q69L1yjriMhdpRpNKreTeQrpRXTJA4EwM3gDCQfkwP5VIzMVTaXFc77c8Cy6\n6hewH+8b/AbdZAuKbaN17zU4voQCUw7+FyS8Dna1bw19twSn+myNVX9WnvF6JH6c\nqA9qIvdePEjJvpgaAKdGFQ45iF8X8/09Azr2SO7Z2Z5ow7a69Fm3XGtIUhiKC0Vm\nYyyuPPKgaquyh6NxYNIVCiVjij63qBgfcUALbk1WmMKZCAQ5dVYQgbXoIrQYHHdL\nCn9wjyvku8X7A0CrQ2Dy50UT4aEyApiVrVQc4kJNXUJgzmgCAUw2Vo5ZqPUVcbc=\n-----END CERTIFICATE-----\n`;
const intCert = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUQ1IaPiRd6/WHoiut8eoxQpQU+HEwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgxMB4XDTIzMDEyNzE3NDU1\nMFoXDTIzMDIyODE3NDYyMFowHTEbMBkGA1UEAxMSU2hvcnQtTGl2ZWQgSW50IFIx\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx2iFrAXAGh5F+deqOlbo\nUk7g96cZjomFs0rLkxAaxJfu4uRF5PRLa9MSr9IsoQ7psXN8Di1LbWw11WwTjmcR\nMZ0f4piFXIT6EgMU6wZIKq+KjNJKffaTbJo9uCeqmr3eRP8pQqwKirqoukVamHYe\nJ9QciQEqydmj6QqO0i277Q8Ag1kThwGxq/fzxcmOSdbc38Vipu4qMHFU/t73GyxD\nMAQBbIpWd19noYY5KGQqcb7EmYajSMaiVXb65oL+SIjQb1DZQOTvMINWuZ7YwaM9\ngVRhnL/Ixb8Vj8yrJgpaUvs2BbwdoNvKH4S6qUHI2EmLn3zilSlEXq+5tEvT2Y/P\nawIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQU429qOoTfUG9C5BJ9Y95+aR/8gF4wHwYDVR0jBBgwFoAU+++fiUOshjrT\n/O7uN9d46h2pYWkwDQYJKoZIhvcNAQELBQADggEBAGVAgNu0lAWlfg0jD9TEJnOJ\n5/bofcGLb8keZjrw1jZlM7h6Sx+THcoZB9OxJkc4Kcg3dtVEMN4xi6Ypsa2rtFSm\nwwicMTeZhuA/0bwvZV80U1lJ7f45EFvgnlkF/yhYjh3GxLhAfY3tOCJZEvRw9iUL\n0cZF49SMWPME3hBNypv2CxJUHeqWH0IgffhyfYb/lFPRyii1B+uZqOOJbbzVJreR\n7mnHCmG/wblps1fJD4uOKMAdHcubM8LXwfOxib7st8chPVEIGzSfh8Lm3H6FqmqT\nLbTEpx4+jfEy9zGZalUP7NoERQ1E4TUv8zaz0PPTvjrLHNAxKzikDUgWiugoV34=\n-----END CERTIFICATE-----\n`;
const oldLeafCertRaw = `-----BEGIN CERTIFICATE-----\nMIIDTjCCAjagAwIBAgIUISxJR6kWmDj8Wo42lW1IeqoP2wQwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgyMB4XDTIzMDEyNzIxMTIy\nNVoXDTIzMDEyNzIyMTI1NVowFzEVMBMGA1UEAxMMb2xkLWxlYWYuY29tMIIBIjAN\nBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr90tuJJ6W9ZTztQpTsI7jvhF8OYv\norvPCnl5S7S1czChJtVnCPQ0NkFySlJgH8QYazIF+iZzKgK0kYgkXCHWPvJ8J6X5\ndCcGaXz+g51rWMPPU4kz7gwo2BSH7Ds14bj9QyQyJW62xdom2R3Pqf8gc3i5H85S\nPrAXYv610ZonxNm3S1Dj6qPhpsJxgrWxpGZUnP5ctRR+fqJhmwinifRblo2SmyQA\n9HDmYwSFPBiTc1tgLqK3R2cQdaRvGr0FrJm7mQpn6IX8bIR2h3OHOaw9idX/xoKJ\nKsdaUN7xH1xGov8NW2WXoRyOn1dxkXaNOB773k5QIJKhqS5W/fB8UerfjQIDAQAB\no4GLMIGIMA4GA1UdDwEB/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYB\nBQUHAwIwHQYDVR0OBBYEFO7NfpqE8mE+LjcLhwmYjIV64xXpMB8GA1UdIwQYMBaA\nFClhY4YNoXINAOx5KiGE3ymNbe6iMBcGA1UdEQQQMA6CDG9sZC1sZWFmLmNvbTAN\nBgkqhkiG9w0BAQsFAAOCAQEAJjpfvUKy0f+v5p5/9wOQjIur7FAyq1wLKihkdak4\nHQgQtX79bhCJA02vZzCUobLktDQf0AhinBs6ofr7gub0iI92Xm7OurDlerwnsbtm\nhtuzJ9zmtoAimJSBOEuMcJSX4WtJq9efjhW5+gJvGP/FBuEMNd0l6HxUDgKWAq9Y\nVMFcgQJDI71I9pV5spacCSPoX86kuR6WH0zbXVpP6hIltFqYvyUh4vQ4o2Nsai/P\nCtAVlJEqS4Re8YhqPA3FHcljtjbTOOovWbkOKnybtIeLZn6stLfbrRrYpPsxnRWJ\nWRfGEmVygqmMC5WFCaJ5vioIMosR82TeYAcwHE+eaDNuZA==\n-----END CERTIFICATE-----\n`;
const intCertCrl = `-----BEGIN X509 CRL-----\nMIIBlzCBgAIBATANBgkqhkiG9w0BAQsFADAdMRswGQYDVQQDExJTaG9ydC1MaXZl\nZCBJbnQgUjEXDTIzMDEyNzE5MjgxOFoXDTIzMDEzMDE5MjgxOFqgLzAtMB8GA1Ud\nIwQYMBaAFONvajqE31BvQuQSfWPefmkf/IBeMAoGA1UdFAQDAgEDMA0GCSqGSIb3\nDQEBCwUAA4IBAQCh+/9aCKGAH/W3YFXKLIAdtugzMRImVLJZH7R+tiJ80gjYsL1g\n9moP8W4DTt31LZJVPpkceLEqw9glKiNmsh5kXm+/9cV9E2zjHIcI4fJdmDSw6RWI\n/aIcxHJmQ5nrhjptBuhdpGmqI5RA1omuqfYt5Gfysa3EHqPBO/VgU/nTvSEGGkJw\nWFdWpV1ncZv995YeJU1Zx3N3TBW5nLGl0McJt0fREHLv/G9Cj2PGSEQxKNCAkbt+\nhHDnVfw4paFKoeW4txMMPi2UKvkN0ypjQ+NDYrGVsCJz/8ghdwNWU0NBtJEWyjcI\nf8LVIsAE4aniWBPN/EVqZWwYMER08UxgINDM\n-----END X509 CRL-----\n`;
const newIntCertRaw = `-----BEGIN CERTIFICATE-----\nMIIDKzCCAhOgAwIBAgIUe6Piv67pg7ZgV4EYpADbG8O2SGAwDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSTG9uZy1MaXZlZCBSb290IFgyMB4XDTIzMDEyNzE5Mjc0\nOFoXDTIzMDIyODE5MjgxOFowHTEbMBkGA1UEAxMSU2hvcnQtTGl2ZWQgSW50IFIx\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx2iFrAXAGh5F+deqOlbo\nUk7g96cZjomFs0rLkxAaxJfu4uRF5PRLa9MSr9IsoQ7psXN8Di1LbWw11WwTjmcR\nMZ0f4piFXIT6EgMU6wZIKq+KjNJKffaTbJo9uCeqmr3eRP8pQqwKirqoukVamHYe\nJ9QciQEqydmj6QqO0i277Q8Ag1kThwGxq/fzxcmOSdbc38Vipu4qMHFU/t73GyxD\nMAQBbIpWd19noYY5KGQqcb7EmYajSMaiVXb65oL+SIjQb1DZQOTvMINWuZ7YwaM9\ngVRhnL/Ixb8Vj8yrJgpaUvs2BbwdoNvKH4S6qUHI2EmLn3zilSlEXq+5tEvT2Y/P\nawIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV\nHQ4EFgQU429qOoTfUG9C5BJ9Y95+aR/8gF4wHwYDVR0jBBgwFoAUKWFjhg2hcg0A\n7HkqIYTfKY1t7qIwDQYJKoZIhvcNAQELBQADggEBADRTZICB4qi6ZOp3CozVa1en\nf5gXxSllm2hsRy15kWhPsfRm+rsRz8wL1Ciwk/HrXVuEAZCyYevogZ5BaLvRygeY\na+ejKUWbuYqhLmGh3hLTCxsJGOGKksXvIXmuN7gxbfTO6/2MXmNaTt8ZVSKnn25F\nNkmTWpcxfAh92SHE6c/9MY9cRbT4EEegqMSg6wpxk7kH7D44444LZDJce48gHV1n\nFz8G/n67lfV4P4wA5HmYL4aDhjqqeNNGH/IwCweFOX1oXwLtD6/L8eqeriViEp4T\nGZTBfw0fgfv2Wo8+IsOOGmLgSmAPzSWB/CsSqWmEX8JB5n2XUkhLRNOiL0jV/uc=\n-----END CERTIFICATE-----\n`;
const newIntCertCrl = `-----BEGIN X509 CRL-----\nMIIBlzCBgAIBATANBgkqhkiG9w0BAQsFADAdMRswGQYDVQQDExJTaG9ydC1MaXZl\nZCBJbnQgUjEXDTIzMDEyNzE5MjgxOFoXDTIzMDEzMDE5MjgxOFqgLzAtMB8GA1Ud\nIwQYMBaAFONvajqE31BvQuQSfWPefmkf/IBeMAoGA1UdFAQDAgEDMA0GCSqGSIb3\nDQEBCwUAA4IBAQCh+/9aCKGAH/W3YFXKLIAdtugzMRImVLJZH7R+tiJ80gjYsL1g\n9moP8W4DTt31LZJVPpkceLEqw9glKiNmsh5kXm+/9cV9E2zjHIcI4fJdmDSw6RWI\n/aIcxHJmQ5nrhjptBuhdpGmqI5RA1omuqfYt5Gfysa3EHqPBO/VgU/nTvSEGGkJw\nWFdWpV1ncZv995YeJU1Zx3N3TBW5nLGl0McJt0fREHLv/G9Cj2PGSEQxKNCAkbt+\nhHDnVfw4paFKoeW4txMMPi2UKvkN0ypjQ+NDYrGVsCJz/8ghdwNWU0NBtJEWyjcI\nf8LVIsAE4aniWBPN/EVqZWwYMER08UxgINDM\n-----END X509 CRL-----\n`;
const newLeafCertRaw = `-----BEGIN CERTIFICATE-----\nMIIDTjCCAjagAwIBAgIULwsmGqV/hjRMjXkwfVbWYQD3cWowDQYJKoZIhvcNAQEL\nBQAwHTEbMBkGA1UEAxMSU2hvcnQtTGl2ZWQgSW50IFIxMB4XDTIzMDEyNzIxMTky\nMVoXDTIzMDEyNzIyMTk1MVowFzEVMBMGA1UEAxMMbmV3LWxlYWYuY29tMIIBIjAN\nBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0o1yb0nnM4vyf8YpPoR4LzpHDJzU\n4eILsm9a0CEqT/9cZj+fkZY8txZ9lF6inyrM9uthym8Y8obzRF8MNfo2N1+BB0vG\nLcQzZG9FMLNsCEJ61QojkYCHvOFv5M25qBCwbe6oobhfndySl0GN7l63ixXYOQ/6\nDPWvixey1gdn8SZvpu4y+Kk7ggH5TQucIpkPUfvpcEeNdIzMAbfTlet7rtE2fp6Z\nIYJQjqiTpTCBKzGN+y5FPTcdkabNLi2887A/bSaTPgMbMwb/p8+OLxqRHAXI7mO6\nqeUciSeCMQVkBnyLfo9jsqbvj2EC1Bes9NdiOIliyYEgNezovshYDeEtKwIDAQAB\no4GLMIGIMA4GA1UdDwEB/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYB\nBQUHAwIwHQYDVR0OBBYEFHE0pwD2PoZmmhbcH/qHK4WpRvU5MB8GA1UdIwQYMBaA\nFONvajqE31BvQuQSfWPefmkf/IBeMBcGA1UdEQQQMA6CDG5ldy1sZWFmLmNvbTAN\nBgkqhkiG9w0BAQsFAAOCAQEAgz/J8gT9x9fdE371YvkVTgci1kFObSLjr3Jts9Sn\nhp9Yn0aEOmak24FD9vJgQHsDzGYKYsnSo1IvCSuQ6n9s1PuK3bjdxxGr6P/r5yEr\nSKBaRjFWmVzIMKKFilDyViveOFZBoxMZLAzgO8M2uiRkHpdRXAKCViyDtLvLjR9H\nfpSxx++VfhLE2HDgraWtVcmC4r9GGdLzQMf3LcMC3/g+VcFzwLVHkdqQdO8h9o68\nD24u2X/XYsUIJlGnPBJ3JGm+DEeqKJaBFIdvxyPFxkheYngghsf/j15QfD/E8df2\n3ihmvNxS66QzqWY8zVI0zJEkjldk1BCfMgTfDlYjGaK7/w==\n-----END CERTIFICATE-----`;

const newParentMountCrl = `-----BEGIN X509 CRL-----\nMIIBlzCBgAIBATANBgkqhkiG9w0BAQsFADAdMRswGQYDVQQDExJMb25nLUxpdmVk\nIFJvb3QgWDIXDTIzMDEyNzE3NDYxOVoXDTIzMDEzMDE3NDYxOVqgLzAtMB8GA1Ud\nIwQYMBaAFClhY4YNoXINAOx5KiGE3ymNbe6iMAoGA1UdFAQDAgEBMA0GCSqGSIb3\nDQEBCwUAA4IBAQAizl1tvpc3dc6cmHnDREr66U1tX4Tun92foaVU/jaY4NoSjP3l\naJmMTsYO39dubOIkt4DSXcqK8YoMMwX1FN+VIkmTWmkKavYuUgnQob1sWfDN8PGR\nshQ6wla3C351rUYRU9xei+MVLdIXBm+CnqpNFGHqGOtUiVDUfyRWzNbXNbq8ORv5\nQwTv7ujDgEhBZYm6wGZPzmLDbfxajgop/BYJMuduPXAXEOP48lB6Lmz0DUqkOlSc\na+xCvUBu2FcFlkMpyTrgyGQiJtio7JqZ/S0rswSake3ADYWb5KibUikBa6+2B08H\n/BIiNwxKRiGKUyLJhbDYUj5CeOV2pF0m2Knd\n-----END X509 CRL-----\n`;
const oldParentMountCrl = `-----BEGIN X509 CRL-----\nMIIBlzCBgAIBATANBgkqhkiG9w0BAQsFADAdMRswGQYDVQQDExJMb25nLUxpdmVk\nIFJvb3QgWDEXDTIzMDEyNzE3NDYxOVoXDTIzMDEzMDE3NDYxOVqgLzAtMB8GA1Ud\nIwQYMBaAFPvvn4lDrIY60/zu7jfXeOodqWFpMAoGA1UdFAQDAgEDMA0GCSqGSIb3\nDQEBCwUAA4IBAQATD1GQuwxyF0AN527e7Zr3at5qgblPZdnkv/A2FwO3w9aee0fb\nvey9jzCKJP+RC9WaUMU0F1j3fh9wJnik1BKdjQRyKUgqMvE/JicVvoZf47yBMYZJ\n3EoH8SzHXoBmau3a1tjRZgZGSybbt4O2wpGg7gTBxyTRCYkni7YhHUUbcVB22A6F\nqmJI36pm2tnvQY5kxLJtx57j7TppJ9FxNCKOMlnn4R/UIFOZ1wGAxZfXeYSEcvJX\njfHI0eK7nzcHQS+sITnuLeIYpdrXWXt8YpURSvxrZjFC5i4oS/7ZS0TZ3m403hku\nQbz1nxncaMZ3FS9nVm4YbuVao6p18bFhckQM\n-----END X509 CRL-----\n`;
function jsonToParsableObject(jsonString) {
  const content = jsonString.includes('CRL')
    ? { base64: jsonString?.replace(/(-----(BEGIN|END) X509 CRL-----|\n)/g, ''), type: 'csr' }
    : { base64: jsonString?.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, ''), type: 'cert' };
  const der = fromBase64(content.base64);
  const asn1 = asn1js.fromBER(stringToArrayBuffer(der));
  switch (content.type) {
    case 'cert':
      return new Certificate({ schema: asn1.result });
    case 'csr':
      return new CertificateRevocationList({ schema: asn1.result });
    default:
      return { can_parse: false };
  }
}
async function parseCRL() {
  const rootCa = jsonToParsableObject(rootCertRaw);
  const intermediateCa = jsonToParsableObject(intCert);
  const oldLeafCert = jsonToParsableObject(oldLeafCertRaw);
  const newIntCert = jsonToParsableObject(newIntCertRaw);
  const newLeafCert = jsonToParsableObject(newLeafCertRaw);
  const crlOldInt = jsonToParsableObject(intCertCrl);
  const crlNewInt = jsonToParsableObject(newIntCertCrl);

  const chainEngine = new CertificateChainValidationEngine({
    certs: [rootCa, intermediateCa, newLeafCert],
    crls: [
      crlOldInt,
      crlNewInt,
      jsonToParsableObject(newParentMountCrl),
      jsonToParsableObject(oldParentMountCrl),
    ],
    trustedCerts: [rootCa],
  });

  const chain = await chainEngine.verify();
  console.log(chain);
  // export const chainValidation()
}
