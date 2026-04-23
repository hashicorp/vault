/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export const METRICS_DATA_RESPONSE = {
  request_id: 'c16f7715-e534-7ebe-439f-7bc6005a26da',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    months: [
      {
        month: '2026-01',
        updated_at: '2026-01-14T10:49:00Z',
        usage_metrics: [
          {
            metric_name: 'static_secrets',
            metric_data: {
              total: 10,
              metric_details: [{ type: 'kv', count: 10 }],
            },
          },
          {
            metric_name: 'dynamic_roles',
            metric_data: {
              total: 130,
              metric_details: [
                { type: 'aws_dynamic', count: 10 },
                { type: 'azure_dynamic', count: 10 },
                { type: 'database_dynamic', count: 10 },
                { type: 'gcp_dynamic', count: 10 },
                { type: 'ldap_dynamic', count: 10 },
                { type: 'openldap_dynamic', count: 10 },
                { type: 'alicloud_dynamic', count: 10 },
                { type: 'rabbitmq_dynamic', count: 10 },
                { type: 'consul_dynamic', count: 10 },
                { type: 'nomad_dynamic', count: 10 },
                { type: 'kubernetes_dynamic', count: 10 },
                { type: 'mongodbatlas_dynamic', count: 10 },
                { type: 'terraform_dynamic', count: 10 },
              ],
            },
          },
          {
            metric_name: 'auto_rotated_roles',
            metric_data: {
              total: 70,
              metric_details: [
                { type: 'aws_static', count: 10 },
                { type: 'azure_static', count: 10 },
                { type: 'database_static', count: 10 },
                { type: 'gcp_static', count: 10 },
                { type: 'gcp_impersonated', count: 10 },
                { type: 'ldap_static', count: 10 },
                { type: 'openldap_static', count: 10 },
              ],
            },
          },
          {
            metric_name: 'kmip',
            metric_data: {
              used_in_month: true,
            },
          },
          {
            metric_name: 'pki_units',
            metric_data: { total: 100.1234 },
          },
          {
            metric_name: 'ssh_units',
            metric_data: {
              total: 100.2468,
              metric_details: [
                { type: 'otp_units', count: 50.1234 },
                { type: 'certificate_units', count: 50.1234 },
              ],
            },
          },
          {
            metric_name: 'external_plugins',
            metric_data: { total: 100 },
          },
          {
            metric_name: 'data_protection_calls',
            metric_data: {
              total: 420,
              metric_details: [
                { type: 'transit', count: 200 },
                { type: 'transform', count: 220 },
              ],
            },
          },
          {
            metric_name: 'managed_keys',
            metric_data: {
              total: 430,
              metric_details: [
                { type: 'kmse', count: 210 },
                { type: 'totp', count: 220 },
              ],
            },
          },
        ],
      },
      {
        month: '2025-12',
        updated_at: '2026-01-14T10:49:00Z',
        usage_metrics: [
          {
            metric_name: 'static_secrets',
            metric_data: {
              total: 2,
              metric_details: [{ type: 'kv', count: 2 }],
            },
          },
          {
            metric_name: 'dynamic_roles',
            metric_data: {
              total: 125,
              metric_details: [
                { type: 'aws_dynamic', count: 5 },
                { type: 'azure_dynamic', count: 10 },
                { type: 'database_dynamic', count: 10 },
                { type: 'gcp_dynamic', count: 10 },
                { type: 'ldap_dynamic', count: 10 },
                { type: 'openldap_dynamic', count: 10 },
                { type: 'alicloud_dynamic', count: 10 },
                { type: 'rabbitmq_dynamic', count: 10 },
                { type: 'consul_dynamic', count: 10 },
                { type: 'nomad_dynamic', count: 10 },
                { type: 'kubernetes_dynamic', count: 10 },
                { type: 'mongodbatlas_dynamic', count: 10 },
                { type: 'terraform_dynamic', count: 10 },
              ],
            },
          },
          {
            metric_name: 'auto_rotated_roles',
            metric_data: {
              total: 65,
              metric_details: [
                { type: 'aws_static', count: 5 },
                { type: 'azure_static', count: 10 },
                { type: 'database_static', count: 10 },
                { type: 'gcp_static', count: 10 },
                { type: 'gcp_impersonated', count: 10 },
                { type: 'ldap_static', count: 10 },
                { type: 'openldap_static', count: 10 },
              ],
            },
          },
          {
            metric_name: 'kmip',
            metric_data: {
              used_in_month: false,
            },
          },
          {
            metric_name: 'pki_units',
            metric_data: { total: 100.1234 },
          },
          {
            metric_name: 'ssh_units',
            metric_data: {
              total: 100.2468,
              metric_details: [
                { type: 'otp_units', count: 50.1234 },
                { type: 'certificate_units', count: 50.1234 },
              ],
            },
          },
          {
            metric_name: 'external_plugins',
            metric_data: { total: 100 },
          },
          {
            metric_name: 'data_protection_calls',
            metric_data: {
              total: 220,
              metric_details: [
                { type: 'transit', count: 200 },
                { type: 'transform', count: 220 },
              ],
            },
          },
          {
            metric_name: 'managed_keys',
            metric_data: {
              total: 220,
              metric_details: [
                { type: 'kmse', count: 200 },
                { type: 'totp', count: 220 },
              ],
            },
          },
        ],
      },
    ],
  },
};
