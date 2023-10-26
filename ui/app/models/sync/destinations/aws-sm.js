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
};
const fields = ['name', 'region', 'accessKeyId', 'secretAccessKey'];
@withModelValidations(validations)
@withFormFields(fields)
export default class SyncDestinationsAwsSecretsManagerModel extends SyncDestinationModel {
  @attr('string', {
    label: 'Access key ID',
    subText: 'Access key ID to authenticate against the secrets manager.',
  })
  accessKeyId;

  @attr('string', {
    label: 'Secret access key',
    subText: 'Secret access key to authenticate against the secrets manager.',
  })
  secretAccessKey;

  @attr('string', {
    subText: 'For AWS secrets manager, the name of the region must be supplied, something like “us-west-1.”',
  })
  region;
}
