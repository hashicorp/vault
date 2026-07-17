/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export enum MetricNameEnum {
  STATIC_SECRETS = 'static_secrets',
  DATA_PROTECTION_CALLS = 'data_protection_calls',
  MANAGED_KEYS = 'managed_keys',
  KMIP = 'kmip',
  EXTERNAL_PLUGINS = 'external_plugins',
  DYNAMIC_ROLES = 'dynamic_roles',
  PKI_UNITS = 'pki_units',
  SSH_UNITS = 'ssh_units',
}

export interface Month {
  month: string;
  updated_at: string;
  usage_metrics: MetricData[];
}

export interface MetricData {
  metric_name: MetricNameEnum;
  metric_data: {
    metric_details: Array<{ type: string; count: number }>;
    used_in_month?: boolean;
    total: number;
  };
}

export interface NormalizedMetricsData {
  [key: string]: number | boolean | undefined;
}
