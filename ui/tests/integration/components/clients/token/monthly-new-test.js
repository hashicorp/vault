/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { formatRFC3339, getUnixTime } from 'date-fns';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { SELECTORS } from 'vault/tests/helpers/clients';
import clientsHandler from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import { dateFormat } from 'core/helpers/date-format';

const { tokenTab, charts } = SELECTORS;
const START_TIME = getUnixTime(new Date('2023-10-01T00:00:00Z'));

module('Integration | Component | clients/token/monthly-new', function (hooks) {
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
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it should render with full month activity data', async function (assert) {
    const expectedNewEntity = formatNumber([calculateAverage(this.newActivity, 'entity_clients')]);
    const expectedNewNonEntity = formatNumber([calculateAverage(this.newActivity, 'non_entity_clients')]);

    await render(hbs`
      <Clients::Token::MonthlyNew
        @byMonthActivityData={{this.activity.byMonth}}
        @mountPath={{this.selectedAuthMethod}}
        @runningTotals={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
      />
    `);

    assert
      .dom(tokenTab.entity)
      .hasText(
        `Average new entity clients per month ${expectedNewEntity}`,
        `renders correct new entity stat ${expectedNewEntity}`
      );
    assert
      .dom(tokenTab.nonentity)
      .hasText(
        `Average new non-entity clients per month ${expectedNewNonEntity}`,
        `renders correct new nonentity stat ${expectedNewNonEntity}`
      );
    // assert bar chart is correct
    findAll('[data-test-vertical-chart="x-axis-labels"] text').forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for bar chart: ${this.activity.byMonth[i].month}`
        );
    });
    assert
      .dom('[data-test-vertical-chart="data-bar"]')
      .exists(
        { count: this.activity.byMonth.filter((m) => m.counts !== null).length * 2 },
        'renders correct number of data bars'
      );
  });

  test('it should render empty state for no new monthly data', async function (assert) {
    this.monthlyWithoutNew = this.activity.byMonth.map((d) => ({
      ...d,
      new_clients: { month: d.month },
    }));

    await render(hbs`
      <Clients::Token::MonthlyNew
        @byMonthActivityData={{this.monthlyWithoutNew}}
        @mountPath={{this.selectedAuthMethod}}
        @runningTotals={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
      />
    `);

    assert.dom(charts.verticalBar).doesNotExist('vertical bar chart does not render');
    assert.dom(tokenTab.legend).doesNotExist('legend does not render');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No new clients');
    const formattedTimestamp = dateFormat([this.timestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert.dom(charts.timestamp).hasText(`Updated ${formattedTimestamp}`, 'renders timestamp');
    assert.dom(tokenTab.entity).doesNotExist('new client counts does not exist');
    assert.dom(tokenTab.nonentity).doesNotExist('average new client counts does not exist');
  });
});
