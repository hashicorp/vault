/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate } from 'pkijs';
import { differenceInHours, getUnixTime } from 'date-fns';
import {
  EXTENSION_OIDs,
  IGNORED_OIDs,
  KEY_USAGE_BITS,
  SAN_TYPES,
  SIGNATURE_ALGORITHM_OIDs,
  SUBJECT_OIDs,
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

export function jsonToCertObject(jsonString) {
  const cert_base64 = jsonString.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, '');
  const cert_der = fromBase64(cert_base64);
  const cert_asn1 = asn1js.fromBER(stringToArrayBuffer(cert_der));
  return new Certificate({ schema: cert_asn1.result });
}

export function parseCertificate(certificateContent) {
  let cert;
  try {
    cert = jsonToCertObject(certificateContent);
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

  return {
    ...subjValues,
    ...extValues,
    parsing_errors,
    exclude_cn_from_sans,
  };
}

/*
How to use the verify function for cross-signing: 
(See setup script here: https://github.com/hashicorp/vault-tools/blob/main/vault-ui/pki/pki-cross-sign-config.sh)
1. A trust chain exists between "old-parent-issuer-name" -> "old-intermediate"
2. Cross-sign "old-intermediate" against "my-parent-issuer-name" creating a new certificate: "newly-cross-signed-int-name"
3. Generate a leaf certificate from "newly-cross-signed-int-name", let's call it "baby-leaf"
4. Verify that "baby-leaf" validates against both chains: 
"old-parent-issuer-name" -> "old-intermediate" -> "baby-leaf"
"my-parent-issuer-name" -> "newly-cross-signed-int-name" -> "baby-leaf"

A valid cross-signing would mean BOTH of the following return true:
verifyCertificates(oldParentCert, oldIntCert, leaf)
verifyCertificates(newParentCert, crossSignedCert, leaf)

each arg is the JSON string certificate value
*/
export async function verifyCertificates(certA, certB, leaf) {
  const parsedCertA = jsonToCertObject(certA);
  const parsedCertB = jsonToCertObject(certB);
  if (leaf) {
    const parsedLeaf = jsonToCertObject(leaf);
    const chainA = await parsedLeaf.verify(parsedCertA);
    const chainB = await parsedLeaf.verify(parsedCertB);
    // the leaf's issuer should be equal the subject data of the intermediate certs
    const isEqualA = parsedLeaf.issuer.isEqual(parsedCertA.subject);
    const isEqualB = parsedLeaf.issuer.isEqual(parsedCertB.subject);
    return chainA && chainB && isEqualA && isEqualB;
  }
  // can be used to validate if a certificate is self-signed, by passing it as both certA and B (i.e. a root cert)
  return (await parsedCertA.verify(parsedCertB)) && parsedCertA.issuer.isEqual(parsedCertB.subject);
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

  const isUnexpectedSubjectOid = (rdn) => !Object.values(SUBJECT_OIDs).includes(rdn.type);
  if (subject.any(isUnexpectedSubjectOid)) {
    const unknown = subject.filter(isUnexpectedSubjectOid).map((rdn) => rdn.type);
    errors.push(new Error('certificate contains unsupported subject OIDs: ' + unknown.join(', ')));
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
  const isUnknownExtension = (ext) => !allowedOids.includes(ext.extnID);
  if (extensions.any(isUnknownExtension)) {
    const unknown = extensions.filter(isUnknownExtension).map((ext) => ext.extnID);
    errors.push(new Error('certificate contains unsupported extension OIDs: ' + unknown.join(', ')));
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
      errors.push(new Error('certificate contains an unsupported subjectAltName construction'));
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
    const parsed_ips = [];
    for (const ip_san of values.ip_sans) {
      const unused = ip_san.valueBlock.unusedBits;
      if (unused !== undefined && unused !== null && unused !== 0) {
        errors.push(new Error('unsupported ip_san value: non-zero unused bits in encoding'));
        continue;
      }

      const ip = new Uint8Array(ip_san.valueBlock.valueHex);

      // Length of the IP determines the type: 4 bytes for IPv4, 16 bytes for
      // IPv6.
      if (ip.length === 4) {
        const ip_addr = ip.join('.');
        parsed_ips.push(ip_addr);
      } else if (ip.length === 16) {
        const src = new Array(...ip);
        const hex = src.map((value) => '0' + new Number(value).toString(16));
        const trimmed = hex.map((value) => value.slice(value.length - 2, 3));
        // add a colon after every other number (those with an odd index)
        let ip_addr = trimmed.map((value, index) => (index % 2 === 0 ? value : value + ':')).join('');
        // Remove trailing :, if any.
        ip_addr = ip_addr.slice(-1) === ':' ? ip_addr.slice(0, -1) : ip_addr;
        parsed_ips.push(ip_addr);
      } else {
        errors.push(
          new Error(
            'unsupported ip_san value: unknown IP address size (should be 4 or 16 bytes, was ' +
              parseInt(ip.length / 2) +
              ')'
          )
        );
      }
    }
    values.ip_sans = parsed_ips;
  }

  if (values.key_usage) {
    // KeyUsage is a big-endian bit-packed enum. Unused right-most bits are
    // truncated. So, a KeyUsage with CertSign+CRLSign would be "000001100",
    // with the right two bits truncated, and packed into an 8-bit, one-byte
    // string ("00000011"), introducing a leading zero. unused indicates that
    // this bit can be discard, shifting our result over by one, to go back
    // to its original form (minus trailing zeros).
    //
    // We can thus take our enumeration (KEY_USAGE_BITS), check whether the
    // bits are asserted, and push in our pretty names as appropriate.
    const unused = values.key_usage.valueBlock.unusedBits;
    const keyUsage = new Uint8Array(values.key_usage.valueBlock.valueHex);

    const computedKeyUsages = [];
    for (const enumIndex in KEY_USAGE_BITS) {
      // May span two bytes.
      const byteIndex = parseInt(enumIndex / 8);
      const bitIndex = parseInt(enumIndex % 8);
      const enumName = KEY_USAGE_BITS[enumIndex];
      const mask = 1 << (8 - bitIndex); // Big endian.
      if (byteIndex >= keyUsage.length) {
        // DecipherOnly is rare and would push into a second byte, but we
        // don't have one so exit.
        break;
      }

      let enumByte = keyUsage[byteIndex];
      const needsAdjust = byteIndex + 1 === keyUsage.length && unused > 0;
      if (needsAdjust) {
        enumByte = parseInt(enumByte << unused);
      }

      const isSet = (mask & enumByte) === mask;
      if (isSet) {
        computedKeyUsages.push(enumName);
      }
    }

    // Vault currently doesn't allow setting key_usage during issuer
    // generation, but will allow it if it comes in via an externally
    // generated CSR. Validate that key_usage matches expectations and
    // prune accordingly.
    const expectedUsages = ['CertSign', 'CRLSign'];
    const isUnexpectedKeyUsage = (ext) => !expectedUsages.includes(ext);

    if (computedKeyUsages.any(isUnexpectedKeyUsage)) {
      const unknown = computedKeyUsages.filter(isUnexpectedKeyUsage);
      errors.push(new Error('unsupported key usage value on issuer certificate: ' + unknown.join(', ')));
    }

    values.key_usage = computedKeyUsages;
  }

  if (values.other_sans) {
    // We need to parse these into their server-side values.
    const parsed_sans = [];
    for (const san of values.other_sans) {
      let [objectId, constructed] = san.valueBlock.value;
      objectId = objectId.toJSON().valueBlock.value;
      constructed = constructed.valueBlock.value[0].toJSON(); // can I just grab the first element here?
      const { blockName } = constructed;
      const value = constructed.valueBlock.value;
      parsed_sans.push(`${objectId};${blockName.replace('String', '')}:${value}`);
    }
    values.other_sans = parsed_sans;
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
    "key_usage": ['CertSign', 'CRLSign'],
    "ip_sans": ['192.158.1.38', '1234:fd2:5621:1:89::4500'],
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
