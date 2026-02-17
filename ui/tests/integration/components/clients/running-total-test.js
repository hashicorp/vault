/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { getUnixTime } from 'date-fns';
import { findAll } from '@ember/test-helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import { CLIENT_COUNT, CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import {
  destructureClientCounts,
  formatByMonths,
  formatByNamespace,
} from 'core/utils/client-counts/serializers';

const START_TIME = getUnixTime(LICENSE_START);

module('Integration | Component | clients/running-total', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.flags = this.owner.lookup('service:flags');
    this.version = this.owner.lookup('service:version');
    this.flags.activatedFlags = ['secrets-sync'];
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    const activityResponse = await this.owner
      .lookup('service:api')
      .sys.internalClientActivityReportCounts(undefined, getUnixTime(timestamp.now()), undefined, START_TIME);
    this.activity = {
      ...activityResponse,
      by_namespace: formatByNamespace(activityResponse.by_namespace),
      by_month: formatByMonths(activityResponse.months),
      total: destructureClientCounts(activityResponse.total),
    };
    this.byMonthClients = this.activity.by_month.map((d) => d.new_clients);

    this.renderComponent = async () => {
      await render(hbs`
      <Clients::RunningTotal
        @byMonthClients={{this.byMonthClients}}
        @runningTotals={{this.activity.total}}
      />
    `);
    };
  });

  test('it text for community versions', async function (assert) {
    this.version.type = 'community';
    await this.renderComponent();
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} p`)
      .hasText(
        'Number of clients in the date range by client type, and a breakdown of new clients per month during the date range.'
      );
  });

  // abbreviating "ent" so the test is not filtered out in CE repo runs
  test('it renders text for ent versions', async function (assert) {
    this.version.type = 'enterprise';
    await this.renderComponent();
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} p`)
      .hasText(
        'Number of clients in the billing period by client type, and a breakdown of new clients per month during the billing period.'
      );
  });

  test('it renders text for HVD managed versions', async function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    await this.renderComponent();
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} p`)
      .hasText(
        'Number of total unique clients in the data period by client type, and total number of unique clients per month. The monthly total is the relevant billing metric.'
      );
  });

  test('it renders with full monthly activity data', async function (assert) {
    await this.renderComponent();

    assert.dom(CLIENT_COUNT.card('Client usage trends')).exists('running total component renders');
    assert.dom(CHARTS.chart('Client usage by month')).exists('bar chart renders');
    assert.dom(CHARTS.legend).hasText('New clients');
    const expectedColor = 'rgb(28, 52, 95)';
    const color = getComputedStyle(find(CHARTS.legendDot(1))).backgroundColor;
    assert.strictEqual(color, expectedColor, `actual color: ${color}, expected color: ${expectedColor}`);

    const expectedValues = {
      'Entity clients': formatNumber([this.activity.total.entity_clients]),
      'Non-entity clients': formatNumber([this.activity.total.non_entity_clients]),
      'ACME clients': formatNumber([this.activity.total.acme_clients]),
      'Secret sync clients': formatNumber([this.activity.total.secret_syncs]),
    };
    for (const label in expectedValues) {
      assert
        .dom(CLIENT_COUNT.statLegendValue(label))
        .hasText(
          `${expectedValues[label]} ${label}`,
          `stat label: ${label} renders correct total: ${expectedValues[label]}`
        );
    }

    // assert bar chart is correct
    findAll(CHARTS.xAxisLabel).forEach((e, i) => {
      const timestamp = this.byMonthClients[i].timestamp;
      const displayMonth = parseAPITimestamp(timestamp, 'M/yy');
      assert.dom(e).hasText(displayMonth, `renders x-axis labels for bar chart: ${displayMonth}`);
    });
    assert
      .dom(CHARTS.verticalBar)
      .exists({ count: this.byMonthClients.length }, 'renders correct number of bars ');
  });

  test('it toggles to split chart by client type', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.inputByAttr('toggle view'));

    assert.dom(CLIENT_COUNT.card('Client usage trends')).exists('running total component renders');
    assert.dom(CHARTS.chart('Client usage by month')).exists('bar chart renders');
    assert
      .dom(CHARTS.legend)
      .hasText(
        'Entity clients Non-entity clients ACME clients Secret sync clients',
        'it renders legend in order that matches the stacked bar data and secret sync clients is last'
      );

    // assert each legend item is correct
    const expectedLegend = [
      { label: 'Entity clients', color: 'rgb(66, 105, 208)' },
      { label: 'Non-entity clients', color: 'rgb(239, 177, 23)' },
      { label: 'ACME clients', color: 'rgb(255, 114, 92)' },
      { label: 'Secret sync clients', color: 'rgb(108, 197, 176)' },
    ];

    findAll('.legend-item').forEach((e, i) => {
      const { label, color } = expectedLegend[i];
      assert.dom(e).hasText(label, `legend renders label: ${label}`);
      const dotColor = getComputedStyle(find(CHARTS.legendDot(i + 1))).backgroundColor;
      assert.strictEqual(dotColor, color, `${label} - actual color: ${dotColor}, expected: ${color}`);
    });

    // assert bar chart is correct
    findAll(CHARTS.xAxisLabel).forEach((e, i) => {
      const timestamp = this.byMonthClients[i].timestamp;
      const displayMonth = parseAPITimestamp(timestamp, 'M/yy');
      assert.dom(e).hasText(`${displayMonth}`, `renders x-axis labels for bar chart: ${displayMonth}`);
    });

    const months = this.byMonthClients.length;
    const barsPerMonth = expectedLegend.length;
    assert
      .dom(CHARTS.verticalBar)
      .exists({ count: months * barsPerMonth }, `renders ${barsPerMonth} bars per month`);
  });

  test('it renders when no monthly breakdown is available', async function (assert) {
    this.byMonthClients = [];
    await this.renderComponent();
    const expectedStats = {
      Entity: formatNumber([this.activity.total.entity_clients]),
      'Non-entity': formatNumber([this.activity.total.non_entity_clients]),
      ACME: formatNumber([this.activity.total.acme_clients]),
      'Secret sync': formatNumber([this.activity.total.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(CLIENT_COUNT.statTextValue(label))
        .hasText(
          `${expectedStats[label]}`,
          `stat label: ${label} renders single month new clients: ${expectedStats[label]}`
        );
    }
    assert.dom(CHARTS.chart('Client usage by month')).doesNotExist('bar chart does not render');
    assert.dom(CLIENT_COUNT.statTextValue()).exists({ count: 5 }, 'renders 5 stat text containers');
  });

  test('it hides secret sync totals when feature is not activated', async function (assert) {
    this.flags.activatedFlags = [];
    // reset secret sync clients to 0
    this.byMonthClients = this.byMonthClients.map((obj) => ({ ...obj, secret_syncs: 0 }));

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.card('Client usage trends')).exists('running total component renders');
    assert.dom(CHARTS.chart('Client usage by month')).exists('bar chart renders');
    assert.dom(CLIENT_COUNT.statLegendValue('Entity clients')).exists();
    assert.dom(CLIENT_COUNT.statLegendValue('Non-entity clients')).exists();
    assert
      .dom(CLIENT_COUNT.statLegendValue('Secret sync clients'))
      .doesNotExist('does not render secret syncs');

    // check toggle view
    await click(GENERAL.inputByAttr('toggle view'));
    assert
      .dom(CHARTS.legend)
      .hasText('Entity clients Non-entity clients ACME clients', 'legend does not include sync clients');

    // assert each legend item is correct
    const expectedLegend = [
      { label: 'Entity clients', color: 'rgb(66, 105, 208)' },
      { label: 'Non-entity clients', color: 'rgb(239, 177, 23)' },
      { label: 'ACME clients', color: 'rgb(255, 114, 92)' },
    ];

    findAll('.legend-item').forEach((e, i) => {
      const { label, color } = expectedLegend[i];
      assert.dom(e).hasText(label, `legend renders label: ${label}`);
      const dotColor = getComputedStyle(find(CHARTS.legendDot(i + 1))).backgroundColor;
      assert.strictEqual(dotColor, color, `${label} - actual color: ${dotColor}, expected: ${color}`);
    });

    const months = this.byMonthClients.length;
    const barsPerMonth = expectedLegend.length;
    assert
      .dom(CHARTS.verticalBar)
      .exists({ count: months * barsPerMonth }, `renders ${barsPerMonth} bars per month`);
  });
});
