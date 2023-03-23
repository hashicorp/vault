/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';

export default class PkiCrlModel extends Model {
  // This model uses the backend value as the model ID
  get useOpenAPI() {
    return true;
  }

  @attr('string') expiry;
  @attr('boolean') autoRebuild;
  @attr('string') ocspExpiry;
  @attr('boolean') ocspDisable;
}
