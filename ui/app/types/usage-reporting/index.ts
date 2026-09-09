/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface SimpleDatum {
  value: number;
  label: string;
}

/*
  States for disaster recovery and performance
  (More states might be added once this is hooked up to the backend)
*/
export enum REPLICATION_ENABLED_STATE {
  PRIMARY = 'primary',
  SECONDARY = 'secondary',
  BOOTSTRAPPING = 'bootstrapping',
}
export const REPLICATION_DISABLED_STATE = 'disabled';

export interface UsageDashboardData {
  authMethods: Record<string, number>;
  leasesByAuthMethod: Record<string, number>;
  kvv1Secrets: number;
  kvv2Secrets: number;
  leaseCountQuotas: {
    globalLeaseCountQuota: {
      capacity: number;
      count: number;
      name: string;
    };
    totalLeaseCountQuotas: number;
  };
  namespaces: number;
  secretSync: {
    totalDestinations: number;
    destinations: Record<string, number>;
  };
  pki: {
    totalIssuers: number;
    totalRoles: number;
  };
  replicationStatus: {
    drPrimary: boolean;
    drState: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
    prPrimary: boolean;
    prState: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
  };
  secretEngines: Record<string, number>;
}

export interface NamespaceData {
  keys: string[];
}

/* Fetches data for the dashboard (For a specific namespace if provided) */
export type getUsageDataFunction = (namespace?: string) => Promise<UsageDashboardData>;

/* Fetches a list of namespaces for the namespace picker */
export type getNamespaceDataFunction = () => Promise<NamespaceData>;
