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
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { dateFormat } from 'core/helpers/date-format';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);
const { syncTab, charts, usageStats } = CLIENT_COUNT;

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
    // set this to 0
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    this.startTimestamp = START_TIME;
    this.endTimestamp = END_TIME;
    this.isSecretsSyncActivated = true;

    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Sync
        @isSecretsSyncActivated={{this.isSecretsSyncActivated}}
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
    assert.strictEqual(
      dataBars.length,
      this.activity.byMonth.filter((m) => m.clients).length,
      'it renders a bar for each non-zero month'
    );
  });

  test('it should render an empty state for no monthly data', async function (assert) {
    assert.expect(5);
    this.activity.set('byMonth', []);

    await this.renderComponent();

    assert.dom(charts.chart('Secrets sync usage')).doesNotExist('vertical bar chart does not render');
    assert.dom(GENERAL.emptyStateTitle).hasText('No monthly secrets sync clients');
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

  test('it should render an empty state if secrets sync is not activated', async function (assert) {
    this.isSecretsSyncActivated = false;

    await this.renderComponent();

    assert.dom(GENERAL.emptyStateTitle).hasText('No Secrets Sync clients');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('No data is available because Secrets Sync has not been activated.');
    assert.dom(GENERAL.emptyStateActions).hasText('Activate Secrets Sync');

    assert.dom(charts.chart('Secrets sync usage')).doesNotExist();
    assert.dom(syncTab.total).doesNotExist();
    assert.dom(syncTab.average).doesNotExist();
  });

  test('it should render an empty chart if secrets sync is activated but no secrets synced', async function (assert) {
    this.isSecretsSyncActivated = true;
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
    await this.renderComponent();

    assert
      .dom(syncTab.total)
      .hasText(
        'Total sync clients The total number of secrets synced from Vault to other destinations during this date range. 0'
      );
    assert.dom(syncTab.average).doesNotExist('Does not render average if the calculation is 0');
  });
});
