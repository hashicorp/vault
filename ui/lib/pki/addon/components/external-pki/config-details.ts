/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toLabel } from 'core/helpers/to-label';

import type {
  PkiExternalCaReadConfigAcmeAccountResponse,
  PkiExternalCaReadConfigDnsAwsRoute53Response,
  PkiExternalCaReadConfigDnsAzureDnsResponse,
  PkiExternalCaReadConfigDnsGoogleCloudDnsResponse,
  PkiExternalCaReadConfigDnsRfc2136Response,
  PkiExternalCaReadLookupOrderResponse,
  PkiExternalCaReadRoleResponse,
} from '@hashicorp/vault-client-typescript';

interface Args {
  config:
    | PkiExternalCaReadConfigAcmeAccountResponse
    | PkiExternalCaReadLookupOrderResponse
    | PkiExternalCaReadConfigDnsAwsRoute53Response
    | PkiExternalCaReadConfigDnsAzureDnsResponse
    | PkiExternalCaReadConfigDnsGoogleCloudDnsResponse
    | PkiExternalCaReadConfigDnsRfc2136Response
    | PkiExternalCaReadRoleResponse;
}

export default class ExternalPkiConfigDetailsComponent extends Component<Args> {
  excludedFields = ['name', 'account_keys', 'challenges', 'identifiers', 'last_error'];

  get errorMessage() {
    if (this.args.config && 'last_error' in this.args.config) {
      return this.args.config.last_error;
    }
    return '';
  }

  label = (field: string) => {
    const label = toLabel([field]);
    const transformedLabel = label.replace(/\b[Ii]d\b/g, 'ID');

    return (
      {
        acme_account_name: 'ACME account name',
        assume_role_arn: 'IAM role ARN to assume',
        ca_chain: 'CA chain',
        csr_generate_key_type: 'CSR key generation type',
        csr_identifier_population: 'CSR identifier population',
        directory_url: 'Directory URL',
        dns_provider_name: 'DNS provider name',
        dns_provider_type: 'DNS provider type',
        key_type: 'Active key type',
        nameserver: 'DNS server address',
        not_after: 'Valid until',
        not_before: 'Valid after',
        trusted_ca: 'Trusted CA',
        tsig_algorithm: 'TSIG algorithm',
        tsig_key_name: 'TSIG key name',
        ttl: 'Time to live',
      }[field] || transformedLabel
    );
  };

  renderCopyButton = (field: string) => ['order_id', 'serial_number'].includes(field);

  renderEncodedDataCard = (field: string) =>
    ['ca_chain', 'certificate', 'private_key', 'trusted_ca'].includes(field);

  isDate = (field: string) =>
    ['creation_date', 'expires', 'last_update', 'next_work_date', 'not_after', 'not_before'].includes(field);
}
