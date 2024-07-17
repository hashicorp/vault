/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, { STATIC_NOW, LICENSE_START, UPGRADE_DATE } from 'vault/mirage/handlers/clients';
import syncHandler from 'vault/mirage/handlers/sync';
import sinon from 'sinon';
import { visit, click, findAll, settled, fillIn } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CHARTS, CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import ss from 'vault/tests/pages/components/search-select';
import { format } from 'date-fns';

const searchSelect = create(ss);

module('Acceptance | clients | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should render charts', async function (assert) {
    assert
      .dom(CLIENT_COUNT.dateRange.dateDisplay('start'))
      .hasText('July 2023', 'billing start month is correctly parsed from license');
    assert
      .dom(CLIENT_COUNT.dateRange.dateDisplay('end'))
      .hasText('January 2024', 'billing start month is correctly parsed from license');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');
    assert
      .dom(CHARTS.container('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(`${CHARTS.container('Vault client counts')} ${CHARTS.xAxisLabel}`)
      .hasText('7/23', 'x-axis labels start with billing start date');
    assert.strictEqual(findAll(CHARTS.plotPoint).length, 5, 'line chart plots 5 points to match query');
  });

  test('it should update charts when querying date ranges', async function (assert) {
    // query for single, historical month with no new counts (July 2023)
    const licenseStartMonth = format(LICENSE_START, 'yyyy-MM');
    const upgradeMonth = format(UPGRADE_DATE, 'yyyy-MM');
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), licenseStartMonth);
    await fillIn(CLIENT_COUNT.dateRange.editDate('end'), licenseStartMonth);

    await click(GENERAL.saveButton);
    assert
      .dom(CLIENT_COUNT.usageStats('Vault client counts'))
      .doesNotExist('running total single month stat boxes do not show');
    assert
      .dom(CHARTS.container('Vault client counts'))
      .doesNotExist('running total month over month charts do not show');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('attribution area shows');
    assert
      .dom(`${CHARTS.container('new-clients')} ${GENERAL.emptyStateTitle}`)
      .exists('new client attribution has empty state');
    assert
      .dom(GENERAL.emptyStateSubtitle)
      .hasText('There are no new clients for this namespace during this time period.    ');
    assert.dom(CHARTS.container('total-clients')).exists('total client attribution chart shows');

    // reset to billing period
    await click(CLIENT_COUNT.dateRange.edit);
    await click(CLIENT_COUNT.dateRange.reset);
    await click(GENERAL.saveButton);

    // change to start on month/year of upgrade to 1.10
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), upgradeMonth);
    await click(GENERAL.saveButton);

    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');
    assert
      .dom(CHARTS.container('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(`${CHARTS.container('Vault client counts')} ${CHARTS.xAxisLabel}`)
      .hasText('9/23', 'x-axis labels start with queried start month (upgrade date)');
    assert.strictEqual(findAll(CHARTS.plotPoint).length, 5, 'line chart plots 5 points to match query');

    // query for single, historical month (upgrade month)
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), upgradeMonth);
    await fillIn(CLIENT_COUNT.dateRange.editDate('end'), upgradeMonth);
    await click(GENERAL.saveButton);

    assert
      .dom(CLIENT_COUNT.usageStats('Vault client counts'))
      .exists('running total single month usage stats show');
    assert
      .dom(CHARTS.container('Vault client counts'))
      .doesNotExist('running total month over month charts do not show');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('attribution area shows');
    assert.dom(CHARTS.container('new-clients')).exists('new client attribution chart shows');
    assert.dom(CHARTS.container('total-clients')).exists('total client attribution chart shows');

    // query historical date range (from September 2023 to December 2023)
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2023-09');
    await fillIn(CLIENT_COUNT.dateRange.editDate('end'), '2023-12');
    await click(GENERAL.saveButton);

    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');
    assert
      .dom(CHARTS.container('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert.strictEqual(findAll(CHARTS.plotPoint).length, 4, 'line chart plots 4 points to match query');
    const xAxisLabels = findAll(CHARTS.xAxisLabel);
    assert
      .dom(xAxisLabels[xAxisLabels.length - 1])
      .hasText('12/23', 'x-axis labels end with queried end month');

    // reset to billing period
    await click(CLIENT_COUNT.dateRange.edit);
    await click(CLIENT_COUNT.dateRange.reset);
    await click(GENERAL.saveButton);

    // query month older than count start date
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2020-07');
    await click(GENERAL.saveButton);

    assert
      .dom(CLIENT_COUNT.counts.startDiscrepancy)
      .hasTextContaining(
        'You requested data from July 2020. We only have data from January 2023, and that is what is being shown here.',
        'warning banner displays that date queried was prior to count start date'
      );
  });

  test('totals filter correctly with full data', async function (assert) {
    // stub secrets sync being activated
    this.owner.lookup('service:flags').activatedFlags = ['secrets-sync'];

    assert
      .dom(CHARTS.container('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');

    const response = await this.store.peekRecord('clients/activity', 'some-activity-id');
    // FILTER BY NAMESPACE
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    const topNamespace = response.byNamespace[0];
    const topMount = topNamespace.mounts[0];

    assert.dom(CLIENT_COUNT.selectedNs).hasText(topNamespace.label, 'selects top namespace');
    assert.dom('[data-test-top-attribution]').includesText('Top auth method');
    assert
      .dom('[data-test-attribution-clients] p')
      .includesText(`${formatNumber([topMount.clients])}`, 'top attribution clients accurate');

    let expectedStats = {
      Entity: formatNumber([topNamespace.entity_clients]),
      'Non-entity': formatNumber([topNamespace.non_entity_clients]),
      ACME: formatNumber([topNamespace.acme_clients]),
      'Secret sync': formatNumber([topNamespace.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(CLIENT_COUNT.statTextValue(label))
        .includesText(`${expectedStats[label]}`, `label: ${label} renders accurate namespace client counts`);
    }

    // FILTER BY AUTH METHOD
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    assert.ok(true, 'Filter by first auth method');
    assert.dom(CLIENT_COUNT.selectedAuthMount).hasText(topMount.label, 'selects top mount');
    assert.dom(CLIENT_COUNT.attributionBlock).doesNotExist('Does not show attribution block');

    expectedStats = {
      Entity: formatNumber([topMount.entity_clients]),
      'Non-entity': formatNumber([topMount.non_entity_clients]),
      ACME: formatNumber([topMount.acme_clients]),
      'Secret sync': formatNumber([topMount.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(CLIENT_COUNT.statTextValue(label))
        .includesText(`${expectedStats[label]}`, `label: "${label} "renders accurate mount client counts`);
    }

    await click(GENERAL.searchSelect.removeSelected);
    assert.ok(true, 'Remove namespace filter without first removing auth method filter');
    assert.dom('[data-test-top-attribution]').includesText('Top namespace');
    assert
      .dom('[data-test-attribution-clients]')
      .hasTextContaining(
        `${formatNumber([topNamespace.clients])}`,
        'top attribution clients back to unfiltered value'
      );

    expectedStats = {
      Entity: formatNumber([response.total.entity_clients]),
      'Non-entity': formatNumber([response.total.non_entity_clients]),
      ACME: formatNumber([response.total.acme_clients]),
      'Secret sync': formatNumber([response.total.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(CLIENT_COUNT.statTextValue(label))
        .includesText(`${expectedStats[label]}`, `label: ${label} is back to unfiltered value`);
    }
  });
});

module('Acceptance | clients | overview | sync in license, activated', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));

    syncHandler(this.server);

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should render the correct tabs', async function (assert) {
    assert.dom(GENERAL.tab('sync')).exists('shows the sync tab');
  });

  test('it should show secrets sync stats', async function (assert) {
    assert.dom(CLIENT_COUNT.statTextValue('Secret sync')).exists('shows secret sync data on overview');
  });

  test('it should navigate to secrets sync page', async function (assert) {
    await click(GENERAL.tab('sync'));

    assert.dom(GENERAL.tab('sync')).hasClass('active');
    assert.dom(GENERAL.emptyStateTitle).doesNotExist();

    assert
      .dom(CHARTS.chart('Secrets sync usage'))
      .exists('chart is shown because feature is active and has data');
  });
});

