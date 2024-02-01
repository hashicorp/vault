/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { visit, click, findAll, find, settled } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { SELECTORS, STATIC_NOW, LICENSE_START, UPGRADE_DATE } from 'vault/tests/helpers/clients';
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

  test('it should render charts', async function (assert) {
    assert
      .dom(SELECTORS.counts.startMonth)
      .hasText('July 2022', 'billing start month is correctly parsed from license');
    assert
      .dom(SELECTORS.rangeDropdown)
      .hasText(`Jul 2022 - Jan 2023`, 'Date range shows dates correctly parsed activity response');
    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert
      .dom(SELECTORS.charts.chart('running total'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(SELECTORS.charts.line.xAxisLabel)
      .hasText(`7/22`, 'x-axis labels start with billing start date');
    assert.strictEqual(
      findAll('[data-test-line-chart="plot-point"]').length,
      6,
      `line chart plots 6 points to match query`
    );
  });

  test('it should update charts when querying date ranges', async function (assert) {
    // query for single, historical month with no new counts
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-show-calendar]');
    if (parseInt(find('[data-test-display-year]').innerText) > LICENSE_START.getFullYear()) {
      await click('[data-test-previous-year]');
    }
    await click(`[data-test-calendar-month=${ARRAY_OF_MONTHS[LICENSE_START.getMonth()]}]`);
    assert
      .dom(SELECTORS.runningTotalMonthStats)
      .doesNotExist('running total single month stat boxes do not show');
    assert
      .dom(SELECTORS.charts.chart('running total'))
      .doesNotExist('running total month over month charts do not show');
    assert.dom(SELECTORS.attributionBlock).exists('attribution area shows');
    assert
      .dom('[data-test-chart-container="new-clients"] [data-test-component="empty-state"]')
      .exists('new client attribution has empty state');
    assert
      .dom('[data-test-empty-state-subtext]')
      .hasText('There are no new clients for this namespace during this time period.    ');
    assert.dom('[data-test-chart-container="total-clients"]').exists('total client attribution chart shows');

    // reset to billing period
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-current-billing-period]');

    // change billing start to month/year of first upgrade
    await click(SELECTORS.counts.startEdit);
    await click(SELECTORS.monthDropdown);
    await click(`[data-test-dropdown-month="${ARRAY_OF_MONTHS[UPGRADE_DATE.getMonth()]}"]`);
    await click(SELECTORS.yearDropdown);
    await click(`[data-test-dropdown-year="${UPGRADE_DATE.getFullYear()}"]`);
    await click('[data-test-date-dropdown-submit]');
    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert
      .dom(SELECTORS.charts.chart('running total'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(SELECTORS.charts.line.xAxisLabel)
      .hasText(`8/22`, 'x-axis labels start with updated billing start month');
    assert.strictEqual(
      findAll('[data-test-line-chart="plot-point"]').length,
      6,
      `line chart plots 6 points to match query`
    );

    // query three months ago (Oct 2022)
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-show-calendar]');
    await click('[data-test-previous-year]');
    await click('[data-test-calendar-month="October"]');

    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert
      .dom(SELECTORS.charts.chart('running total'))
      .exists('Shows running totals with monthly breakdown charts');
    assert.strictEqual(
      findAll('[data-test-line-chart="plot-point"]').length,
      3,
      `line chart plots 3 points to match query`
    );
    const xAxisLabels = findAll(SELECTORS.charts.line.xAxisLabel);
    assert
      .dom(xAxisLabels[xAxisLabels.length - 1])
      .hasText(`10/22`, 'x-axis labels end with queried end month');

    // query for single, historical month (upgrade month)
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-show-calendar]');
    assert.dom('[data-test-display-year]').hasText('2022');
    await click('[data-test-calendar-month="August"]');

    assert.dom(SELECTORS.runningTotalMonthStats).exists('running total single month stat boxes show');
    assert
      .dom(SELECTORS.charts.chart('running total'))
      .doesNotExist('running total month over month charts do not show');
    assert.dom(SELECTORS.attributionBlock).exists('attribution area shows');
    assert.dom('[data-test-chart-container="new-clients"]').exists('new client attribution chart shows');
    assert.dom('[data-test-chart-container="total-clients"]').exists('total client attribution chart shows');

    // reset to billing period
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-current-billing-period]');
    // query month older than count start date
    await click(SELECTORS.counts.startEdit);
    await click(SELECTORS.monthDropdown);
    await click(`[data-test-dropdown-month="${ARRAY_OF_MONTHS[LICENSE_START.getMonth()]}"]`);
    await click(SELECTORS.yearDropdown);
    await click(`[data-test-dropdown-year="${LICENSE_START.getFullYear() - 3}"]`);
    await click('[data-test-date-dropdown-submit]');
    assert
      .dom(SELECTORS.counts.startDiscrepancy)
      .hasTextContaining(
        `We only have data from January 2022`,
        'warning banner displays that date queried was prior to count start date'
      );
  });

  test('totals filter correctly with full data', async function (assert) {
    assert
      .dom(SELECTORS.charts.chart('running total'))
      .exists('Shows running totals with monthly breakdown charts');
    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');

    const response = await this.store.peekRecord('clients/activity', 'some-activity-id');
    // FILTER BY NAMESPACE
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    const topNamespace = response.byNamespace[0];
    const topMount = topNamespace.mounts[0];

    assert.dom(SELECTORS.selectedNs).hasText(topNamespace.label, 'selects top namespace');
    assert.dom('[data-test-top-attribution]').includesText('Top auth method');
    assert
      .dom(SELECTORS.charts.statTextValue('Entity clients'))
      .includesText(`${formatNumber([topNamespace.entity_clients])}`, 'total entity clients is accurate');
    assert
      .dom(SELECTORS.charts.statTextValue('Non-entity clients'))
      .includesText(
        `${formatNumber([topNamespace.non_entity_clients])}`,
        'total non-entity clients is accurate'
      );
    // * unavailable during SYNC BETA (1.16.0), planned for 1.16.1 release
    // assert
    //   .dom(SELECTORS.charts.statTextValue('Secrets sync clients'))
    //   .includesText(`${formatNumber([topNamespace.secret_syncs])}`, 'total sync clients is accurate');
    assert
      .dom('[data-test-attribution-clients] p')
      .includesText(`${formatNumber([topMount.clients])}`, 'top attribution clients accurate');

    // FILTER BY AUTH METHOD
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();

    assert.ok(true, 'Filter by first auth method');
    assert.dom(SELECTORS.selectedAuthMount).hasText(topMount.label, 'selects top mount');
    assert
      .dom(SELECTORS.charts.statTextValue('Entity clients'))
      .includesText(`${formatNumber([topMount.entity_clients])}`, 'total entity clients is accurate');
    assert
      .dom(SELECTORS.charts.statTextValue('Non-entity clients'))
      .includesText(`${formatNumber([topMount.non_entity_clients])}`, 'total non-entity clients is accurate');
    // * unavailable during SYNC BETA (1.16.0), planned for 1.16.1 release
    // assert
    //   .dom(SELECTORS.charts.statTextValue('Secrets sync clients'))
    //   .includesText(`${formatNumber([topMount.secret_syncs])}`, 'total sync clients is accurate');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');

    await click('#namespace-search-select [data-test-selected-list-button="delete"]');
    assert.ok(true, 'Remove namespace filter without first removing auth method filter');
    assert.dom('[data-test-top-attribution]').includesText('Top namespace');
    assert
      .dom(SELECTORS.charts.statTextValue('Entity clients'))
      .hasTextContaining(
        `${formatNumber([response.total.entity_clients])}`,
        'total entity clients is back to unfiltered value'
      );
    assert
      .dom(SELECTORS.charts.statTextValue('Non-entity clients'))
      .hasTextContaining(
        `${formatNumber([formatNumber([response.total.non_entity_clients])])}`,
        'total non-entity clients is back to unfiltered value'
      );
    // * unavailable during SYNC BETA (1.16.0), planned for 1.16.1 release
    // assert
    //   .dom(SELECTORS.charts.statTextValue('Secrets sync clients'))
    //   .hasTextContaining(
    //     `${formatNumber([formatNumber([response.total.secret_syncs])])}`,
    //     'total sync clients is back to unfiltered value'
    //   );
    assert
      .dom('[data-test-attribution-clients]')
      .hasTextContaining(
        `${formatNumber([topNamespace.clients])}`,
        'top attribution clients back to unfiltered value'
      );
  });
});
