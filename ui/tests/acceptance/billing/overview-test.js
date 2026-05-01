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
import { dateFormat } from 'core/helpers/date-format';
import { METRICS_DATA_RESPONSE } from 'vault/tests/helpers/billing/stubs';

const SELECTORS = {
  metricDetail: (metricKey) => `[data-test-metric-detail="${metricKey}"]`,
  metricDetailValue: (metricKey) => `[data-test-metric-detail-value="${metricKey}"]`,
};

module('Acceptance | billing/overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.version = this.owner.lookup('service:version');
    this.mockMetrics = METRICS_DATA_RESPONSE.data;
    this.server.get('/sys/billing/overview', () => this.mockMetrics);

    // Stub the API service
    const api = this.owner.lookup('service:api');
    this.billingStub = sinon.stub(api.sys, 'systemReadBillingOverview').resolves(this.mockMetrics);
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
        'Data reflects usage across this Vault cluster. Billing metrics determine license utilization.'
      );
    assert.dom(GENERAL.textBody('Last updated date time')).hasText(
      `Values update every 10 minutes. Last updated: ${dateFormat(
        [this.mockMetrics.months[0].updated_at, 'MMMM d, yyyy, hh:mm:ss aaa'],
        {
          withTimeZone: true,
        }
      )}`
    );

    assert.dom(GENERAL.cardContainer('Summary')).exists();

    assert.dom(GENERAL.cardContainer('Secrets')).exists();

    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.STATIC_SECRETS_KV)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.STATIC_SECRETS_KV)).hasText('10');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.DYNAMIC_ROLES_TOTAL)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.DYNAMIC_ROLES_TOTAL)).hasText('130');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.AUTO_ROTATED_ROLES_TOTAL)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.AUTO_ROTATED_ROLES_TOTAL)).hasText('70');

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
      .hasText('220');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSIT)).exists();
    assert
      .dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSIT))
      .hasText('200');

    assert.dom(GENERAL.cardContainer('Managed keys')).exists();
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.MANAGED_KEYS_TOTP)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.MANAGED_KEYS_TOTP)).hasText('220');
    assert.dom(SELECTORS.metricDetail(NormalizedBillingMetrics.MANAGED_KEYS_KMSE)).exists();
    assert.dom(SELECTORS.metricDetailValue(NormalizedBillingMetrics.MANAGED_KEYS_KMSE)).hasText('210');
    await logout();
  });

  test('display no data available when updated_at is invalid', async function (assert) {
    this.server.get('/sys/license/features', () => ({ features: ['Consumption Billing'] }));
    const mockMetricsInvalidDate = { ...this.mockMetrics };
    mockMetricsInvalidDate.months[1].updated_at = '0001-01-01T00:00:00Z';
    this.server.get('/sys/billing/overview', () => mockMetricsInvalidDate);
    await login();
    assert.dom(GENERAL.navLink('Billing metrics')).hasText('Billing metrics');
    await click(GENERAL.navLink('Billing metrics'));
    await click(GENERAL.dropdownToggle('Date range'));
    await click(GENERAL.menuItem('2025-12'));
    assert.dom(GENERAL.textBody('Last updated date time')).hasText('No data available.');
    await logout();
  });

  test('hide billing/overview when license endpoint does not have consumption billing', async function (assert) {
    this.server.get('/sys/license/features', () => ({ features: [] }));
    await login();

    assert.dom(GENERAL.navLink('Billing metrics')).doesNotExist('Billing metrics nav link is not present');
    await logout();
  });
});
