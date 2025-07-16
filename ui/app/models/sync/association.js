/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class SyncAssociationModel extends Model {
  @attr mount;
  @attr secretName;
  @attr syncStatus;
  @attr updatedAt;
  // destination related properties that are not serialized to payload
  @attr destinationName;
  @attr destinationType;
  @attr subKey; // this property is added if a destination has 'secret-key' granularity

  @lazyCapabilities(
    apiPath`sys/sync/destinations/${'destinationType'}/${'destinationName'}/associations/set`,
    'destinationType',
    'destinationName'
  )
  setAssociationPath;

  @lazyCapabilities(
    apiPath`sys/sync/destinations/${'destinationType'}/${'destinationName'}/associations/remove`,
    'destinationType',
    'destinationName'
  )
  removeAssociationPath;

  get canSync() {
    return this.setAssociationPath.get('canUpdate') !== false;
  }

  get canUnsync() {
    return this.removeAssociationPath.get('canUpdate') !== false;
  }
}
