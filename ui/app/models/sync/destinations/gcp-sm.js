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
  credentials: [{ type: 'presence', message: 'Credentials are required.' }],
};
const fields = ['name', 'credentials'];

@withModelValidations(validations)
@withFormFields(fields)
export default class SyncDestinationsGoogleCloudSecretManagerModel extends SyncDestinationModel {
  @attr('string', {
    subText: 'JSON credentials for GCP secret manager.',
    editType: 'file',
    docLink: '/vault/docs/secrets/gcp#authentication',
  })
  credentials;
}
