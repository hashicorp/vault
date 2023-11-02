/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { findDestination } from 'vault/helpers/sync-destinations';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

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

  @lazyCapabilities(apiPath`sys/sync/destinations/${'type'}/${'name'}`, 'type', 'name') destinationPath;
  @lazyCapabilities(apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/set`, 'type', 'name')
  setAssociationPath;

  get canCreate() {
    return this.destinationPath.get('canCreate') !== false;
  }
  get canDelete() {
    return this.destinationPath.get('canDelete') !== false;
  }
  get canEdit() {
    return this.destinationPath.get('canUpdate') !== false;
  }
  get canRead() {
    return this.destinationPath.get('canRead') !== false;
  }
  get canList() {
    return this.destinationPath.get('canList') !== false;
  }
  get canSync() {
    return this.setAssociationPath.get('canUpdate') !== false;
  }
}
