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
import { calculateAverage } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { dateFormat } from 'core/helpers/date-format';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT as ts } from 'vault/tests/helpers/clients/client-count-helpers';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);

module('Integration | Component | clients | Page::Token', function (hooks) {
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
  });

  test('it should render monthly total chart', async function (assert) {
    const getAverage = (data) => {
      const average = ['entity_clients', 'non_entity_clients'].reduce((count, key) => {
        return (count += calculateAverage(data, key) || 0);
      }, 0);
      return formatNumber([average]);
    };
    const expectedTotal = getAverage(this.activity.byMonth);
    const expectedNew = getAverage(this.newActivity);
    const chart = ts.charts.chart('monthly total');

    await this.renderComponent();

    assert
      .dom(ts.charts.statTextValue('Average total clients per month'))
      .hasText(expectedTotal, 'renders correct total clients');
    assert
      .dom(ts.charts.statTextValue('Average new clients per month'))
      .hasText(expectedNew, 'renders correct new clients');
    // assert bar chart is correct
    findAll(`${chart} ${ts.charts.bar.xAxisLabel}`).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for bar chart: ${this.activity.byMonth[i].month}`
        );
    });
    assert
      .dom(`${chart} ${ts.charts.bar.dataBar}`)
      .exists(
        { count: this.activity.byMonth.filter((m) => m.counts !== null).length * 2 },
        'renders correct number of data bars'
      );
    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert
      .dom(`${chart} ${ts.charts.timestamp}`)
      .hasText(`Updated ${formattedTimestamp}`, 'renders timestamp');
    assert.dom(`${chart} ${ts.charts.legendLabel(1)}`).hasText('Entity clients', 'Legend label renders');
    assert.dom(`${chart} ${ts.charts.legendLabel(2)}`).hasText('Non-entity clients', 'Legend label renders');
  });

  test('it should render monthly new chart', async function (assert) {
    const expectedNewEntity = formatNumber([calculateAverage(this.newActivity, 'entity_clients')]);
    const expectedNewNonEntity = formatNumber([calculateAverage(this.newActivity, 'non_entity_clients')]);
    const chart = ts.charts.chart('monthly new');

    await this.renderComponent();

    assert
      .dom(ts.charts.statTextValue('Average new entity clients per month'))
      .hasText(expectedNewEntity, 'renders correct new entity clients');
    assert
      .dom(ts.charts.statTextValue('Average new non-entity clients per month'))
      .hasText(expectedNewNonEntity, 'renders correct new nonentity clients');
    // assert bar chart is correct
    findAll(`${chart} ${ts.charts.bar.xAxisLabel}`).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.activity.byMonth[i].month}`,
          `renders x-axis labels for bar chart: ${this.activity.byMonth[i].month}`
        );
    });
    assert
      .dom(`${chart} ${ts.charts.bar.dataBar}`)
      .exists(
        { count: this.activity.byMonth.filter((m) => m.counts !== null).length * 2 },
        'renders correct number of data bars'
      );
    const formattedTimestamp = dateFormat([this.activity.responseTimestamp, 'MMM d yyyy, h:mm:ss aaa'], {
      withTimeZone: true,
    });
    assert
      .dom(`${chart} ${ts.charts.timestamp}`)
      .hasText(`Updated ${formattedTimestamp}`, 'renders timestamp');
    assert.dom(`${chart} ${ts.charts.legendLabel(1)}`).hasText('Entity clients', 'Legend label renders');
    assert.dom(`${chart} ${ts.charts.legendLabel(2)}`).hasText('Non-entity clients', 'Legend label renders');
  });

  test('it should render empty state for no new monthly data', async function (assert) {
    this.activity.byMonth = this.activity.byMonth.map((d) => ({
      ...d,
      new_clients: { month: d.month },
    }));
    const chart = ts.charts.chart('monthly-new');

    await this.renderComponent();

    assert.dom(`${chart} ${ts.charts.verticalBar}`).doesNotExist('Chart does not render');
    assert.dom(`${chart} ${ts.charts.legend}`).doesNotExist('Legend does not render');
    assert.dom(GENERAL.emptyStateTitle).hasText('No new clients');
    assert.dom(ts.tokenTab.entity).doesNotExist('New client counts does not exist');
    assert.dom(ts.tokenTab.nonentity).doesNotExist('Average new client counts does not exist');
  });

  test('it should render usage stats', async function (assert) {
    assert.expect(6);

    this.activity.endTime = this.activity.startTime;
    const {
      total: { entity_clients, non_entity_clients },
    } = this.activity;

    const checkUsage = () => {
      assert
        .dom(ts.charts.statTextValue('Total clients'))
        .hasText(formatNumber([entity_clients + non_entity_clients]), 'Total clients value renders');
      assert
        .dom(ts.charts.statTextValue('Entity clients'))
        .hasText(formatNumber([entity_clients]), 'Entity clients value renders');
      assert
        .dom(ts.charts.statTextValue('Non-entity clients'))
        .hasText(formatNumber([non_entity_clients]), 'Non-entity clients value renders');
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
