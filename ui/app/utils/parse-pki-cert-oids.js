//* lookup OIDs: http://oid-info.com/basic-search.htm

export const SUBJECT_OIDs = {
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

export const EXTENSION_OIDs = {
  key_usage: '2.5.29.15', // contains keyUsage values
  subject_alt_name: '2.5.29.17', // contains SAN_TYPES below
  basic_constraints: '2.5.29.19', // contains max_path_length
  name_constraints: '2.5.29.30', // contains permitted_dns_domains
  ext_key_usage: '2.5.29.37', // contains EXT_KEY_USAGE_OIDs OIDs below
};

// these are allowed ext oids, but not parsed and passed to cross-signed certs
export const IGNORED_OIDs = {
  // These two extensions are controlled by the parent authority.
  authority_key_identifier: '2.5.29.35',
  authority_access_info: '1.3.6.1.5.5.7.1.1',
  // This extension is based off the key material of the new issuer, which
  // will automatically match the existing issuer's key material.
  subject_key_identifier: '2.5.29.14',
};

// SubjectAltName/GeneralName types (scroll up to page 38 -> https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.7 )
export const SAN_TYPES = {
  alt_names: 2, // dNSName
  uri_sans: 6, // uniformResourceIdentifier
  ip_sans: 7, // iPAddress - OCTET STRING
};

export const SIGNATURE_ALGORITHM_OIDs = {
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

export const KEY_USAGE_BITS = [
  'DigitalSignature',
  'ContentCommitment',
  'KeyEncipherment',
  'DataEncipherment',
  'KeyAgreement',
  'CertSign',
  'CRLSign',
  'EncipherOnly',
  'DecipherOnly',
];

export const EXT_KEY_USAGE_OIDs = {
  '2.5.29.37.0': 'Any',
  '1.3.6.1.5.5.7.3.1': 'ServerAuth',
  '1.3.6.1.5.5.7.3.2': 'ClientAuth',
  '1.3.6.1.5.5.7.3.3': 'CodeSigning',
  '1.3.6.1.5.5.7.3.4': 'EmailProtection',
  '1.3.6.1.5.5.7.3.5': 'IPSECEndSystem',
  '1.3.6.1.5.5.7.3.6': 'IPSECTunnel',
  '1.3.6.1.5.5.7.3.7': 'IPSECUser',
  '1.3.6.1.5.5.7.3.8': 'TimeStamping',
  '1.3.6.1.5.5.7.3.9': 'OCSPSigning',
};
