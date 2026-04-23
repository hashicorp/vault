/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentURL } from '@ember/test-helpers';
import sinon from 'sinon';

import { login, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { NormalizedBillingMetrics } from 'vault/utils/metrics-helpers';

const mockMetrics = {
  months: [
    {
      month: '2026-02',
      updated_at: '2026-01-14T10:49:00Z', // Signal partial-month data
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
            total: 60,
            metric_details: [
              { type: 'aws_dynamic', count: 22 },
              { type: 'azure_dynamic', count: 20 },
              { type: 'database_dynamic', count: 30 },
            ],
          },
        },
        {
          metric_name: 'auto_rotated_roles',
          metric_data: {
            total: 30,
            metric_details: [
              { type: 'aws_static', count: 22 },
              { type: 'azure_static', count: 20 },
            ],
          },
        },
        { metric_name: 'kmip', metric_data: { used_in_month: true } },
        { metric_name: 'pki_units', metric_data: { total: 100.1234 } },
        { metric_name: 'data_protection_calls', metric_data: { total: 12 } },
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
          metric_name: 'managed_keys',
          metric_data: {
            total: 82,
            metric_details: [
              { type: 'totp', count: 52 },
              { type: 'kmse', count: 30 },
            ],
          },
        },
        /**
         * Other metrics:
         * - external_plugins
         * - managed_keys (types: totp, kmse)
         * - data_protection_calls (types: transit, transform)
         * - id_token_units (types: oidc, spiffe) // not adding in 2.0.0
         *
         * Additional metrics to be added for new features in Vault 1.22/2.0.
         */
      ],
    },
    {
      month: '2026-01',
      updated_at: '2026-01-14T10:49:00Z', // Signal partial-month data
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
            total: 60,
            metric_details: [
              { type: 'aws_dynamic', count: 10 },
              { type: 'azure_dynamic', count: 20 },
              { type: 'database_dynamic', count: 30 },
            ],
          },
        },
        {
          metric_name: 'auto_rotated_roles',
          metric_data: {
            total: 30,
            metric_details: [
              { type: 'aws_static', count: 10 },
              { type: 'azure_static', count: 20 },
            ],
          },
        },
        { metric_name: 'kmip', metric_data: { used_in_month: true } },
        { metric_name: 'pki_units', metric_data: { total: 100.1234 } },
        { metric_name: 'data_protection_calls', metric_data: { total: 22 } },
        { metric_name: 'managed_keys', metric_data: { total: 44 } },
        {
          metric_name: 'ssh_units',
          metric_data: {
            total: 100.2468,
            metric_details: [
              { type: 'otp_units', count: 50.1234 },
              { type: 'certificate_units', count: 51.22 },
            ],
          },
        },
        {
          metric_name: 'managed_keys',
          metric_data: {
            total: 82,
            metric_details: [
              { type: 'totp', count: 2 },
              { type: 'kmse', count: 1 },
            ],
          },
        },
        /**
         * Other metrics:
         * - external_plugins
         * - managed_keys (types: totp, kmse)
         * - data_protection_calls (types: transit, transform)
         * - id_token_units (types: oidc, spiffe) // not adding in 2.0.0
         *
         * Additional metrics to be added for new features in Vault 1.22/2.0.
         */
      ],
    },
  ],
};

const SELECTORS = {
  metricDetail: (metricKey) => `[data-test-metric-detail="${metricKey}"]`,
  metricDetailValue: (metricKey) => `[data-test-metric-detail-value="${metricKey}"]`,
};

module('Acceptance | billing/overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.version = this.owner.lookup('service:version');
    this.server.get('/sys/billing/overview', () => mockMetrics);

    // Stub the API service
    const api = this.owner.lookup('service:api');
    this.billingStub = sinon.stub(api.sys, 'systemReadBillingOverview').resolves(mockMetrics);
  });

  hooks.afterEach(function () {
    this.billingStub?.restore();
  });

  test('display billing/overview when license endpoint has consumption billing', async function (assert) {
    this.server.get('/sys/license/features', () => ({ features: ['Consumption Billing'] }));
    await login();

    assert.dom(GENERAL.navLink('Billing metrics')).exists('Billing metrics nav link is present');
    assert.dom(GENERAL.navLink('Billing metrics')).hasText('Billing metrics');
    await click(GENERAL.navLink('Billing metrics'));
    assert.strictEqual(currentURL(), '/vault/billing/overview');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Billing metrics');
    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .hasText(
        'Data reflects usage across this Vault cluster. Billing metrics are used in license utilization.'
      );
    assert.dom(GENERAL.cardContainer('Summary')).exists();

    assert.dom(GENERAL.cardContainer('Secrets')).exists();

    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.STATIC_SECRETS_KV)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.STATIC_SECRETS_KV)).hasText('10');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.DYNAMIC_ROLES_TOTAL)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.DYNAMIC_ROLES_TOTAL)).hasText('60');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.AUTO_ROTATED_ROLES_TOTAL)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.AUTO_ROTATED_ROLES_TOTAL)).hasText('30');

    assert.dom(GENERAL.cardContainer('Credential units')).exists();
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.PKI_UNITS_TOTAL)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.PKI_UNITS_TOTAL)).hasText('100.1234');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.SSH_UNITS_OTP_UNITS)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.SSH_UNITS_OTP_UNITS)).hasText('50.1234');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.SSH_UNITS_CERTIFICATE_UNITS)).exists();
    assert
      .dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.SSH_UNITS_CERTIFICATE_UNITS))
      .hasText('50.1234');

    assert.dom(GENERAL.cardContainer('Data protection calls')).exists();
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSFORM)).exists();
    assert
      .dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSFORM))
      .hasText('0');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSIT)).exists();
    assert
      .dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSIT))
      .hasText('0');

    assert.dom(GENERAL.cardContainer('Managed keys')).exists();
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.MANAGED_KEYS_TOTP)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.MANAGED_KEYS_TOTP)).hasText('52');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.MANAGED_KEYS_KMSE)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.MANAGED_KEYS_KMSE)).hasText('30');
    await logout();
  });
  test('hide billing/overview when license endpoint does not have consumption billing', async function (assert) {
    this.server.get('/sys/license/features', () => ({ features: [] }));
    await login();

    assert.dom(GENERAL.navLink('Billing metrics')).doesNotExist('Billing metrics nav link is not present');
    await logout();
  });
});
