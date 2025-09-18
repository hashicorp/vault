/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { CLIENT_TYPES, ClientFilters, EXPORT_CLIENT_TYPES } from 'core/utils/client-count-utils';

// At time of writing ClientTypes are: 'acme_clients' | 'clients' | 'entity_clients' | 'non_entity_clients' | 'secret_syncs'
export type ClientTypes = (typeof CLIENT_TYPES)[number];

// 'namespace_path' | 'mount_path' | 'mount_type' | 'month'
export type ClientFilterTypes = (typeof ClientFilters)[keyof typeof ClientFilters];

// client_type in the exported activity data differs slightly from the types of client keys
// returned by sys/internal/counters/activity endpoint (:
// 'non-entity-token' | 'pki-acme' | 'secret-sync' | 'entity'
type ActivityExportClientTypes = (typeof EXPORT_CLIENT_TYPES)[number];

export interface TotalClients {
  clients: number;
  entity_clients: number;
  non_entity_clients: number;
  secret_syncs: number;
  acme_clients: number;
}

// extend this type when the counts are optional (eg for new clients)
interface TotalClientsSometimes {
  clients?: number;
  entity_clients?: number;
  non_entity_clients?: number;
  secret_syncs?: number;
  acme_clients?: number;
}

export interface ByNamespaceClients extends TotalClients {
  label: string;
  mounts: MountClients[];
}

export interface MountClients extends TotalClients {
  label: string;
  mount_path: string;
  mount_type: string;
  namespace_path: string;
}

export interface ByMonthClients extends TotalClients {
  timestamp: string;
  namespaces: ByNamespaceClients[];
  new_clients: ByMonthNewClients;
}

export interface ByMonthNewClients extends TotalClientsSometimes {
  timestamp: string;
  namespaces: ByNamespaceClients[];
}

export interface NamespaceNewClients extends TotalClientsSometimes {
  timestamp: string;
  label: string;
  mounts: MountClients[];
}

export interface MountNewClients extends TotalClientsSometimes {
  timestamp: string;
  label: string;
}

// SERIALIZED RESPONSE DATA from activity/export API
export interface ActivityExportData {
  client_id: string;
  client_type: ActivityExportClientTypes;
  namespace_id: string;
  namespace_path: string;
  mount_accessor: string;
  mount_type: string;
  mount_path: string;
  token_creation_time: string;
  client_first_used_time: string;
}

// API RESPONSE SHAPE (prior to serialization)
export interface NamespaceObject {
  namespace_id: string;
  namespace_path: string;
  counts: Counts;
  mounts: { mount_path: string; counts: Counts; mount_type: string }[];
}

export type ActivityMonthStandard = {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: Counts;
    namespaces: NamespaceObject[];
    timestamp: string;
  };
};

export type ActivityMonthNoNewClients = {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: null;
    namespaces: null;
  };
};

export type ActivityMonthEmpty = {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: null;
  namespaces: null;
  new_clients: null;
};

export type ActivityMonthBlock = ActivityMonthEmpty | ActivityMonthNoNewClients | ActivityMonthStandard;

export interface Counts {
  acme_clients: number;
  clients: number;
  entity_clients: number;
  non_entity_clients: number;
  secret_syncs: number;
}
