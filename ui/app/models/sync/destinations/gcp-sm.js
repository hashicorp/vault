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
const displayFields = ['name', 'credentials'];
const formFieldGroups = [{ default: ['name'] }, { Credentials: ['credentials'] }];
@withModelValidations(validations)
@withFormFields(displayFields, formFieldGroups)
export default class SyncDestinationsGoogleCloudSecretManagerModel extends SyncDestinationModel {
  @attr('string', {
    label: 'JSON credentials',
    subText:
      'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
    editType: 'file',
    docLink: '/vault/docs/secrets/gcp#authentication',
  })
  credentials; // obfuscated, never returned by API

  // TODO - confirm if project_id is going to be added to READ response (not editable)
}
