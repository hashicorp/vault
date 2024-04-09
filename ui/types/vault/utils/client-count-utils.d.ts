/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
The client count utils are responsible for serializing the sys/internal/counters/activity API response
The initial API response shape (line ~82) and serialized types are defined here and used to
defines the activity model in models/clients/activity.d.ts

To help visualize the response there are sample responses in ui/tests/helpers/clients.js
*/

// TYPES RETURNED BY UTILS (serialized)

export interface TotalClients {
  clients: number;
  entity_clients: number;
  non_entity_clients: number;
  secret_syncs: number;
  acme_clients: number;
}

export interface ByNamespaceClients extends TotalClients {
  label: string;
  mounts: MountClients[];
}

export interface MountClients extends TotalClients {
  label: string;
}

export interface ByMonthClients extends TotalClients {
  month: string;
  timestamp: string;
  namespaces: ByNamespaceClients[];
  namespaces_by_key: { [key: string]: NamespaceByKey };
  new_clients: ByMonthNewClients;
}

export interface ByMonthNewClients extends TotalClients {
  month: string;
  timestamp: string;
  namespaces: ByNamespaceClients[];
}

export interface NamespaceByKey extends TotalClients {
  month: string;
  timestamp: string;
  mounts_by_key: { [key: string]: MountByKey };
  new_clients: NamespaceNewClients;
}

export interface NamespaceNewClients extends TotalClients {
  month: string;
  label: string;
  mounts: MountClients[];
}

export interface MountByKey extends TotalClients {
  month: string;
  timestamp: string;
  label: string;
  new_clients: MountNewClients;
}

export interface MountNewClients extends TotalClients {
  month: string;
  label: string;
}

// remove?
export interface EmptyByMonthClients {
  month: string;
  timestamp: string;
  namespaces: [];
  namespaces_by_key: Record<string, never>;
  new_clients: {
    month: string;
    timestamp: string;
    namespaces: [];
  };
}

// API RESPONSE SHAPE (prior to serialization)

interface SysInternalCountersActivityResponse {
  start_time: string;
  end_time: string;
  total: Counts;
  by_namespace: NamespaceObject[];
  months: ActivityMonthBlock[];
}

export interface NamespaceObject {
  namespace_id: string;
  namespace_path: string;
  counts: Counts;
  mounts: { mount_path: string; counts: Counts }[];
}

export interface ActivityMonthBlock {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: Counts;
  namespaces: NamespaceObject[];
  new_clients: {
    counts: Counts;
    namespaces: NamespaceObject[];
    timestamp: string;
  };
}

export interface EmptyActivityMonthBlock {
  timestamp: string; // YYYY-MM-01T00:00:00Z (always the first day of the month)
  counts: null;
  namespaces: null;
  new_clients: null;
}

export interface Counts {
  acme_clients: number;
  clients: number;
  distinct_entities: number;
  entity_clients: number;
  non_entity_clients: number;
  non_entity_tokens: number;
  secret_syncs: number;
}
