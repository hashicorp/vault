import Component from '@glimmer/component';

/**
 * Expects to have parsedCertificate from parse-pki-cert util passed in, and will render attributes if values exist
 */
export default class ParsedCertificateInfoRowsComponent extends Component {
  get possibleFields() {
    return [
      { key: 'serial_number' },
      { key: 'key_usage' },
      { key: 'other_sans', label: 'Other SANs' },
      { key: 'alt_names', label: 'Subject Alternative Names (SANs)' },
      { key: 'uri_sans', label: 'URI SANs' },
      { key: 'ip_sans', label: 'IP SANs' },
      { key: 'permitted_dns_domains', label: 'Permitted DNS domains' },
      { key: 'max_path_length' },
      { key: 'exclude_cn_from_sans', label: 'Exclude CN from SANs' },
      { key: 'signature_bits' },
      { key: 'use_pss', label: 'Use PSS' },
      { key: 'not_valid_after', formatDate: 'MMM d yyyy HH:mm:ss a zzzz' },
      { key: 'not_valid_before', formatDate: 'MMM d yyyy HH:mm:ss a zzzz' },
      { key: 'ttl', label: 'TTL' },
      { key: 'ou', label: 'Organizational units (OU)' },
      { key: 'organization' },
      { key: 'country' },
      { key: 'locality' },
      { key: 'province' },
      { key: 'street_address' },
      { key: 'postal_code' },
    ];
  }
}
