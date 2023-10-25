/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { findDestination } from 'vault/helpers/sync-destinations';

// Base model for all secret sync destination types
export default class SyncDestinationModel extends Model {
  @attr('string', { subText: 'Specifies the name for this destination.' }) name;
  @attr type;

  // findDestination returns static attributes for each destination type
  get icon() {
    return findDestination(this.type)?.icon;
  }

  get typeDisplayName() {
    return findDestination(this.type)?.name;
  }

  get maskedParams() {
    return findDestination(this.type)?.maskedParams;
  }
}
