/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import syncHandler from 'vault/mirage/handlers/sync';
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

module('Integration | Component | clients | Clients::Page::Sync', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
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

  module('with secrets sync not activated', function () {
    test('it should render an empty state', async function (assert) {
      await this.renderComponent();

      assert.dom(GENERAL.emptyStateTitle).hasText('No Secrets Sync clients');
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText('No data is available because Secrets Sync has not been activated.');
      assert.dom(GENERAL.emptyStateActions).hasText('Activate Secrets Sync');

      assert.dom(CHARTS.chart('Secrets sync usage')).doesNotExist();
      assert.dom(statText('Total sync clients')).doesNotExist();
      assert.dom(statText('Average sync clients per month')).doesNotExist();
    });
  });

  module('with secrets sync activated', function (hooks) {
    hooks.beforeEach(async function () {
      syncHandler(this.server);
      this.owner.lookup('service:flags').activatedFlags = ['secrets-sync'];

      this.store = this.owner.lookup('service:store');
      const activityQuery = {
        start_time: { timestamp: START_TIME },
        end_time: { timestamp: END_TIME },
      };
      // set this to 0
      this.activity = await this.store.queryRecord('clients/activity', activityQuery);
      this.startTimestamp = START_TIME;
      this.endTimestamp = END_TIME;
    });

    test('it should render with full month activity data', async function (assert) {
      const monthCount = this.activity.byMonth.length;
      assert.expect(8 + monthCount * 2);
      const expectedTotal = formatNumber([this.activity.total.secret_syncs]);
      const expectedAvg = formatNumber([calculateAverage(this.activity.byMonth, 'secret_syncs')]);
      const expectedNewAvg = formatNumber([
        calculateAverage(
          this.activity.byMonth.map((m) => m?.new_clients),
          'secret_syncs'
        ),
      ]);
      await this.renderComponent();

      assert
        .dom(statText('Total sync clients'))
        .hasText(
          `Total sync clients The total number of secrets synced from Vault to other destinations during this date range. ${expectedTotal}`,
          `renders correct total sync stat ${expectedTotal}`
        );
      assert
        .dom(statText('Average sync clients per month'))
        .hasText(
          `Average sync clients per month ${expectedAvg}`,
          `renders correct average sync stat ${expectedAvg}`
        );
      assert.dom(statText('Average new sync clients per month')).hasTextContaining(`${expectedNewAvg}`);

      const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
        withTimeZone: true,
      });
      assert.dom(CHARTS.timestamp).hasText(`Updated ${formattedTimestamp}`, 'renders response timestamp');

      assertBarChart(assert, 'Secrets sync usage', this.activity.byMonth);
      assertBarChart(assert, 'Monthly new', this.activity.byMonth);
    });

    test('it should render stats without chart for a single month', async function (assert) {
      assert.expect(5);
      const activityQuery = { start_time: { timestamp: END_TIME }, end_time: { timestamp: END_TIME } };
      this.activity = await this.store.queryRecord('clients/activity', activityQuery);
      const expectedTotal = formatNumber([this.activity.total.secret_syncs]);
      await this.renderComponent();

      assert.dom(CHARTS.chart('Secrets sync usage')).doesNotExist('total usage chart does not render');
      assert.dom(CHARTS.container('Monthly new')).doesNotExist('monthly new chart does not render');
      assert.dom(statText('Average sync clients per month')).doesNotExist();
      assert.dom(statText('Average new sync clients per month')).doesNotExist();
      assert
        .dom(usageStats('Secrets sync usage'))
        .hasText(
          `Secrets sync usage Usage metrics tutorial This data can be used to understand how many secrets sync clients have been used for this date range. Each Vault secret that is synced to at least one destination counts as one Vault client. Total sync clients ${expectedTotal}`,
          'it renders usage stats with single month copy'
        );
    });

    // EMPTY STATES
    test('it should render empty state when sync data does not exist for a date range', async function (assert) {
      assert.expect(8);
      // this happens when a user queries historical data that predates the monthly breakdown (added in 1.11)
      // only entity + non-entity clients existed then, so we show an empty state for sync clients
      // because the activity response just returns { secret_syncs: 0 } which isn't very clear
      this.activity.byMonth = [];

      await this.renderComponent();

      assert.dom(GENERAL.emptyStateTitle).hasText('No secrets sync clients');
      assert.dom(GENERAL.emptyStateMessage).hasText('There is no sync data available for this date range.');

      assert.dom(CHARTS.chart('Secrets sync usage')).doesNotExist('vertical bar chart does not render');
      assert.dom(CHARTS.container('Monthly new')).doesNotExist('monthly new chart does not render');
      assert.dom(statText('Total sync clients')).doesNotExist();
      assert.dom(statText('Average sync clients per month')).doesNotExist();
      assert.dom(statText('Average new sync clients per month')).doesNotExist();
      assert.dom(usageStats('Secrets sync usage')).doesNotExist();
    });

    test('it should render empty state when sync data does not exist for a single month', async function (assert) {
      assert.expect(1);
      const activityQuery = { start_time: { timestamp: START_TIME }, end_time: { timestamp: START_TIME } };
      this.activity = await this.store.queryRecord('clients/activity', activityQuery);
      this.activity.byMonth = [];
      await this.renderComponent();

      assert.dom(GENERAL.emptyStateMessage).hasText('There is no sync data available for this month.');
    });

    test('it should render an empty total usage chart  if secrets sync is activated but monthly syncs are null or 0', async function (assert) {
      // manually stub because mirage isn't setup to handle mixed data yet
      const counts = {
        clients: 10,
        entity_clients: 4,
        non_entity_clients: 6,
        secret_syncs: 0,
      };
      const monthData = {
        month: '1/24',
        timestamp: '2024-01-01T00:00:00-08:00',
        ...counts,
        namespaces: [
          {
            label: 'root',
            ...counts,
            mounts: [],
          },
        ],
      };
      this.activity.byMonth = [
        {
          ...monthData,
          namespaces_by_key: {
            root: {
              ...monthData,
              mounts_by_key: {},
            },
          },
          new_clients: {
            ...monthData,
          },
        },
      ];
      this.activity.total = counts;
      const monthCount = this.activity.byMonth.length;
      assert.expect(6 + monthCount * 2);
      await this.renderComponent();

      assert.dom(CHARTS.chart('Secrets sync usage')).exists('renders empty sync usage chart');
      assert
        .dom(statText('Total sync clients'))
        .hasText(
          'Total sync clients The total number of secrets synced from Vault to other destinations during this date range. 0'
        );
      assert
        .dom(statText('Average sync clients per month'))
        .doesNotExist('Does not render average if the calculation is 0');
      findAll(`${CHARTS.chart('Secrets sync usage')} ${CHARTS.xAxisLabel}`).forEach((e, i) => {
        assert
          .dom(e)
          .hasText(
            `${this.activity.byMonth[i].month}`,
            `renders x-axis labels for empty bar chart: ${this.activity.byMonth[i].month}`
          );
      });
      findAll(`${CHARTS.chart('Secrets sync usage')} ${CHARTS.verticalBar}`).forEach((e, i) => {
        assert.dom(e).isNotVisible(`does not render data bar for: ${this.activity.byMonth[i].month}`);
      });

      assert
        .dom(CHARTS.container('Monthly new'))
        .doesNotExist('empty monthly new chart does not render at all');
      assert.dom(statText('Average sync clients per month')).doesNotExist();
      assert.dom(statText('Average new sync clients per month')).doesNotExist();
    });
  });
});
