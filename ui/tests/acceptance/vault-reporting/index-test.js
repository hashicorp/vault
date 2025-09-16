/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { visit, currentURL, waitFor, click, fillIn } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { mockedResponseWithData, mockedEmptyResponse } from 'vault/tests/helpers/vault-usage/mocks';
import { createPolicyCmd, createTokenCmd, runCmd } from 'vault/tests/helpers/commands';

const loginWithReportingToken = async (capability = 'read') => {
  const policyName = 'show-vault-reporting';
  const policy = `
      path "sys/utilization-report" {
        capabilities = ["${capability}"]
      }
    `;

  const commands = [createPolicyCmd(policyName, policy), createTokenCmd(policyName)];
  // Use lower privileged token
  const token = await runCmd(commands);
  await login(token);
};

module('Acceptance | enterprise vault-reporting', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
  });

  test('it visits the usage reporting dashboard and renders the header', async function (assert) {
    // Log in with lower privileged token
    await loginWithReportingToken('read');
    await visit('/vault/dashboard');
    await click('[data-test-sidebar-nav-link="Vault Usage"]');
    assert.strictEqual(currentURL(), '/vault/usage-reporting', 'navigates to usage reporting dashboard');
    assert.dom('.hds-page-header').includesText('Vault Usage', 'renders the "Vault Usage" header');
  });

  test('it hides the nav item if policy does not allow access to sys/utilization-report', async function (assert) {
    await loginWithReportingToken('deny');
    await visit('/vault/dashboard');
    assert.dom('[data-test-sidebar-nav-link="Vault Usage"]').doesNotExist('sidebar nav link is hidden');
  });

  test('it renders the counters dashboard block with all expected counters', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-vault-reporting-dashboard-counters]');
    assert
      .dom('[data-test-vault-reporting-dashboard-counters]')
      .exists('renders the counters dashboard block');

    const expectedCounters = ['Child namespaces', 'KV secrets', 'PKI roles'];

    expectedCounters.forEach((counterLabel) => {
      assert
        .dom(`[data-test-vault-reporting-counter="${counterLabel}"]`)
        .exists(`counter "${counterLabel}" is rendered`);
    });

    assert.dom('[data-test-vault-reporting-counter="Child namespaces"]').includesText('1');
    assert.dom('[data-test-vault-reporting-counter="KV secrets"]').includesText('100');
  });

  test('dashboard card: Secret engines', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-vault-reporting-dashboard-secret-engines]');

    assert.dom('[data-test-vault-reporting-dashboard-secret-engines]').exists('renders Secret engines card');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-secret-engines] [data-test-vault-reporting-dashboard-card-title]'
      )
      .exists('title is present')
      .hasText('Secret engines', 'title is correct');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-secret-engines] [data-test-vault-reporting-dashboard-card-title-link]'
      )
      .exists('title link is present')
      .hasAttribute('href', 'secrets', 'link points to secrets');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-secret-engines] [data-test-vault-reporting-dashboard-card-description]'
      )
      .exists('description is present')
      .hasText('Enabled secret engines for this cluster.', 'description is correct');
    assert
      .dom('[data-test-vault-reporting-dashboard-secret-engines]')
      .includesText('aws nomad cubbyhole 47 46 45');
  });

  test('dashboard card: Authentication methods', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-vault-reporting-dashboard-auth-methods]');

    assert
      .dom('[data-test-vault-reporting-dashboard-auth-methods]')
      .exists('renders Authentication methods card');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-auth-methods] [data-test-vault-reporting-dashboard-card-title]'
      )
      .exists('title is present')
      .hasText('Authentication methods', 'title is correct');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-auth-methods] [data-test-vault-reporting-dashboard-card-title-link]'
      )
      .exists('title link is present')
      .hasAttribute('href', 'access', 'link points to access');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-auth-methods] [data-test-vault-reporting-dashboard-card-description]'
      )
      .exists('description is present')
      .hasText('Enabled authentication methods for this cluster.', 'description is correct');
    assert
      .dom('[data-test-vault-reporting-dashboard-auth-methods]')
      .includesText('kubernetes userpass aws 44 43 42');
  });

  test('dashboard card: Global lease count quota', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-vault-reporting-dashboard-lease-count]');

    assert
      .dom('[data-test-vault-reporting-dashboard-lease-count]')
      .exists('renders Global lease count quota card');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-lease-count] [data-test-vault-reporting-dashboard-card-title]'
      )
      .exists('title is present')
      .hasText('Global lease count quota', 'title is correct');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-lease-count] [data-test-vault-reporting-dashboard-card-title-link]'
      )
      .exists('title link is present')
      .hasAttribute(
        'href',
        'https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota',
        'link points to lease count docs'
      );
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-lease-count] [data-test-vault-reporting-dashboard-card-description]'
      )
      .exists('description is present')
      .hasText('Total number of active leases for this quota.', 'description is correct');
    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('50%', 'percentage is correct');
    assert
      .dom('[data-test-vault-reporting-global-lease-count-text]')
      .hasText('210K / 420K', 'lease count is correct');
  });

  test('dashboard card: Secrets sync', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-vault-reporting-dashboard-secrets-sync]');

    assert.dom('[data-test-vault-reporting-dashboard-secrets-sync]').exists('renders Secrets sync card');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-secrets-sync] [data-test-vault-reporting-dashboard-card-title]'
      )
      .hasText('Secrets sync', 'title is correct');

    assert
      .dom(
        '[data-test-vault-reporting-dashboard-secrets-sync] [data-test-vault-reporting-dashboard-card-description]'
      )
      .hasText(
        'Total number of destinations (e.g. third-party integrations) synced with secrets',
        'description is correct'
      );
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .includesText('Destination', 'Destinations header is present');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .includesText('aws: 1', 'aws destination is present');
  });

  test('dashboard card: Cluster replication status', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-vault-reporting-dashboard-cluster-replication]');

    assert
      .dom('[data-test-vault-reporting-dashboard-cluster-replication]')
      .exists('renders Cluster replication status card');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-cluster-replication] [data-test-vault-reporting-dashboard-card-title]'
      )
      .exists('title is present')
      .hasText('Cluster replication', 'title is correct');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-cluster-replication] [data-test-vault-reporting-dashboard-card-title-link]'
      )
      .exists('title link is present')
      .hasAttribute('href', 'replication', 'link points to replication');
    assert
      .dom(
        '[data-test-vault-reporting-dashboard-cluster-replication] [data-test-vault-reporting-dashboard-card-description]'
      )
      .exists('description is present')
      .hasText('Status of disaster recovery and performance replication.', 'description is correct');
  });

  test('empty states display expected text', async function (assert) {
    this.server.get('sys/utilization-report', () => mockedEmptyResponse);
    await visit('/vault/usage-reporting');

    await waitFor('[data-test-vault-reporting-dashboard-secret-engines]');
    assert
      .dom('[data-test-vault-reporting-dashboard-secret-engines]')
      .includesText('None enabled', 'Secret engines empty state: title is shown');
    assert
      .dom('[data-test-vault-reporting-dashboard-secret-engines]')
      .includesText(
        'Secret engines in this namespace will appear here.',
        'Secret engines empty state: body is shown'
      );
    assert
      .dom('[data-test-vault-reporting-dashboard-secret-engines]')
      .includesText('Enable secret engines', 'Secret engines empty state: CTA is shown');

    await waitFor('[data-test-vault-reporting-dashboard-auth-methods]');
    assert
      .dom('[data-test-vault-reporting-dashboard-auth-methods]')
      .includesText('None enabled', 'Auth methods empty state: title is shown');
    assert
      .dom('[data-test-vault-reporting-dashboard-auth-methods]')
      .includesText(
        'Authentication methods in this namespace will appear here.',
        'Auth methods empty state: body is shown'
      );
    assert
      .dom('[data-test-vault-reporting-dashboard-auth-methods]')
      .includesText('Enable authentication methods', 'Auth methods empty state: CTA is shown');

    await waitFor('[data-test-vault-reporting-dashboard-lease-count]');
    assert
      .dom('[data-test-vault-reporting-dashboard-lease-count]')
      .includesText(
        `Lease quotas enforce limits on active secrets and tokens. It's recommended to enable this to protect stability for this Vault cluster.`,
        'Lease quota empty state description is shown'
      );
    assert
      .dom('[data-test-vault-reporting-dashboard-lease-count]')
      .includesText('Global lease count quota', 'Lease quota empty state: docs link is shown');
  });

  test('namespace lookup functionality', async function (assert) {
    this.server.get('sys/internal/ui/namespaces', {
      data: {
        keys: ['child-ns1/', 'child-ns2/'],
      },
    });

    // Mock different responses for different namespaces
    const defaultMockResponse = mockedResponseWithData;
    const childNs1MockResponse = {
      data: {
        ...mockedResponseWithData.data,
        kvv2_secrets: 200,
        kvv1_secrets: 50,
      },
    };

    // Initially load default data
    this.server.get('sys/utilization-report', () => defaultMockResponse);

    const namespaceService = this.owner.lookup('service:namespace');
    namespaceService.set('accessibleNamespaces', ['child-ns1/', 'child-ns2/']);
    namespaceService.set('path', 'parent-ns');

    await visit('/vault/usage-reporting');

    // Verify initial KV secrets count (should be 40 from default mock)
    await waitFor('[data-test-vault-reporting-counter="KV secrets"]');
    assert
      .dom('[data-test-vault-reporting-counter="KV secrets"]')
      .includesText('100', 'Initial KV secrets count is 100');

    // Update mock to return different data for child-ns1
    this.server.get('sys/utilization-report', () => childNs1MockResponse);

    // Click the namespace picker dropdown to open it
    await click('.hds-dropdown.ssu-namespace-picker button');

    // Verify all namespaces are visible initially
    assert.dom('[data-test-vault-reporting-namespace-menu-item="root"]').exists('root namespace is visible');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="child-ns1"]')
      .exists('child-ns1 namespace is visible');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="child-ns2"]')
      .exists('child-ns2 namespace is visible');

    // Use the search bar to search for "ns1"
    await fillIn('[data-test-vault-reporting-namespace-search]', 'ns1');

    // Verify only child-ns1 is visible after search
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="child-ns1"]')
      .exists('child-ns1 is visible after search');

    // Verify root and child-ns2 are filtered out
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="root"]')
      .doesNotExist('root namespace is filtered out');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="child-ns2"]')
      .doesNotExist('child-ns2 namespace is filtered out');

    // Click on child-ns1 to select it
    await click('[data-test-vault-reporting-namespace-menu-item="child-ns1"]');

    // Verify that child-ns1 is now displayed as the selected namespace in the closed dropdown
    assert
      .dom('.hds-dropdown.ssu-namespace-picker')
      .includesText('child-ns1', 'child-ns1 is displayed as the selected namespace');

    // Wait for data to be refetched and verify the KV secrets count has changed
    await waitFor('[data-test-vault-reporting-counter="KV secrets"]');
    assert
      .dom('[data-test-vault-reporting-counter="KV secrets"]')
      .includesText('250', 'KV secrets count updated to 250 after selecting child-ns1');
  });
});
