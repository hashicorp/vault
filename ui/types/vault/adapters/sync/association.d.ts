/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

import type { SyncDestinationQueryData } from './destination';

export interface SyncStatus {
  destinationType: string;
  destinationName: string;
  syncStatus: string;
  updatedAt: string;
}

export interface SyncDestinationAssociationMetrics {
  icon: string;
  name: string;
  associationCount: number;
  status: string;
  lastUpdated: Date;
}

export interface SyncAssociationMetrics {
  total_associations: number;
  total_secrets: number;
}

export default interface LdapLibraryAdapter extends AdapterRegistry {
  queryAll(): Promise<SyncAssociationMetrics>;
  fetchSyncStatus(mount: string, secretName: string): SyncStatus[];
  fetchByDestinations(destinations: SyncDestinationQueryData[]): Promise<SyncDestinationAssociationMetrics[]>;
}
