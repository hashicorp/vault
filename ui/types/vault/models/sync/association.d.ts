/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Model } from 'vault/app-types';

export default interface SyncAssocationModel extends Model {
  mount: string;
  secretName: string;
  syncStatus: string;
  updatedAt: string;
  destinationName: string;
  destinationType: string;

  get canSync(): boolean;
  get canUnsync(): boolean;
}
