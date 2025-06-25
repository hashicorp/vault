/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { commonFields, getPayload } from './shared';

import type { SystemWriteSyncDestinationsAzureKvNameRequest } from '@hashicorp/vault-client-typescript';

type AzureKvFormData = SystemWriteSyncDestinationsAzureKvNameRequest & {
  name: string;
};

export default class AzureKvForm extends Form<AzureKvFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      commonFields.name,
      new FormField('keyVaultUri', 'string', {
        label: 'Key Vault URI',
        subText:
          'URI of an existing Azure Key Vault instance. If empty, Vault will use the KEY_VAULT_URI environment variable if configured.',
        editDisabled: true,
      }),
      new FormField('tenantId', 'string', {
        label: 'Tenant ID',
        subText:
          'ID of the target Azure tenant. If empty, Vault will use the AZURE_TENANT_ID environment variable if configured.',
        editDisabled: true,
      }),
      new FormField('cloud', 'string', {
        subText: 'Specifies a cloud for the client. The default is Azure Public Cloud.',
        editDisabled: true,
      }),
      new FormField('clientId', 'string', {
        label: 'Client ID',
        subText:
          'Client ID of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_ID environment variable if configured.',
      }),
    ]),
    new FormFieldGroup('Credentials', [
      new FormField('clientSecret', 'string', {
        subText:
          'Client secret of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_SECRET environment variable if configured.',
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
    const data = getPayload<AzureKvFormData>('azure-kv', this.data, this.isNew);
    return { ...formState, data };
  }
}
