/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class PkiTidyModel extends Model {
  // This model uses the backend value as the model ID
  get useOpenAPI() {
    return true;
  }

  @attr('boolean', { defaultValue: false }) tidyCertStore;
  @attr('boolean', { defaultValue: false }) tidyRevocationQueue;
  @attr('string') safetyBuffer;
}
