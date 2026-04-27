/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { regions } from 'vault/helpers/aws-regions';
import { CredentialType, DestinationType } from 'sync/utils/constants';
import CreateDestinationForm from './create-destination';

import type { SystemWriteSyncDestinationsAwsSmNameRequest } from '@hashicorp/vault-client-typescript';

type AwsSmFormData = SystemWriteSyncDestinationsAwsSmNameRequest & {
  name: string;
  credential_type: CredentialType;
};

export default class AwsSmForm extends CreateDestinationForm<AwsSmFormData> {
  get isAccountPluginConfigured() {
    return !!this.data.access_key_id;
  }

  get isWifPluginConfigured() {
    const { identity_token_audience, identity_token_ttl, role_arn } = this.data;
    return !!identity_token_audience || !!identity_token_ttl || !!role_arn;
  }

  accountCredentialGroup = new FormFieldGroup('IAM credentials', [
    new FormField('access_key_id', 'string', {
      label: 'Access key ID',
      subText:
        'Access key ID to authenticate against the secrets manager. If empty, Vault will use the AWS_ACCESS_KEY_ID environment variable if configured.',
      sensitive: true,
      noCopy: true,
    }),
    new FormField('secret_access_key', 'string', {
      label: 'Secret access key',
      subText:
        'Secret access key to authenticate against the secrets manager. If empty, Vault will use the AWS_SECRET_ACCESS_KEY environment variable if configured.',
      sensitive: true,
      noCopy: true,
    }),
  ]);

  get wifCredentialGroup() {
    return this.createWifCredentialGroup();
  }

  get formFieldGroups() {
    const credentialGroup =
      this.credentialType === CredentialType.ACCOUNT ? this.accountCredentialGroup : this.wifCredentialGroup;
    return [
      new FormFieldGroup('Destination details', [
        this.commonFields.name,
        new FormField('region', 'string', {
          possibleValues: regions(),
          noDefault: true,
          subText:
            'For AWS secrets manager, the name of the region must be supplied, something like “us-west-1.” If empty, Vault will use the AWS_REGION environment variable if configured.',
          editDisabled: true,
        }),
        new FormField('role_arn', 'string', {
          label: 'Role ARN',
          subText:
            'Specifies a role to assume when connecting to AWS. When assuming a role, Vault uses temporary STS credentials to authenticate.',
        }),
        new FormField('external_id', 'string', {
          label: 'External ID',
          subText:
            'Optional extra protection that must match the trust policy granting access to the AWS IAM role ARN. We recommend using a different random UUID per destination.',
        }),
      ]),
      credentialGroup,
      new FormFieldGroup('Advanced configuration', [
        this.commonFields.granularity,
        this.commonFields.secretNameTemplate,
        this.commonFields.customTags,
      ]),
    ];
  }

  toJSON() {
    const formState = super.toJSON();
    const data = this.getPayload<AwsSmFormData>(DestinationType.AwsSm, this.data, this.isNew);
    return { ...formState, data };
  }
}
