/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { calculateSum } from 'vault/utils/chart-helpers';

import type { Month, NormalizedMetricsData } from 'vault/vault/billing/overview';

export enum NormalizedBillingMetrics {
  AUTO_ROTATED_ROLES_TOTAL = 'auto_rotated_roles_total',
  CREDENTIAL_UNITS_TOTAL = 'credential_units_total',
  DATA_PROTECTION_CALLS_TOTAL = 'data_protection_calls_total',
  DATA_PROTECTION_CALLS_TRANSFORM = 'data_protection_calls_transform',
  DATA_PROTECTION_CALLS_TRANSIT = 'data_protection_calls_transit',
  DATA_PROTECTION_CALLS_GCPKMS = 'data_protection_calls_gcpkms',
  DYNAMIC_ROLES_TOTAL = 'dynamic_roles_total',
  EXTERNAL_PLUGINS_TOTAL = 'external_plugins_total',
  ID_TOKEN_UNITS_TOTAL = 'id_token_units_total',
  ID_TOKEN_UNITS_OIDC = 'id_token_units_oidc',
  ID_TOKEN_UNITS_SPIFFE = 'id_token_units_spiffe',
  KMIP_USED_IN_MONTH = 'kmip_used_in_month',
  MANAGED_KEYS = 'managed_keys',
  MANAGED_KEYS_KMSE = 'managed_keys_kmse',
  MANAGED_KEYS_TOTAL = 'managed_keys_total',
  MANAGED_KEYS_TOTP = 'managed_keys_totp',
  PKI_UNITS_TOTAL = 'pki_units_total',
  SSH_UNITS = 'ssh_units',
  SSH_UNITS_CERTIFICATE_UNITS = 'ssh_units_certificate_units',
  SSH_UNITS_OTP_UNITS = 'ssh_units_otp_units',
  SSH_UNITS_TOTAL = 'ssh_units_total',
  STATIC_SECRETS_KV = 'static_secrets_kv',
  STATIC_SECRETS_TOTAL = 'static_secrets_total',
}

export enum BillingMetricsKeys {
  USED_IN_MONTH = 'used_in_month',
  KMIP = 'kmip',
  TOTAL = 'total',
}

export function normalizeMetricData(metric: Month | null | undefined) {
  const { usage_metrics } = metric || {};
  if (!usage_metrics) return;

  const normalized: NormalizedMetricsData = {};

  for (const metric of usage_metrics) {
    if (metric.metric_name === 'kmip') {
      const kmipKey = `${metric.metric_name}_used_in_month`;
      normalized[kmipKey] = metric.metric_data.used_in_month;
    }

    const metricName = metric.metric_name;
    const total = metric.metric_data?.total;

    if (typeof total === 'number') {
      normalized[`${metricName}_total`] = total;
    }

    for (const detail of metric.metric_data?.metric_details ?? []) {
      // Skip detail entries that are missing a type or a numeric count — both are required to build a valid normalized key.
      if (!detail.type || typeof detail.count !== 'number') continue;
      // Prefix parent metric_name to detail "type" to avoid future naming collisions.
      // For example the 'kv' type in the `metrics_details`
      // becomes `static_secrets_kv`:
      // {
      //   metric_name: 'static_secrets',
      //   metric_data: {
      //     total: 10,
      //     metric_details: [{ type: 'kv', count: 10 }],
      //   },
      // },
      const detailName = `${metricName}_${detail.type}`;
      normalized[detailName] = detail.count;
    }
  }

  // Calculate credential_units_total as the sum of ssh_units, pki_units, and id_token_units
  const sshUnitsTotal =
    typeof normalized[NormalizedBillingMetrics.SSH_UNITS_TOTAL] === 'number'
      ? normalized[NormalizedBillingMetrics.SSH_UNITS_TOTAL]
      : 0;
  const pkiUnitsTotal =
    typeof normalized[NormalizedBillingMetrics.PKI_UNITS_TOTAL] === 'number'
      ? normalized[NormalizedBillingMetrics.PKI_UNITS_TOTAL]
      : 0;
  const idTokenUnitsTotal =
    typeof normalized[NormalizedBillingMetrics.ID_TOKEN_UNITS_TOTAL] === 'number'
      ? normalized[NormalizedBillingMetrics.ID_TOKEN_UNITS_TOTAL]
      : 0;
  normalized[NormalizedBillingMetrics.CREDENTIAL_UNITS_TOTAL] =
    calculateSum([sshUnitsTotal, pkiUnitsTotal, idTokenUnitsTotal], 4) ?? 0;

  // Explicitly set any missing metric keys to 0.
  for (const metricsKey of Object.values(NormalizedBillingMetrics)) {
    if (!(metricsKey in normalized)) {
      normalized[metricsKey] = 0;
    }
  }

  return normalized;
}
