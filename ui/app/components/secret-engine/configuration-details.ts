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
    'roleArn',
    'identityTokenAudience',
    'identityTokenTtl',
    'accessKey',
    'region',
    'iamEndpoint',
    'stsEndpoint',
    'maxRetries',
    'lease',
    'leaseMax',
    'issuer',
  ];

  azureFields = [
    'subscriptionId',
    'tenantId',
    'clientId',
    'identityTokenAudience',
    'identityTokenTtl',
    'rootPasswordTtl',
    'environment',
    'issuer',
  ];

  gcpFields = ['serviceAccountEmail', 'ttl', 'maxTtl', 'identityTokenAudience', 'identityTokenTtl', 'issuer'];

  sshFields = ['publicKey', 'generateSigningKey'];

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
        leaseMax: 'Max Lease TTL',
        ttl: 'Config TTL',
      }[field] || formattedLabel
    );
  };

  isDuration = (field: string) => {
    return ['identityTokenTtl', 'rootPasswordTtl', 'lease', 'leaseMax', 'ttl', 'maxTtl'].includes(field);
  };
}
