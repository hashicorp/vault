/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MOCK_DATA = {
  authMethods: { aws: 42, userpass: 43 },
  leasesByAuthMethod: {},
  kvv1Secrets: 60,
  kvv2Secrets: 40,
  leaseCountQuotas: {
    globalLeaseCountQuota: { capacity: 420000, count: 210000, name: 'default' },
    totalLeaseCountQuotas: 1,
  },
  namespaces: 1,
  pki: { totalIssuers: 0, totalRoles: 5 },
  replicationStatus: { drPrimary: true, drState: 'enabled', prPrimary: false, prState: 'disabled' },
  secretEngines: { cubbyhole: 45, aws: 47 },
  secretSync: { totalDestinations: 1, destinations: { aws: 1 } },
};

module('Integration | Component | usage-reporting/dashboard/export', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the export toggle button when @data is provided', async function (assert) {
    this.set('data', MOCK_DATA);
    await render(hbs`<UsageReporting::Dashboard::Export @data={{this.data}} />`);
    assert.dom('[data-test-vault-reporting-export-toggle]').exists('export toggle button renders');
  });

  test('it does not render when @data is not provided', async function (assert) {
    await render(hbs`<UsageReporting::Dashboard::Export />`);
    assert
      .dom('[data-test-vault-reporting-export-toggle]')
      .doesNotExist('export button is hidden when no data');
  });

  test('it renders JSON and CSV download options in the dropdown', async function (assert) {
    this.set('data', MOCK_DATA);
    await render(hbs`<UsageReporting::Dashboard::Export @data={{this.data}} />`);
    await click('[data-test-vault-reporting-export-toggle]');
    assert.dom('[data-test-vault-reporting-export-json]').exists('JSON download option renders');
    assert.dom('[data-test-vault-reporting-export-csv]').exists('CSV download option renders');
  });

  test('JSON download link has the correct filename attribute', async function (assert) {
    this.set('data', MOCK_DATA);
    await render(hbs`<UsageReporting::Dashboard::Export @data={{this.data}} />`);
    await click('[data-test-vault-reporting-export-toggle]');
    assert
      .dom('[data-test-vault-reporting-export-json]')
      .hasAttribute('download', 'vault-usage-dashboard.json', 'JSON link has correct download filename');
  });

  test('CSV download link has the correct filename attribute', async function (assert) {
    this.set('data', MOCK_DATA);
    await render(hbs`<UsageReporting::Dashboard::Export @data={{this.data}} />`);
    await click('[data-test-vault-reporting-export-toggle]');
    assert
      .dom('[data-test-vault-reporting-export-csv]')
      .hasAttribute('download', 'vault-usage-dashboard.csv', 'CSV link has correct download filename');
  });

  test('CSV download content uses sentence case labels', async function (assert) {
    assert.expect(10);

    this.set('data', {
      ...MOCK_DATA,
      authMethods: { ...MOCK_DATA.authMethods, ldap: 2 },
      secretEngines: { ...MOCK_DATA.secretEngines, gcp: 9, rabbitmq: 3 },
      secretSync: {
        ...MOCK_DATA.secretSync,
        destinations: { ...MOCK_DATA.secretSync.destinations, gcp: 2 },
      },
    });
    await render(hbs`<UsageReporting::Dashboard::Export @data={{this.data}} />`);
    await click('[data-test-vault-reporting-export-toggle]');

    const csvHref = document.querySelector('[data-test-vault-reporting-export-csv]').getAttribute('href');
    const csvText = decodeURIComponent(csvHref.split(',')[1]);

    assert.true(csvText.includes('"Metric","Count/breakdown"'), 'uses sentence case header');
    assert.true(csvText.includes('"Child namespaces","1"'), 'uses sentence case for child namespaces');
    assert.true(csvText.includes('"Total KV secrets","100"'), 'uses sentence case for total KV secrets');
    assert.true(csvText.includes('"PKI roles","5"'), 'uses sentence case for PKI roles');
    assert.true(csvText.includes('"Secret engine AWS","47"'), 'capitalizes AWS in nested secret engine rows');
    assert.true(csvText.includes('"Secret engine GCP","9"'), 'capitalizes GCP in nested secret engine rows');
    assert.true(
      csvText.includes('"Secret engine RabbitMQ","3"'),
      'applies branded RabbitMQ override in nested secret engine rows'
    );
    assert.true(
      csvText.includes('"Auth method userpass","43"'),
      'keeps non-acronym nested auth method rows unchanged'
    );
    assert.true(csvText.includes('"Auth method LDAP","2"'), 'capitalizes LDAP in nested auth method rows');
    assert.true(
      csvText.includes('"Secrets sync destination AWS","1"'),
      'capitalizes AWS in nested secrets sync rows'
    );
  });
});
