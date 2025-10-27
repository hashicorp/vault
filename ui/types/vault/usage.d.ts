/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { GenerateUtilizationReportResponse } from '@hashicorp/vault-client-typescript';
import type {
  REPLICATION_ENABLED_STATE,
  REPLICATION_DISABLED_STATE,
} from '@hashicorp-internal/vault-reporting/types/index';

export type GlobalLeaseCountQuota = {
  capacity: number;
  count: number;
  name: string;
};

export type LeaseCountQuotas = {
  global_lease_count_quota: GlobalLeaseCountQuota;
  total_lease_count_quotas: number;
};

export type SecretSync = {
  total_destinations: number;
  destinations: Record<string, number>;
};

export type Pki = {
  total_issuers: number;
  total_roles: number;
};

export type ReplicationStatus = {
  dr_primary: boolean;
  dr_state: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
  pr_primary: boolean;
  pr_state: REPLICATION_ENABLED_STATE | typeof REPLICATION_DISABLED_STATE;
};

export type UtilizationReport = GenerateUtilizationReportResponse & {
  lease_count_quotas: LeaseCountQuotas;
  secret_sync: SecretSync;
  pki: Pki;
  replication_status: ReplicationStatus;
};
