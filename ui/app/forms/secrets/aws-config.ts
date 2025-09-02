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
          const { lease, lease_max } = data;
          return (lease && lease_max) || (!lease && !lease_max) ? true : false;
        },
        message: 'Lease TTL and Max Lease TTL are both required if one of them is set.',
      },
    ],
  };

  get isAccountPluginConfigured() {
    return !!this.data.access_key;
  }

  get isWifPluginConfigured() {
    const { identity_token_audience, identity_token_ttl, role_arn } = this.data;
    return !!identity_token_audience || !!identity_token_ttl || !!role_arn;
  }

  accountFields = [
    new FormField('access_key', 'string'),
    new FormField('secret_key', 'string', { sensitive: true }),
  ];

  optionFields = [
    new FormField('region', 'string', {
      possibleValues: regions(),
      subText:
        'Specifies the AWS region. If not set it will use the AWS_REGION env var, AWS_DEFAULT_REGION env var, or us-east-1 in that order.',
    }),
    new FormField('iam_endpoint', 'string', { label: 'IAM endpoint' }),
    new FormField('sts_endpoint', 'string', { label: 'STS endpoint' }),
    new FormField('max_retries', 'number', {
      subText: 'Number of max retries the client should use for recoverable errors. Default is -1.',
    }),
  ];

  wifFields = [
    this.commonWifFields.issuer,
    new FormField('role_arn', 'string', {
      label: 'Role ARN',
      subText: 'Role ARN to assume for plugin workload identity federation.',
    }),
    this.commonWifFields.identity_token_audience,
    this.commonWifFields.identity_token_ttl,
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
    new FormField('lease_max', 'string', {
      label: 'Max Lease TTL',
      editType: 'ttl',
    }),
  ];
}
