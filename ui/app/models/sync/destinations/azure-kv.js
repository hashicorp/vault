/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
  keyVaultUri: [{ type: 'presence', message: 'Key Vault URI is required.' }],
  clientId: [{ type: 'presence', message: 'Client ID is required.' }],
  clientSecret: [{ type: 'presence', message: 'Client secret is required.' }],
  tenantId: [{ type: 'presence', message: 'Tenant ID is required.' }],
};
const fields = ['name', 'keyVaultUri', 'tenantId', 'cloud', 'clientId', 'clientSecret'];

@withModelValidations(validations)
@withFormFields(fields)
export default class SyncDestinationsAzureKeyVaultModel extends SyncDestinationModel {
  @attr('string', { label: 'Key Vault URI', subText: 'URI of an existing Azure Key Vault instance.' })
  keyVaultUri;

  @attr('string', { label: 'Client ID', subText: 'Client ID of an Azure app registration.' })
  clientId;

  @attr('string', { subText: 'Client secret of an Azure app registration.' })
  clientSecret;

  @attr('string', { label: 'Tenant ID', subText: 'ID of the target Azure tenant.' })
  tenantId;

  @attr('string', {
    subText: 'Specifies a cloud for the client. The default is Azure Public Cloud.',
    defaultValue: 'cloud',
  })
  cloud;
}
