/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationsBaseModel from './base';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class SyncDestinationsAwsSecretsManagerModel extends SyncDestinationsBaseModel {
  @attr('string') accessKeyId;
  @attr('string') secretAccessKey;
  @attr('string') region;

  get type() {
    return 'aws-sm';
  }

  get icon() {
    return 'aws-color';
  }
}
