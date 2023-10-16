/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';

// Base model for all secret sync destination types
@withFormFields()
export default class SyncDestinationsBaseModel extends Model {
  @attr('string') name;
}
