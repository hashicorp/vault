/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MOCK_DATA = {
  authMethods: { aws: 42, userpass: 43, kubernetes: 44 },
  leasesByAuthMethod: {},
  kvv1Secrets: 60,
  kvv2Secrets: 40,
  leaseCountQuotas: {
    globalLeaseCountQuota: { capacity: 420000, count: 210000, name: 'default' },
    totalLeaseCountQuotas: 1,
  },
  namespaces: 1,
  pki: { totalIssuers: 0, totalRoles: 5 },
  replicationStatus: { drPrimary: true, drState: 'primary', prPrimary: false, prState: 'disabled' },
  secretEngines: { cubbyhole: 45, nomad: 46, aws: 47 },
  secretSync: { totalDestinations: 1, destinations: { aws: 1 } },
};

const MOCK_NAMESPACE_DATA = { keys: ['root', 'ns1'] };

module('Integration | Component | usage-reporting/views/dashboard', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the Vault Usage page header', async function (assert) {
    this.set('fetchUsageData', async () => MOCK_DATA);
    this.set('fetchNamespaceData', async () => MOCK_NAMESPACE_DATA);
    await render(hbs`
      <UsageReporting::Views::Dashboard
        @onFetchUsageData={{this.fetchUsageData}}
        @onFetchNamespaceData={{this.fetchNamespaceData}}
      />
    `);
    await waitFor('[data-test-usage-dashboard-container]');
    assert.dom('[data-test-usage-dashboard-container]').exists('dashboard container renders');
  });

  test('it renders the counters section', async function (assert) {
    this.set('fetchUsageData', async () => MOCK_DATA);
    this.set('fetchNamespaceData', async () => MOCK_NAMESPACE_DATA);
    await render(hbs`
      <UsageReporting::Views::Dashboard
        @onFetchUsageData={{this.fetchUsageData}}
        @onFetchNamespaceData={{this.fetchNamespaceData}}
      />
    `);
    await waitFor('[data-test-vault-reporting-dashboard-counters]');
    assert
      .dom('[data-test-vault-reporting-counter="Child namespaces"]')
      .exists('Child namespaces counter renders');
    assert
      .dom('[data-test-vault-reporting-counter="KV secrets"]')
      .exists('KV secrets counter renders')
      .includesText('100', 'KV secrets shows combined v1+v2 count');
    assert.dom('[data-test-vault-reporting-counter="PKI roles"]').exists('PKI roles counter renders');
  });

  test('it renders the namespace picker when namespace keys are present', async function (assert) {
    this.set('fetchUsageData', async () => MOCK_DATA);
    this.set('fetchNamespaceData', async () => MOCK_NAMESPACE_DATA);
    await render(hbs`
      <UsageReporting::Views::Dashboard
        @onFetchUsageData={{this.fetchUsageData}}
        @onFetchNamespaceData={{this.fetchNamespaceData}}
      />
    `);
    await waitFor('[data-test-vault-reporting-namespace-picker]');
    assert
      .dom('[data-test-vault-reporting-namespace-picker]')
      .exists('namespace picker renders when keys are available');
  });

  test('it does not render the namespace picker when there are no namespace keys', async function (assert) {
    this.set('fetchUsageData', async () => MOCK_DATA);
    this.set('fetchNamespaceData', async () => ({ keys: [] }));
    await render(hbs`
      <UsageReporting::Views::Dashboard
        @onFetchUsageData={{this.fetchUsageData}}
        @onFetchNamespaceData={{this.fetchNamespaceData}}
      />
    `);
    await waitFor('[data-test-vault-reporting-dashboard-counters]');
    assert
      .dom('[data-test-vault-reporting-namespace-picker]')
      .doesNotExist('namespace picker is hidden when no keys');
  });

  test('it renders cluster-level cards', async function (assert) {
    this.set('fetchUsageData', async () => MOCK_DATA);
    this.set('fetchNamespaceData', async () => MOCK_NAMESPACE_DATA);
    await render(hbs`
      <UsageReporting::Views::Dashboard
        @onFetchUsageData={{this.fetchUsageData}}
        @onFetchNamespaceData={{this.fetchNamespaceData}}
      />
    `);
    await waitFor('[data-test-vault-reporting-dashboard-cluster-viz-blocks]');
    assert
      .dom('[data-test-vault-reporting-dashboard-cluster-replication]')
      .exists('cluster replication card renders');
    assert.dom('[data-test-vault-reporting-dashboard-lease-count]').exists('global lease count card renders');
    assert.dom('[data-test-vault-reporting-dashboard-secrets-sync]').exists('secrets sync card renders');
  });

  test('it renders an inline error alert when fetch fails', async function (assert) {
    this.set('fetchUsageData', async () => {
      throw new Error('API error');
    });
    this.set('fetchNamespaceData', async () => MOCK_NAMESPACE_DATA);
    await render(hbs`
      <UsageReporting::Views::Dashboard
        @onFetchUsageData={{this.fetchUsageData}}
        @onFetchNamespaceData={{this.fetchNamespaceData}}
      />
    `);
    await waitFor('[data-test-vault-reporting-dashboard-error]');
    assert.dom('[data-test-vault-reporting-dashboard-error]').exists('error alert renders on fetch failure');
    assert
      .dom('[data-test-vault-reporting-dashboard-error-description]')
      .hasText('An error occurred, please try again.', 'renders error description');
  });
});