module('Acceptance | clients | overview | sync in license, not activated', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.get('/sys/license/features', () => ({ features: ['Secrets Sync'] }));

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should show the secrets sync tab', async function (assert) {
    assert.dom(GENERAL.tab('sync')).exists('sync tab is shown because feature is in license');
  });

  test('it should hide secrets sync stats', async function (assert) {
    assert
      .dom(CLIENT_COUNT.statTextValue('Secret sync'))
      .doesNotExist('stat is hidden because feature is not activated');
    assert.dom(CLIENT_COUNT.statTextValue('Entity')).exists('other stats are still visible');
  });
});

module('Acceptance | clients | overview | sync not in license', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    // mocks endpoint for no additional license modules
    this.server.get('/sys/license/features', () => ({ features: [] }));

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should hide the secrets sync tab', async function (assert) {
    assert.dom(GENERAL.tab('sync')).doesNotExist();
  });

  test('it should hide secrets sync stats', async function (assert) {
    assert.dom(CLIENT_COUNT.statTextValue('Secret sync')).doesNotExist();
    assert.dom(CLIENT_COUNT.statTextValue('Entity')).exists('other stats are still visible');
  });
});

module('Acceptance | clients | overview | HVD', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    syncHandler(this.server);
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should show the secrets sync tab', async function (assert) {
    assert.dom(GENERAL.tab('sync')).exists();
  });

  test('it should show secrets sync stats', async function (assert) {
    assert.dom(CLIENT_COUNT.statTextValue('Secret sync')).exists();
  });
});
