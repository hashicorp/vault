/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { render, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import { getUnixTime } from 'date-fns';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT, CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';
import { formatNumber } from 'core/helpers/format-number';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { dateFormat } from 'core/helpers/date-format';
import { assertBarChart } from 'vault/tests/helpers/clients/client-count-helpers';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);
const { statText, usageStats } = CLIENT_COUNT;

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
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it should render with full month activity data charts', async function (assert) {
    const monthCount = this.activity.byMonth.length;
    assert.expect(8 + monthCount * 2);
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
    assert.dom(CHARTS.timestamp).hasText(`Updated ${formattedTimestamp}`, 'renders response timestamp');

    assertBarChart(assert, 'ACME usage', this.activity.byMonth);
    assertBarChart(assert, 'Monthly new', this.activity.byMonth);
  });

  test('it should render stats without chart for a single month', async function (assert) {
    assert.expect(5);
    const activityQuery = { start_time: { timestamp: END_TIME }, end_time: { timestamp: END_TIME } };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    const expectedTotal = formatNumber([this.activity.total.acme_clients]);
    await this.renderComponent();

    assert.dom(CHARTS.chart('ACME usage')).doesNotExist('total usage chart does not render');
    assert.dom(CHARTS.container('Monthly new')).doesNotExist('monthly new chart does not render');
    assert.dom(statText('Average ACME clients per month')).doesNotExist();
    assert.dom(statText('Average new ACME clients per month')).doesNotExist();
    assert
      .dom(usageStats('ACME usage'))
      .hasText(
        `ACME usage Usage metrics tutorial This data can be used to understand how many ACME clients have been used for the queried month. Each ACME request is counted as one client. Total ACME clients ${expectedTotal}`,
        'it renders usage stats with single month copy'
      );
  });

  // EMPTY STATES
  test('it should render empty state when ACME data does not exist for a date range', async function (assert) {
    assert.expect(8);
    // this happens when a user queries historical data that predates the monthly breakdown (added in 1.11)
    // only entity + non-entity clients existed then, so we show an empty state for ACME clients
    // because the activity response just returns { acme_clients: 0 } which isn't very clear
    this.activity.byMonth = [];

    await this.renderComponent();

    assert.dom(GENERAL.emptyStateTitle).hasText('No ACME clients');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('There is no ACME client data available for this date range.');

    assert.dom(CHARTS.chart('ACME usage')).doesNotExist('vertical bar chart does not render');
    assert.dom(CHARTS.container('Monthly new')).doesNotExist('monthly new chart does not render');
    assert.dom(statText('Total ACME clients')).doesNotExist();
    assert.dom(statText('Average ACME clients per month')).doesNotExist();
    assert.dom(statText('Average new ACME clients per month')).doesNotExist();
    assert.dom(usageStats('ACME usage')).doesNotExist();
  });

  test('it should render empty state when ACME data does not exist for a single month', async function (assert) {
    assert.expect(1);
    const activityQuery = { start_time: { timestamp: START_TIME }, end_time: { timestamp: START_TIME } };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    this.activity.byMonth = [];

    await this.renderComponent();

    assert.dom(GENERAL.emptyStateMessage).hasText('There is no ACME client data available for this month.');
  });

  test('it should render empty total usage chart when monthly counts are null or 0', async function (assert) {
    assert.expect(9);
    // manually stub because mirage isn't setup to handle mixed data yet
    const counts = {
      acme_clients: 0,
      clients: 19,
      entity_clients: 0,
      non_entity_clients: 19,
      secret_syncs: 0,
    };
    this.activity.byMonth = [
      {
        month: '3/24',
        timestamp: '2024-03-01T00:00:00Z',
        namespaces: [],
        namespaces_by_key: {},
        new_clients: {
          month: '3/24',
          timestamp: '2024-03-01T00:00:00Z',
          namespaces: [],
        },
      },
      {
        month: '4/24',
        timestamp: '2024-04-01T00:00:00Z',
        ...counts,
        namespaces: [],
        namespaces_by_key: {},
        new_clients: {
          month: '4/24',
          timestamp: '2024-04-01T00:00:00Z',
          namespaces: [],
        },
      },
    ];
    this.activity.total = counts;

    await this.renderComponent();

    assert.dom(CHARTS.chart('ACME usage')).exists('renders empty ACME usage chart');
    assert
      .dom(statText('Total ACME clients'))
      .hasTextContaining('The total number of ACME requests made to Vault during this time period. 0');
    findAll(`${CHARTS.chart('ACME usage')} ${CHARTS.xAxisLabel}`).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for empty bar chart: ${this.activity.byMonth[i].month}`
        );
    });
    findAll(`${CHARTS.chart('ACME usage')} ${CHARTS.verticalBar}`).forEach((e, i) => {
      assert.dom(e).isNotVisible(`does not render data bar for: ${this.activity.byMonth[i].month}`);
    });

    assert
      .dom(CHARTS.container('Monthly new'))
      .doesNotExist('empty monthly new chart does not render at all');
    assert.dom(statText('Average ACME clients per month')).doesNotExist();
    assert.dom(statText('Average new ACME clients per month')).doesNotExist();
  });
});
