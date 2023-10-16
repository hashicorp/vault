/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { syncDestinations } from 'vault/helpers/sync-destinations';

// Base model for all secret sync destination types
export default class SyncDestinationsBaseModel extends Model {
  @attr('string') name;
  @attr type;

  get icon() {
    return syncDestinations().findBy('type', this.type).icon;
  }
}
