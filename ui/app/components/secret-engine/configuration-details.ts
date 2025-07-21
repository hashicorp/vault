/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toLabel } from 'core/helpers/to-label';

import type { AwsConfig, AzureConfig, GcpConfig, SshConfig } from 'vault/vault/secrets/engine';

type Args = {
  config: AwsConfig | AzureConfig | GcpConfig | SshConfig;
  typeDisplay: string;
};

export default class ConfigurationDetails extends Component<Args> {
  awsFields = [
    'role_arn',
    'identity_token_audience',
    'identity_token_ttl',
    'access_key',
    'region',
    'iam_endpoint',
    'sts_endpoint',
    'max_retries',
    'lease',
    'lease_max',
    'issuer',
  ];

  azureFields = [
    'subscription_id',
    'tenant_id',
    'client_id',
    'identity_token_audience',
    'identity_token_ttl',
    'root_password_ttl',
    'environment',
    'issuer',
  ];

  gcpFields = [
    'service_account_email',
    'ttl',
    'max_ttl',
    'identity_token_audience',
    'identity_token_ttl',
    'issuer',
  ];

  sshFields = ['public_key', 'generate_signing_key'];

  get displayFields() {
    switch (this.args.typeDisplay) {
      case 'AWS':
        return this.awsFields;
      case 'Azure':
        return this.azureFields;
      case 'Google Cloud':
        return this.gcpFields;
      case 'SSH':
        return this.sshFields;
      default:
        return [];
    }
  }

  label = (field: string) => {
    const label = toLabel([field]);
    // convert words like id and ttl to uppercase
    const formattedLabel = label
      .split(' ')
      .map((word: string) => {
        const acronyms = ['id', 'ttl', 'arn', 'iam', 'sts'];
        return acronyms.includes(word.toLowerCase()) ? word.toUpperCase() : word;
      })
      .join(' ');
    // map specific fields to custom labels
    return (
      {
        lease: 'Default Lease TTL',
        lease_max: 'Max Lease TTL',
        ttl: 'Config TTL',
      }[field] || formattedLabel
    );
  };

  isDuration = (field: string) => {
    return ['identity_token_ttl', 'root_password_ttl', 'lease', 'lease_max', 'ttl', 'max_ttl'].includes(
      field
    );
  };
}
