/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { CredentialType, DestinationType } from 'sync/utils/constants';

import type { SystemWriteSyncDestinationsAzureKvNameRequest } from '@hashicorp/vault-client-typescript';
import CreateDestinationForm from './create-destination';

type AzureKvFormData = SystemWriteSyncDestinationsAzureKvNameRequest & {
  name: string;
  credential_type: CredentialType;
};

export default class AzureKvForm extends CreateDestinationForm<AzureKvFormData> {
  // the "clientSecret" param is not checked because it's never returned by the API.
  // thus we can never say for sure if the account accessType has been configured so we always return false
  isAccountPluginConfigured = false;

  get isWifPluginConfigured() {
    const { identity_token_audience, identity_token_ttl } = this.data;
    return !!identity_token_audience || !!identity_token_ttl;
  }

  accountCredentialGroup = new FormFieldGroup('Client secret', [
    new FormField('client_secret', 'string', {
      subText:
        'Client secret of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_SECRET environment variable if configured.',
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
        new FormField('key_vault_uri', 'string', {
          label: 'Key Vault URI',
          subText:
            'URI of an existing Azure Key Vault instance. If empty, Vault will use the KEY_VAULT_URI environment variable if configured.',
          editDisabled: true,
        }),
        new FormField('tenant_id', 'string', {
          label: 'Tenant ID',
          subText:
            'ID of the target Azure tenant. If empty, Vault will use the AZURE_TENANT_ID environment variable if configured.',
          editDisabled: true,
        }),
        new FormField('cloud', 'string', {
          subText: 'Specifies a cloud for the client. The default is Azure Public Cloud.',
          editDisabled: true,
        }),
        new FormField('client_id', 'string', {
          label: 'Client ID',
          subText:
            'Client ID of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_ID environment variable if configured.',
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
    const data = this.getPayload<AzureKvFormData>(DestinationType.AzureKv, this.data, this.isNew);
    return { ...formState, data };
  }
}
