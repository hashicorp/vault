/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

const displayFields = [
  // connection details
  'name',
  'keyVaultUri',
  'tenantId',
  'cloud',
  'clientId',
  'clientSecret',
  // vault sync config options
  'granularity',
  'secretNameTemplate',
  'customTags',
];
const formFieldGroups = [
  {
    default: [
      'name',
      'keyVaultUri',
      'tenantId',
      'cloud',
      'clientId',
      'granularity',
      'secretNameTemplate',
      'customTags',
    ],
  },
  { Credentials: ['clientSecret'] },
];
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsAzureKeyVaultModel extends SyncDestinationModel {
  @attr('string', {
    label: 'Key Vault URI',
    subText:
      'URI of an existing Azure Key Vault instance. If empty, Vault will use the KEY_VAULT_URI environment variable if configured.',
    editDisabled: true,
  })
  keyVaultUri;

  @attr('string', {
    label: 'Client ID',
    subText:
      'Client ID of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_ID environment variable if configured.',
  })
  clientId;

  @attr('string', {
    subText:
      'Client secret of an Azure app registration. If empty, Vault will use the AZURE_CLIENT_SECRET environment variable if configured.',
  })
  clientSecret; // obfuscated, never returned by API

  @attr('string', {
    label: 'Tenant ID',
    subText:
      'ID of the target Azure tenant. If empty, Vault will use the AZURE_TENANT_ID environment variable if configured.',
    editDisabled: true,
  })
  tenantId;

  @attr('string', {
    subText: 'Specifies a cloud for the client. The default is Azure Public Cloud.',
    editDisabled: true,
  })
  cloud;

  @attr('object', {
    subText:
      'An optional set of informational key-value pairs added as additional metadata on secrets synced to this destination. Custom tags are merged with built-in tags.',
    editType: 'kv',
  })
  customTags;
}
