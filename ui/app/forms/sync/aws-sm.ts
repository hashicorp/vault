/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { commonFields, getPayload } from './shared';

import type { SystemWriteSyncDestinationsAwsSmNameRequest } from '@hashicorp/vault-client-typescript';

type AwsSmFormData = SystemWriteSyncDestinationsAwsSmNameRequest & {
  name: string;
};

export default class AwsSmForm extends Form<AwsSmFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      commonFields.name,
      new FormField('region', 'string', {
        subText:
          'For AWS secrets manager, the name of the region must be supplied, something like “us-west-1.” If empty, Vault will use the AWS_REGION environment variable if configured.',
        editDisabled: true,
      }),
      new FormField('roleArn', 'string', {
        label: 'Role ARN',
        subText:
          'Specifies a role to assume when connecting to AWS. When assuming a role, Vault uses temporary STS credentials to authenticate.',
      }),
      new FormField('externalId', 'string', {
        label: 'External ID',
        subText:
          'Optional extra protection that must match the trust policy granting access to the AWS IAM role ARN. We recommend using a different random UUID per destination.',
      }),
    ]),
    new FormFieldGroup('Credentials', [
      new FormField('accessKeyId', 'string', {
        label: 'Access key ID',
        subText:
          'Access key ID to authenticate against the secrets manager. If empty, Vault will use the AWS_ACCESS_KEY_ID environment variable if configured.',
        sensitive: true,
        noCopy: true,
      }),
      new FormField('secretAccessKey', 'string', {
        label: 'Secret access key',
        subText:
          'Secret access key to authenticate against the secrets manager. If empty, Vault will use the AWS_SECRET_ACCESS_KEY environment variable if configured.',
        sensitive: true,
        noCopy: true,
      }),
    ]),
    new FormFieldGroup('Advanced configuration', [
      commonFields.granularity,
      commonFields.secretNameTemplate,
      commonFields.customTags,
    ]),
  ];

  toJSON() {
    const formState = super.toJSON();
    const data = getPayload<AwsSmFormData>('aws-sm', this.data, this.isNew);
    return { ...formState, data };
  }
}
