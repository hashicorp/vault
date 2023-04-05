import Component from '@glimmer/component';
import { parsedParameterKeys } from 'vault/utils/parse-pki-cert-oids';

/**
 * Expects to have parsedCertificate from parse-pki-cert util passed in, and will render attributes if values exist
 */
export default class ParsedCertificateInfoRowsComponent extends Component {
  get possibleFields() {
    // We show common name elsewhere on the details view, so no need to render it here
    const fieldKeys = parsedParameterKeys.filter((k) => k !== 'common_name');
    const attrsByKey = {
      other_sans: { label: 'Other SANs' },
      alt_names: { label: 'Subject Alternative Names (SANs)' },
      uri_sans: { label: 'URI SANs' },
      ip_sans: { label: 'IP SANs' },
      permitted_dns_domains: { label: 'Permitted DNS domains' },
      exclude_cn_from_sans: { label: 'Exclude CN from SANs' },
      use_pss: { label: 'Use PSS' },
      ttl: { label: 'TTL' },
      ou: { label: 'Organizational units (OU)' },
      not_valid_after: { formatDate: 'MMM d yyyy HH:mm:ss a zzzz' },
      not_valid_before: { formatDate: 'MMM d yyyy HH:mm:ss a zzzz' },
    };

    return fieldKeys.map((key) => {
      return {
        key,
        ...attrsByKey[key],
      };
    });
  }
}
