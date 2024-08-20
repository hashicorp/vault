/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import { getUnixTime } from 'date-fns';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { dateFormat } from 'core/helpers/date-format';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CHARTS, CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { assertBarChart } from 'vault/tests/helpers/clients/client-count-helpers';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);

module('Integration | Component | clients | Clients::Page::Token', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    const store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: END_TIME },
    };
    this.activity = await store.queryRecord('clients/activity', activityQuery);
    this.newActivity = this.activity.byMonth.map((d) => d.new_clients);
    this.versionHistory = await store
      .findAll('clients/version-history')
      .then((response) => {
        return response.map(({ version, previousVersion, timestampInstalled }) => {
          return {
            version,
            previousVersion,
            timestampInstalled,
          };
        });
      })
      .catch(() => []);
    this.startTimestamp = START_TIME;
    this.endTimestamp = END_TIME;
    this.renderComponent = () =>
      render(hbs`
        <Clients::Page::Token
          @activity={{this.activity}}
          @versionHistory={{this.versionHistory}}
          @startTimestamp={{this.startTimestamp}}
          @endTimestamp={{this.endTimestamp}}
          @namespace={{this.ns}}
          @mountPath={{this.mountPath}}
        />
      `);
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it should render monthly total chart', async function (assert) {
    const count = this.activity.byMonth.length;
    const { entity_clients, non_entity_clients } = this.activity.total;
    assert.expect(count + 8);
    const getAverage = (data) => {
      const average = ['entity_clients', 'non_entity_clients'].reduce((count, key) => {
        return (count += calculateAverage(data, key) || 0);
      }, 0);
      return formatNumber([average]);
    };
    const expectedAvg = getAverage(this.activity.byMonth);
    const expectedTotal = formatNumber([entity_clients + non_entity_clients]);
    const chart = CHARTS.container('Entity/Non-entity clients usage');
    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.statTextValue('Total clients'))
      .hasText(expectedTotal, 'renders correct total clients');
    assert
      .dom(CLIENT_COUNT.statTextValue('Average total clients per month'))
      .hasText(expectedAvg, 'renders correct average clients');

    // assert bar chart is correct
    assert.dom(`${chart} ${CHARTS.xAxis}`).hasText('7/23 8/23 9/23 10/23 11/23 12/23 1/24');
    assertBarChart(assert, 'Entity/Non-entity clients usage', this.activity.byMonth, true);

    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert.dom(`${chart} ${CHARTS.timestamp}`).hasText(`Updated ${formattedTimestamp}`, 'renders timestamp');
    assert.dom(`${chart} ${CHARTS.legendLabel(1)}`).hasText('Entity clients', 'Legend label renders');
    assert.dom(`${chart} ${CHARTS.legendLabel(2)}`).hasText('Non-entity clients', 'Legend label renders');
  });

  test('it should render monthly new chart', async function (assert) {
    const count = this.newActivity.length;
    assert.expect(count + 8);
    const expectedNewEntity = formatNumber([calculateAverage(this.newActivity, 'entity_clients')]);
    const expectedNewNonEntity = formatNumber([calculateAverage(this.newActivity, 'non_entity_clients')]);
    const chart = CHARTS.container('Monthly new');

    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.statTextValue('Average new entity clients per month'))
      .hasText(expectedNewEntity, 'renders correct new entity clients');
    assert
      .dom(CLIENT_COUNT.statTextValue('Average new non-entity clients per month'))
      .hasText(expectedNewNonEntity, 'renders correct new nonentity clients');
    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert.dom(`${chart} ${CHARTS.timestamp}`).hasText(`Updated ${formattedTimestamp}`, 'renders timestamp');
    assert.dom(`${chart} ${CHARTS.legendLabel(1)}`).hasText('Entity clients', 'Legend label renders');
    assert.dom(`${chart} ${CHARTS.legendLabel(2)}`).hasText('Non-entity clients', 'Legend label renders');

    // assert bar chart is correct
    assert.dom(`${chart} ${CHARTS.xAxis}`).hasText('7/23 8/23 9/23 10/23 11/23 12/23 1/24');
    assertBarChart(assert, 'Monthly new', this.newActivity, true);
  });

  test('it should render empty state for no new monthly data', async function (assert) {
    this.activity.byMonth = this.activity.byMonth.map((d) => ({
      ...d,
      new_clients: { month: d.month },
    }));
    const chart = CHARTS.container('Monthly new');

    await this.renderComponent();

    assert.dom(`${chart} ${CHARTS.verticalBar}`).doesNotExist('Chart does not render');
    assert.dom(`${chart} ${CHARTS.legend}`).doesNotExist('Legend does not render');
    assert.dom(GENERAL.emptyStateTitle).hasText('No new clients');
    assert
      .dom(CLIENT_COUNT.statText('Average new entity clients per month'))
      .doesNotExist('New client counts does not exist');
    assert
      .dom(CLIENT_COUNT.statText('Average new non-entity clients per month'))
      .doesNotExist('Average new client counts does not exist');
  });

  test('it should render usage stats', async function (assert) {
    assert.expect(6);

    this.activity.endTime = this.activity.startTime;
    const {
      total: { entity_clients, non_entity_clients },
    } = this.activity;

    const checkUsage = () => {
      assert
        .dom(CLIENT_COUNT.statTextValue('Total clients'))
        .hasText(formatNumber([entity_clients + non_entity_clients]), 'Total clients value renders');
      assert
        .dom(CLIENT_COUNT.statTextValue('Entity'))
        .hasText(formatNumber([entity_clients]), 'Entity value renders');
      assert
        .dom(CLIENT_COUNT.statTextValue('Non-entity'))
        .hasText(formatNumber([non_entity_clients]), 'Non-entity value renders');
    };

    // total usage should display for single month query
    await this.renderComponent();
    checkUsage();

    // total usage should display when there is no monthly data
    this.activity.byMonth = null;
    await this.renderComponent();
    checkUsage();
  });
});
