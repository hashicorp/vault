/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, { STATIC_NOW, LICENSE_START, UPGRADE_DATE } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { visit, click, findAll, settled } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);

module('Acceptance | clients | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it should render the correct tabs', async function (assert) {
    assert.dom(GENERAL.tab('overview')).exists();
    assert.dom(GENERAL.tab('token')).exists();
    assert.dom(GENERAL.tab('acme')).exists();
  });

  test('it should render charts', async function (assert) {
    assert
      .dom(CLIENT_COUNT.counts.startMonth)
      .hasText('July 2023', 'billing start month is correctly parsed from license');
    assert
      .dom(CLIENT_COUNT.rangeDropdown)
      .hasText('Jul 2023 - Jan 2024', 'Date range shows dates correctly parsed activity response');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');
    assert
      .dom(CLIENT_COUNT.chartContainer('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(CLIENT_COUNT.charts.line.xAxisLabel)
      .hasText('7/23', 'x-axis labels start with billing start date');
    assert.strictEqual(
      findAll('[data-test-line-chart="plot-point"]').length,
      5,
      'line chart plots 5 points to match query'
    );
  });

  test('it should update charts when querying date ranges', async function (assert) {
    // query for single, historical month with no new counts (July 2023)
    await click(CLIENT_COUNT.rangeDropdown);
    await click('[data-test-show-calendar]');
    await click('[data-test-previous-year]');
    await click(`[data-test-calendar-month=${ARRAY_OF_MONTHS[LICENSE_START.getMonth()]}]`);
    assert
      .dom(CLIENT_COUNT.runningTotalMonthStats)
      .doesNotExist('running total single month stat boxes do not show');
    assert
      .dom(CLIENT_COUNT.chartContainer('Vault client counts'))
      .doesNotExist('running total month over month charts do not show');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('attribution area shows');
    assert
      .dom('[data-test-chart-container="new-clients"] [data-test-component="empty-state"]')
      .exists('new client attribution has empty state');
    assert
      .dom('[data-test-empty-state-subtext]')
      .hasText('There are no new clients for this namespace during this time period.    ');
    assert.dom('[data-test-chart-container="total-clients"]').exists('total client attribution chart shows');

    // reset to billing period
    await click(CLIENT_COUNT.rangeDropdown);
    await click('[data-test-current-billing-period]');

    // change billing start to month/year of upgrade to 1.10
    await click(CLIENT_COUNT.counts.startEdit);
    await click(CLIENT_COUNT.monthDropdown);
    await click(`[data-test-dropdown-month="${ARRAY_OF_MONTHS[UPGRADE_DATE.getMonth()]}"]`);
    await click(CLIENT_COUNT.yearDropdown);
    await click(`[data-test-dropdown-year="${UPGRADE_DATE.getFullYear()}"]`);
    await click('[data-test-date-dropdown-submit]');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');
    assert
      .dom(CLIENT_COUNT.chartContainer('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(CLIENT_COUNT.charts.line.xAxisLabel)
      .hasText('9/23', 'x-axis labels start with queried start month (upgrade date)');
    assert.strictEqual(
      findAll('[data-test-line-chart="plot-point"]').length,
      5,
      'line chart plots 5 points to match query'
    );

    // query for single, historical month (upgrade month)
    await click(CLIENT_COUNT.rangeDropdown);
    await click('[data-test-show-calendar]');
    assert.dom('[data-test-display-year]').hasText('2024');
    await click('[data-test-previous-year]');
    await click('[data-test-calendar-month="September"]');
    assert.dom(CLIENT_COUNT.runningTotalMonthStats).exists('running total single month stat boxes show');
    assert
      .dom(CLIENT_COUNT.chartContainer('Vault client counts'))
      .doesNotExist('running total month over month charts do not show');
    assert.dom(CLIENT_COUNT.attributionBlock).exists('attribution area shows');
    assert.dom('[data-test-chart-container="new-clients"]').exists('new client attribution chart shows');
    assert.dom('[data-test-chart-container="total-clients"]').exists('total client attribution chart shows');

    // query historical date range (from September 2023 to December 2023)
    await click(CLIENT_COUNT.rangeDropdown);
    await click('[data-test-show-calendar]');
    await click('[data-test-calendar-month="December"]');

    assert.dom(CLIENT_COUNT.attributionBlock).exists('Shows attribution area');
    assert
      .dom(CLIENT_COUNT.chartContainer('Vault client counts'))
      .exists('Shows running totals with monthly breakdown charts');
    assert.strictEqual(
      findAll('[data-test-line-chart="plot-point"]').length,
      4,
      'line chart plots 4 points to match query'
    );
    const xAxisLabels = findAll(CLIENT_COUNT.charts.line.xAxisLabel);
    assert
      .dom(xAxisLabels[xAxisLabels.length - 1])
      .hasText('12/23', 'x-axis labels end with queried end month');

    // reset to billing period
    await click(CLIENT_COUNT.rangeDropdown);
    await click('[data-test-current-billing-period]');
    // query month older than count start date
    await click(CLIENT_COUNT.counts.startEdit);
    await click(CLIENT_COUNT.monthDropdown);
    await click(`[data-test-dropdown-month="${ARRAY_OF_MONTHS[LICENSE_START.getMonth()]}"]`);
    await click(CLIENT_COUNT.yearDropdown);
    await click(`[data-test-dropdown-year="${LICENSE_START.getFullYear() - 3}"]`);
    await click('[data-test-date-dropdown-submit]');
    assert
      .dom(CLIENT_COUNT.counts.startDiscrepancy)
      .hasTextContaining(
        'You requested data from July 2020. We only have data from January 2023, and that is what is being shown here.',
        'warning banner displays that date queried was prior to count start date'
      );
  });

  test('totals filter correctly with full data', async function (assert) {
    assert
      .dom(CLIENT_COUNT.chartContainer('Vault client counts'))
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
      .dom(CLIENT_COUNT.charts.statTextValue('Entity clients'))
      .includesText(`${formatNumber([topNamespace.entity_clients])}`, 'total entity clients is accurate');
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Non-entity clients'))
      .includesText(
        `${formatNumber([topNamespace.non_entity_clients])}`,
        'total non-entity clients is accurate'
      );
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Secrets sync clients'))
      .includesText(`${formatNumber([topNamespace.secret_syncs])}`, 'total sync clients is accurate');
    assert
      .dom('[data-test-attribution-clients] p')
      .includesText(`${formatNumber([topMount.clients])}`, 'top attribution clients accurate');

    // FILTER BY AUTH METHOD
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();

    assert.ok(true, 'Filter by first auth method');
    assert.dom(CLIENT_COUNT.selectedAuthMount).hasText(topMount.label, 'selects top mount');
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Entity clients'))
      .includesText(`${formatNumber([topMount.entity_clients])}`, 'total entity clients is accurate');
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Non-entity clients'))
      .includesText(`${formatNumber([topMount.non_entity_clients])}`, 'total non-entity clients is accurate');
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Secrets sync clients'))
      .includesText(`${formatNumber([topMount.secret_syncs])}`, 'total sync clients is accurate');
    assert.dom(CLIENT_COUNT.attributionBlock).doesNotExist('Does not show attribution block');

    await click('#namespace-search-select [data-test-selected-list-button="delete"]');
    assert.ok(true, 'Remove namespace filter without first removing auth method filter');
    assert.dom('[data-test-top-attribution]').includesText('Top namespace');
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Entity clients'))
      .hasTextContaining(
        `${formatNumber([response.total.entity_clients])}`,
        'total entity clients is back to unfiltered value'
      );
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Non-entity clients'))
      .hasTextContaining(
        `${formatNumber([formatNumber([response.total.non_entity_clients])])}`,
        'total non-entity clients is back to unfiltered value'
      );
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Secrets sync clients'))
      .hasTextContaining(
        `${formatNumber([formatNumber([response.total.secret_syncs])])}`,
        'total sync clients is back to unfiltered value'
      );
    assert
      .dom('[data-test-attribution-clients]')
      .hasTextContaining(
        `${formatNumber([topNamespace.clients])}`,
        'top attribution clients back to unfiltered value'
      );
  });
});

module('Acceptance | clients | overview | sync in license, activated', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');

    // add feature to license
    this.server.get('/sys/license/features', () => ({ features: ['Secrets Sync'] }));
    // activate feature
    this.server.get('/sys/activation-flags', () => ({
      data: { activated: ['secrets-sync'], unactivated: [] },
    }));

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it should render the correct tabs', async function (assert) {
    assert.dom(GENERAL.tab('sync')).exists();
  });

  test('it should show secrets sync data in overview and tab', async function (assert) {
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Secrets sync clients'))
      .exists('shows secret sync data on overview');
    await click(GENERAL.tab('sync'));

    assert.dom(GENERAL.tab('sync')).hasClass('active');
    assert.dom(GENERAL.emptyStateTitle).doesNotExist();
    assert
      .dom(CLIENT_COUNT.charts.chart('Secrets sync usage'))
      .exists('chart is shown because feature is active and has data');
  });
});

module('Acceptance | clients | overview | sync in license, not activated', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.server.get('/sys/license/features', () => ({ features: ['Secrets Sync'] }));

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should show the secrets sync tab', async function (assert) {
    assert.dom(GENERAL.tab('sync')).exists('sync tab is shown because feature is in license');
  });

  test('it should hide secrets sync charts', async function (assert) {
    assert
      .dom(CLIENT_COUNT.charts.chart('Secrets sync usage'))
      .doesNotExist('chart is hidden because feature is not activated');

    assert.dom('[data-test-stat-text="secret-syncs"]').doesNotExist();
  });
});

module('Acceptance | clients | overview | sync not in license', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    // mocks endpoint for no additional license modules
    this.server.get('/sys/license/features', () => ({ features: [] }));

    await authPage.login();
    return visit('/vault/clients/counts/overview');
  });

  test('it should hide the secrets sync tab', async function (assert) {
    assert.dom(GENERAL.tab('sync')).doesNotExist();
  });

  test('it should hide secrets sync charts', async function (assert) {
    assert.dom(CLIENT_COUNT.charts.chart('Secrets sync usage')).doesNotExist();

    assert.dom('[data-test-stat-text="secret-syncs"]').doesNotExist();
  });
});
