/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { visit, currentURL, waitFor } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

const mockedEmptyResponse = {
  data: {
    auth_methods: {},
    kvv1_secrets: 0,
    kvv2_secrets: 0,
    lease_count_quotas: {},
    leases_by_auth_method: {},
    replication_status: {},
    secret_engines: {},
  },
};

const mockedResponseWithData = {
  data: {
    auth_methods: { cats: 42 },
    kvv1_secrets: 60,
    kvv2_secrets: 40,
    lease_count_quotas: {
      global_lease_count_quota: { capacity: 420000, count: 210000, name: 'default' },
      total_lease_count_quotas: 1,
    },
    namespaces: 1,
    replication_status: {
      dr_primary: true,
      dr_state: 'enabled',
      pr_primary: false,
      pr_state: 'enabled',
    },
    secret_engines: { dogs: 43 },
  },
};

module('Acceptance | vault-reporting', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
  });

  test('it visits the usage reporting dashboard and renders the header', async function (assert) {
    await visit('/vault/usage-reporting');
    assert.strictEqual(currentURL(), '/vault/usage-reporting', 'navigates to usage reporting dashboard');
    assert.dom('.hds-page-header').includesText('Vault Usage', 'renders the "Vault Usage" header');
  });

  test('it renders the counters dashboard block with all expected counters', async function (assert) {
    this.server.get('http://localhost:7357/v1/sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-counters]');
    assert.dom('[data-test-dashboard-counters]').exists('renders the counters dashboard block');

    const expectedCounters = ['Child namespaces', 'KV secrets', 'Secrets sync', 'PKI roles'];

    expectedCounters.forEach((counterLabel) => {
      assert.dom(`[data-test-counter="${counterLabel}"]`).exists(`counter "${counterLabel}" is rendered`);
    });

    assert.dom('[data-test-counter="Child namespaces"]').includesText('1');
    assert.dom('[data-test-counter="KV secrets"]').includesText('100');
  });

  test('dashboard card: Secret engines', async function (assert) {
    this.server.get('http://localhost:7357/v1/sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-secret-engines]');

    const card = document.querySelector('[data-test-dashboard-secret-engines]');
    assert.ok(card, 'renders Secret engines card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Secret engines', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(link.getAttribute('href'), 'secrets', 'link points to secrets');

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Enabled secret engines for this cluster.',
      'description is correct'
    );

    assert.dom('[data-test-dashboard-secret-engines]').includesText('dogs 43');
  });

  test('dashboard card: Authentication methods', async function (assert) {
    this.server.get('http://localhost:7357/v1/sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-auth-methods]');

    const card = document.querySelector('[data-test-dashboard-auth-methods]');
    assert.ok(card, 'renders Authentication methods card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Authentication methods', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(link.getAttribute('href'), 'access', 'link points to access');

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Enabled authentication methods for this cluster.',
      'description is correct'
    );

    assert.dom('[data-test-dashboard-auth-methods]').includesText('cats 42');
  });

  test('dashboard card: Global lease count quota', async function (assert) {
    this.server.get('http://localhost:7357/v1/sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-lease-count]');

    const card = document.querySelector('[data-test-dashboard-lease-count]');
    assert.ok(card, 'renders Global lease count quota card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Global lease count quota', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(
      link.getAttribute('href'),
      'https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota',
      'link points to lease count docs'
    );

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Total number of active leases for this quota.',
      'description is correct'
    );

    // Check if the percentage and count are correct
    assert.dom('[data-test-global-lease-percentage-text]').hasText('50%', 'percentage is correct');
    assert.dom('[data-test-global-lease-count-text]').hasText('210K / 420K', 'lease count is correct');
  });

  test('dashboard card: Cluster replication status', async function (assert) {
    this.server.get('http://localhost:7357/v1/sys/utilization-report', () => mockedResponseWithData);
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-cluster-replication]');

    const card = document.querySelector('[data-test-dashboard-cluster-replication]');
    assert.ok(card, 'renders Cluster replication status card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Cluster replication', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(link.getAttribute('href'), 'replication', 'link points to replication');

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Status of disaster recovery and performance replication.',
      'description is correct'
    );
  });

  test('empty states display expected text', async function (assert) {
    this.server.get('http://localhost:7357/v1/sys/utilization-report', () => mockedEmptyResponse);
    await visit('/vault/usage-reporting');

    // Secret Engines
    await waitFor('[data-test-dashboard-secret-engines]');

    assert
      .dom('[data-test-dashboard-secret-engines]')
      .includesText('None enabled', 'Secret engines empty state: title is shown');
    assert
      .dom('[data-test-dashboard-secret-engines]')
      .includesText(
        'Secret engines in this namespace will appear here.',
        'Secret engines empty state: body is shown'
      );
    assert
      .dom('[data-test-dashboard-secret-engines]')
      .includesText('Enable secret engines', 'Secret engines empty state: CTA is shown');

    // Auth Methods
    await waitFor('[data-test-dashboard-auth-methods]');

    assert
      .dom('[data-test-dashboard-auth-methods]')
      .includesText('None enabled', 'Auth methods empty state: title is shown');
    assert
      .dom('[data-test-dashboard-auth-methods]')
      .includesText(
        'Authentication methods in this namespace will appear here.',
        'Auth methods empty state: body is shown'
      );
    assert
      .dom('[data-test-dashboard-auth-methods]')
      .includesText('Enable authentication methods', 'Auth methods empty state: CTA is shown');

    // Lease Count Quota
    await waitFor('[data-test-dashboard-lease-count]');

    assert
      .dom('[data-test-dashboard-lease-count]')
      .includesText(
        `Lease quotas enforce limits on active secrets and tokens. It's recommended to enable this to protect stability for this Vault cluster.`,
        'Lease quota empty state description is shown'
      );

    assert
      .dom('[data-test-dashboard-lease-count]')
      .includesText('Global lease count quota', 'Lease quota empty state: docs link is shown');
  });

  //
});
