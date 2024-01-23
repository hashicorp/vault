/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import clientsHandler from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { formatRFC3339, getUnixTime } from 'date-fns';
import { findAll } from '@ember/test-helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const START_TIME = getUnixTime(new Date('2023-10-01T00:00:00Z'));

module('Integration | Component | clients/running-total', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => new Date('2024-01-31T23:59:59Z'));
  });

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    const store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: getUnixTime(timestamp.now()) },
    };
    this.activity = await store.queryRecord('clients/activity', activityQuery);
    this.newActivity = this.activity.byMonth.map((d) => d.new_clients);
    this.totalUsageCounts = this.activity.total;
    this.set('timestamp', formatRFC3339(timestamp.now()));
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it renders with full monthly activity data', async function (assert) {
    const expectedTotalEntity = formatNumber([this.totalUsageCounts.entity_clients]);
    const expectedTotalNonEntity = formatNumber([this.totalUsageCounts.non_entity_clients]);
    const expectedTotalSync = formatNumber([this.totalUsageCounts.secret_syncs]);

    await render(hbs`
      <Clients::RunningTotal
        @selectedAuthMethod={{this.selectedAuthMethod}}
        @byMonthActivityData={{this.activity.byMonth}}
        @runningTotals={{this.totalUsageCounts}}
        @upgradeData={{this.upgradeDuringActivity}}
        @responseTimestamp={{this.timestamp}}
        @isHistoricalMonth={{false}}
      />
    `);

    assert.dom('[data-test-running-total]').exists('running total component renders');
    assert.dom('[data-test-line-chart]').exists('line chart renders');
    assert
      .dom('[data-test-running-total-entity] p.data-details')
      .hasText(`${expectedTotalEntity}`, `renders correct total entity average ${expectedTotalEntity}`);
    assert
      .dom('[data-test-running-total-nonentity] p.data-details')
      .hasText(
        `${expectedTotalNonEntity}`,
        `renders correct total nonentity average ${expectedTotalNonEntity}`
      );
    assert
      .dom('[data-test-running-total-sync] p.data-details')
      .hasText(`${expectedTotalSync}`, `renders correct total sync ${expectedTotalSync}`);

    // assert line chart is correct
    findAll('[data-test-line-chart="x-axis-labels"] text').forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for line chart: ${this.activity.byMonth[i].month}`
        );
    });
    assert
      .dom('[data-test-line-chart="plot-point"]')
      .exists(
        { count: this.activity.byMonth.filter((m) => m.counts !== null).length },
        'renders correct number of plot points'
      );
  });

  test('it renders with no new monthly data', async function (assert) {
    this.set(
      'monthlyWithoutNew',
      this.activity.byMonth.map((d) => ({
        ...d,
        new_clients: { month: d.month },
      }))
    );
    const expectedTotalEntity = formatNumber([this.totalUsageCounts.entity_clients]);
    const expectedTotalNonEntity = formatNumber([this.totalUsageCounts.non_entity_clients]);
    const expectedTotalSync = formatNumber([this.totalUsageCounts.secret_syncs]);

    await render(hbs`
      <Clients::RunningTotal
        @selectedAuthMethod={{this.selectedAuthMethod}}
        @byMonthActivityData={{this.monthlyWithoutNew}}
        @runningTotals={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @isHistoricalMonth={{false}}
      />
    `);
    assert.dom('[data-test-running-total]').exists('running total component renders');
    assert.dom('[data-test-line-chart]').exists('line chart renders');

    assert
      .dom('[data-test-running-total-entity] p.data-details')
      .hasText(`${expectedTotalEntity}`, `renders correct total entity average ${expectedTotalEntity}`);
    assert
      .dom('[data-test-running-total-nonentity] p.data-details')
      .hasText(
        `${expectedTotalNonEntity}`,
        `renders correct total nonentity average ${expectedTotalNonEntity}`
      );
    assert
      .dom('[data-test-running-total-sync] p.data-details')
      .hasText(`${expectedTotalSync}`, `renders correct total sync ${expectedTotalSync}`);
  });

  test('it renders with single historical month data', async function (assert) {
    const singleMonth = this.activity.byMonth[this.activity.byMonth.length - 1];
    const singleMonthNew = this.newActivity[this.newActivity.length - 1];
    this.set('singleMonth', [singleMonth]);
    const expectedTotalClients = formatNumber([singleMonth.clients]);
    const expectedTotalEntity = formatNumber([singleMonth.entity_clients]);
    const expectedTotalNonEntity = formatNumber([singleMonth.non_entity_clients]);
    const expectedTotalSync = formatNumber([singleMonth.secret_syncs]);
    const expectedNewClients = formatNumber([singleMonthNew.clients]);
    const expectedNewEntity = formatNumber([singleMonthNew.entity_clients]);
    const expectedNewNonEntity = formatNumber([singleMonthNew.non_entity_clients]);
    const expectedNewSyncs = formatNumber([singleMonthNew.secret_syncs]);

    await render(hbs`
      <Clients::RunningTotal
        @selectedAuthMethod={{this.selectedAuthMethod}}
        @byMonthActivityData={{this.singleMonth}}
        @runningTotals={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @isHistoricalMonth={{true}}
      />
    `);
    assert.dom('[data-test-running-total]').exists('running total component renders');
    assert.dom('[data-test-line-chart]').doesNotExist('line chart does not render');
    assert.dom('[data-test-stat-text-container]').exists({ count: 8 }, 'renders stat text containers');
    assert
      .dom('[data-test-new] [data-test-stat-text-container="New clients"] div.stat-value')
      .hasText(`${expectedNewClients}`, `renders correct total new clients: ${expectedNewClients}`);
    assert
      .dom('[data-test-new] [data-test-stat-text-container="Entity clients"] div.stat-value')
      .hasText(`${expectedNewEntity}`, `renders correct total new entity: ${expectedNewEntity}`);
    assert
      .dom('[data-test-new] [data-test-stat-text-container="Non-entity clients"] div.stat-value')
      .hasText(`${expectedNewNonEntity}`, `renders correct total new non-entity: ${expectedNewNonEntity}`);
    assert
      .dom('[data-test-new] [data-test-stat-text-container="Secrets sync clients"] div.stat-value')
      .hasText(`${expectedNewSyncs}`, `renders correct total new non-entity: ${expectedNewSyncs}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Total monthly clients"] div.stat-value')
      .hasText(`${expectedTotalClients}`, `renders correct total clients: ${expectedTotalClients}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Entity clients"] div.stat-value')
      .hasText(`${expectedTotalEntity}`, `renders correct total entity: ${expectedTotalEntity}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Non-entity clients"] div.stat-value')
      .hasText(`${expectedTotalNonEntity}`, `renders correct total non-entity: ${expectedTotalNonEntity}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Secrets sync clients"] div.stat-value')
      .hasText(`${expectedTotalSync}`, `renders correct total sync: ${expectedTotalSync}`);
  });
});
