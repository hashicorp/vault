/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import SyncDestinationsBaseModel from './base';
import { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class SyncDestinationsGoogleCloudSecretManagerModel extends SyncDestinationsBaseModel {
  @attr('string') credentials;

  get type() {
    return 'gcp-sm';
  }

  get icon() {
    return 'gcp-color';
  }
}
