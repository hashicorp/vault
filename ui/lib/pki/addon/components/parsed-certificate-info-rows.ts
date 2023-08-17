/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { parsedParameterKeys } from 'vault/utils/parse-pki-cert-oids';
import { ParsedCertificateData } from 'vault/vault/utils/parse-pki-cert';

/**
 * @module ParsedCertificateInfoRowsComponent
 * Renders attributes parsed from a PKI certificate (provided from parse-pki-cert util). It will only render rows for
 * defined values that match `parsedParameterKeys` imported from the helper. It never renders common_name, even though
 * the value is returned from the parse cert util, because `common_name` is important to PKI and we render it at the top.
 *
 * @example ```js
 * <ParsedCertificateInfoRows @model={{@model.parsedCertificate}} />
 * ```
 *
 * @param {object} model - object of parsed attributes from parse-pki-cert util
 */

interface AttrOptions {
  [key: string]: {
    label?: string;
    formatDate?: string;
  };
}

interface Args {
  model: ParsedCertificateData;
}

export default class ParsedCertificateInfoRowsComponent extends Component<Args> {
  get possibleFields() {
    // We show common name elsewhere on the details view, so no need to render it here
    const fieldKeys = parsedParameterKeys.filter((k) => k !== 'common_name');
    const attrsByKey: AttrOptions = {
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
