/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import { getUnixTime } from 'date-fns';
import { SELECTORS } from 'vault/tests/helpers/clients';
import { formatNumber } from 'core/helpers/format-number';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { dateFormat } from 'core/helpers/date-format';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);
const { syncTab, charts, usageStats } = SELECTORS;

module('Integration | Component | clients | Clients::Page::Sync', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: END_TIME },
    };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    this.startTimestamp = START_TIME;
    this.endTimestamp = END_TIME;
    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Sync
        @activity={{this.activity}}
        @versionHistory={{this.versionHistory}}
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @namespace={{this.countsController.ns}}
        @mountPath={{this.countsController.mountPath}}
      />
    `);
  });

  test('it should render with full month activity data', async function (assert) {
    assert.expect(4 + this.activity.byMonth.length);
    const expectedTotal = formatNumber([this.activity.total.secret_syncs]);
    const expectedAvg = formatNumber([calculateAverage(this.activity.byMonth, 'secret_syncs')]);
    await this.renderComponent();
    assert
      .dom(syncTab.total)
      .hasText(
        `Total sync clients The total number of secrets synced from Vault to other destinations during this date range. ${expectedTotal}`,
        `renders correct total sync stat ${expectedTotal}`
      );
    assert
      .dom(syncTab.average)
      .hasText(
        `Average sync clients per month ${expectedAvg}`,
        `renders correct average sync stat ${expectedAvg}`
      );

    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert.dom(charts.timestamp).hasText(`Updated ${formattedTimestamp}`, 'renders response timestamp');

    // assert bar chart is correct
    findAll(`${charts.chart('Secrets sync usage')} ${charts.xAxisLabel}`).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for bar chart: ${this.activity.byMonth[i].month}`
        );
    });

    const dataBars = findAll(charts.dataBar).filter((b) => b.hasAttribute('height'));
    assert.strictEqual(dataBars.length, this.activity.byMonth.filter((m) => m.counts !== null).length);
  });

  test('it should render empty state for no monthly data', async function (assert) {
    assert.expect(5);
    this.activity.set('byMonth', []);

    await this.renderComponent();

    assert.dom(charts.chart('Secrets sync usage')).doesNotExist('vertical bar chart does not render');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No monthly secrets sync clients');
    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert.dom(charts.timestamp).hasText(`Updated ${formattedTimestamp}`, 'renders timestamp');
    assert.dom(syncTab.total).doesNotExist('total sync counts does not exist');
    assert.dom(syncTab.average).doesNotExist('average sync client counts does not exist');
  });

  test('it should render stats without chart for a single month', async function (assert) {
    assert.expect(4);
    const activityQuery = { start_time: { timestamp: START_TIME }, end_time: { timestamp: START_TIME } };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    const total = formatNumber([this.activity.total.secret_syncs]);
    await this.renderComponent();

    assert.dom(charts.chart('Secrets sync usage')).doesNotExist('vertical bar chart does not render');
    assert
      .dom(usageStats)
      .hasText(
        `Secrets sync usage This data can be used to understand how many secrets sync clients have been used for this date range. Each Vault secret that is synced to at least one destination counts as one Vault client. Total sync clients ${total}`,
        'renders sync stats instead of chart'
      );
    assert.dom(syncTab.total).doesNotExist('total sync counts does not exist');
    assert.dom(syncTab.average).doesNotExist('average sync client counts does not exist');
  });
});
