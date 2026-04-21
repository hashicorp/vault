/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { normalizeMetricData } from 'vault/utils/metrics-helpers';
import { module, test } from 'qunit';

module('Unit | Utility | metric utils', function () {
  test('normalizeMetricData returns undefined for null or undefined input', function (assert) {
    assert.strictEqual(normalizeMetricData(null), undefined, 'Returns undefined for null input');
    assert.strictEqual(normalizeMetricData(undefined), undefined, 'Returns undefined for undefined input');
  });

  test('normalizeMetricData returns all zeros for missing usage_metrics', function (assert) {
    const metric = {
      month: '2026-03',
      updated_at: '2026-04-01T06:59:59Z',
      usage_metrics: [
        {
          metric_data: {
            metric_details: [],
            total: 0,
          },
          metric_name: 'static_secrets',
        },
        {
          metric_data: {
            metric_details: [],
            total: 0,
          },
          metric_name: 'dynamic_roles',
        },
        {
          metric_data: {
            metric_details: [],
            total: 0,
          },
          metric_name: 'auto_rotated_roles',
        },
        {
          metric_data: {
            used_in_month: false,
          },
          metric_name: 'kmip',
        },
        {
          metric_data: {
            total: 0,
          },
          metric_name: 'external_plugins',
        },
        {
          metric_data: {
            metric_details: [],
            total: 0,
          },
          metric_name: 'data_protection_calls',
        },
        {
          metric_data: {
            total: 0,
          },
          metric_name: 'pki_units',
        },
        {
          metric_data: {
            metric_details: [],
            total: 0,
          },
          metric_name: 'managed_keys',
        },
      ],
    };
    const expected = {
      auto_rotated_roles_total: 0,
      data_protection_calls_total: 0,
      data_protection_calls_transform: 0,
      data_protection_calls_transit: 0,
      dynamic_roles_total: 0,
      external_plugins_total: 0,
      kmip_used_in_month: false,
      managed_keys: 0,
      managed_keys_kmse: 0,
      managed_keys_total: 0,
      managed_keys_totp: 0,
      pki_units_total: 0,
      ssh_units: 0,
      ssh_units_certificate_units: 0,
      ssh_units_otp_units: 0,
      ssh_units_total: 0,
      static_secrets_kv: 0,
      static_secrets_total: 0,
    };
    assert.deepEqual(normalizeMetricData(metric), expected, 'Returns all zeros for missing usage_metrics');
  });

  test('normalizeMetricData handles metric_details with missing type or count', function (assert) {
    const metric = {
      usage_metrics: [
        {
          metric_name: 'static_secrets',
          metric_data: {
            total: 5,
            metric_details: [
              { type: 'kv', count: 5 },
              { type: null, count: 3 },
              { type: 'foo' },
              { count: 2 },
            ],
          },
        },
      ],
    };
    const result = normalizeMetricData(metric);
    assert.strictEqual(result.static_secrets_kv, 5, 'Only valid detail is included');
    assert.strictEqual(result.static_secrets_foo, undefined, 'Detail with missing count is ignored');
  });
});
