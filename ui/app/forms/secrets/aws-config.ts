/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import WifConfigForm from './wif-config';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { regions } from 'vault/helpers/aws-regions';

import type { Validations } from 'vault/app-types';
import type { AwsConfigFormData } from 'vault/secrets/engine';

export default class AwsConfigForm extends WifConfigForm<AwsConfigFormData> {
  validations: Validations = {
    lease: [
      {
        validator(data: AwsConfigForm['data']) {
          const { lease, leaseMax } = data;
          return (lease && leaseMax) || (!lease && !leaseMax) ? true : false;
        },
        message: 'Lease TTL and Max Lease TTL are both required if one of them is set.',
      },
    ],
  };

  get isAccountPluginConfigured() {
    return !!this.data.accessKey;
  }

  get isWifPluginConfigured() {
    const { identityTokenAudience, identityTokenTtl, roleArn } = this.data;
    return !!identityTokenAudience || !!identityTokenTtl || !!roleArn;
  }

  accountFields = [
    new FormField('accessKey', 'string'),
    new FormField('secretKey', 'string', { sensitive: true }),
  ];

  optionFields = [
    new FormField('region', 'string', {
      possibleValues: regions(),
      subText:
        'Specifies the AWS region. If not set it will use the AWS_REGION env var, AWS_DEFAULT_REGION env var, or us-east-1 in that order.',
    }),
    new FormField('iamEndpoint', 'string', { label: 'IAM endpoint' }),
    new FormField('stsEndpoint', 'string', { label: 'STS endpoint' }),
    new FormField('maxRetries', 'number', {
      subText: 'Number of max retries the client should use for recoverable errors. Default is -1.',
    }),
  ];

  wifFields = [
    this.commonWifFields.issuer,
    new FormField('roleArn', 'string', {
      label: 'Role ARN',
      subText: 'Role ARN to assume for plugin workload identity federation.',
    }),
    this.commonWifFields.identityTokenAudience,
    this.commonWifFields.identityTokenTtl,
  ];

  // formFieldGroups will render the default and root config option fields
  // formFields will be used to render the additional lease config fields
  // this allows for the full form to validated together when submitted
  get formFieldGroups() {
    const defaultFields = this.accessType === 'account' ? this.accountFields : this.wifFields;
    return [
      new FormFieldGroup('default', defaultFields),
      new FormFieldGroup('Root config options', this.optionFields),
    ];
  }

  formFields = [
    new FormField('lease', 'string', {
      label: 'Default Lease TTL',
      editType: 'ttl',
    }),
    new FormField('leaseMax', 'string', {
      label: 'Max Lease TTL',
      editType: 'ttl',
    }),
  ];
}
