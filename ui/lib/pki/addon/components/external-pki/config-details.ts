/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { PkiExternalCaReadConfigAcmeAccountResponse } from '@hashicorp/vault-client-typescript';
import { toLabel } from 'core/helpers/to-label';

interface Args {
  configs: PkiExternalCaReadConfigAcmeAccountResponse[];
}

export default class ExternalPkiConfigDetailsComponent extends Component<Args> {
  excludedFields = ['name', 'account_keys'];

  label = (field: string) => {
    const label = toLabel([field]);
    const transformedLabel = label.replace(/\b[Ii]d\b/g, 'ID');

    return (
      {
        assume_role_arn: 'IAM role ARN to assume',
        directory_url: 'Directory URL',
        key_type: 'Active key type',
        nameserver: 'DNS server address',
        trusted_ca: 'Trusted CA',
        tsig_algorithm: 'TSIG algorithm',
        tsig_key_name: 'TSIG key name',
        ttl: 'Time to live',
      }[field] || transformedLabel
    );
  };
}
