/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationModel from '../destination';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class SyncDestinationsAwsSecretsManagerModel extends SyncDestinationModel {
  @attr('string') accessKeyId;
  @attr('string') secretAccessKey;
  @attr('string') region;
}
