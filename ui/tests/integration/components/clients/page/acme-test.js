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
import { addMonths, getUnixTime } from 'date-fns';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { formatNumber } from 'core/helpers/format-number';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { dateFormat } from 'core/helpers/date-format';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);
const { statText, charts, usageStats } = CLIENT_COUNT;

module('Integration | Component | clients | Clients::Page::Acme', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: END_TIME },
    };
    // set this to 0
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    this.startTimestamp = START_TIME;
    this.endTimestamp = END_TIME;
    this.isSecretsSyncActivated = true;

    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Acme
        @activity={{this.activity}}
        @versionHistory={{this.versionHistory}}
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @namespace={{this.countsController.ns}}
        @mountPath={{this.countsController.mountPath}}
      />
    `);
  });

  test('it should render with full month activity data charts', async function (assert) {
    assert.expect(6 + this.activity.byMonth.length);
    const expectedTotal = formatNumber([this.activity.total.acme_clients]);
    const expectedAvg = formatNumber([calculateAverage(this.activity.byMonth, 'acme_clients')]);
    const expectedNewAvg = formatNumber([
      calculateAverage(
        this.activity.byMonth.map((m) => m?.new_clients),
        'acme_clients'
      ),
    ]);
    await this.renderComponent();
    assert
      .dom(statText('Total ACME clients'))
      .hasText(
        `Total ACME clients The total number of ACME requests made to Vault during this time period. ${expectedTotal}`,
        `renders correct total acme stat ${expectedTotal}`
      );
    assert.dom(statText('Average ACME clients per month')).hasTextContaining(`${expectedAvg}`);
    assert.dom(statText('Average new ACME clients per month')).hasTextContaining(`${expectedNewAvg}`);

    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert.dom(charts.timestamp).hasText(`Updated ${formattedTimestamp}`, 'renders response timestamp');

    // assert bar chart is correct
    findAll(`${charts.chart('ACME usage')} ${charts.xAxisLabel}`).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for bar chart: ${this.activity.byMonth[i].month}`
        );
    });

    const totalUsageBars = findAll(`${charts.chart('ACME usage')} ${charts.dataBar}`).filter((b) =>
      b.hasAttribute('height')
    );
    assert.strictEqual(
      totalUsageBars.length,
      this.activity.byMonth.filter((m) => m.clients).length,
      'it renders a bar for each non-zero month in total acme usage chart'
    );

    const monthlyNewBars = findAll(`${charts.chart('Monthly new ACME clients')} ${charts.dataBar}`).filter(
      (b) => b.hasAttribute('height')
    );
    assert.strictEqual(
      monthlyNewBars.length,
      this.activity.byMonth.filter((m) => m.clients).length,
      'it renders a bar for each non-zero month in monthly new acme chart'
    );
  });

  test('it should render usage stats for no monthly data', async function (assert) {
    assert.expect(5);
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: getUnixTime(addMonths(LICENSE_START, 1)) },
    };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    const expectedTotal = formatNumber([this.activity.total.acme_clients]);
    await this.renderComponent();

    assert.dom(charts.chart('ACME usage')).doesNotExist('vertical bar chart does not render');
    assert.dom(charts.chart('Monthly new ACME clients')).doesNotExist('monthly new chart does not render');
    assert.dom(statText('Average ACME clients per month')).doesNotExist();
    assert.dom(statText('Average new ACME clients per month')).doesNotExist();
    assert
      .dom(usageStats)
      .hasText(
        `ACME usage This data can be used to understand how many ACME clients have been used for the queried date range. Each ACME request is counted as one client. Total ACME clients ${expectedTotal}`,
        'it renders usage stats with date range copy'
      );
  });

  test('it should render stats without chart for a single month', async function (assert) {
    assert.expect(5);
    const activityQuery = { start_time: { timestamp: END_TIME }, end_time: { timestamp: END_TIME } };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    const expectedTotal = formatNumber([this.activity.total.acme_clients]);
    await this.renderComponent();

    assert.dom(charts.chart('ACME usage')).doesNotExist('total usage chart does not render');
    assert.dom(charts.chart('Monthly new ACME clients')).doesNotExist('monthly new chart does not render');
    assert.dom(statText('Average ACME clients per month')).doesNotExist();
    assert.dom(statText('Average new ACME clients per month')).doesNotExist();
    assert
      .dom(usageStats)
      .hasText(
        `ACME usage This data can be used to understand how many ACME clients have been used for the queried month. Each ACME request is counted as one client. Total ACME clients ${expectedTotal}`,
        'it renders usage stats with single month copy'
      );
  });
});
