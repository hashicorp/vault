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

  get icon() {
    return findDestination(this.type)?.icon;
  }
}
