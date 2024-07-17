/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, findAll, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import { getUnixTime } from 'date-fns';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { selectChoose } from 'ember-power-select/test-support';
import timestamp from 'core/utils/timestamp';
import sinon from 'sinon';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

const START_TIME = getUnixTime(LICENSE_START);
const END_TIME = getUnixTime(STATIC_NOW);

module('Integration | Component | clients | Page::Counts', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: END_TIME },
    };
    this.activity = await this.store.queryRecord('clients/activity', activityQuery);
    this.config = await this.store.queryRecord('clients/config', {});
    this.startTimestamp = START_TIME;
    this.endTimestamp = END_TIME;
    this.versionHistory = [];
    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Counts
        @activity={{this.activity}}
        @activityError={{this.activityError}}
        @config={{this.config}}
        @versionHistory={{this.versionHistory}}
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @namespace={{this.namespace}}
        @mountPath={{this.mountPath}}
        @onFilterChange={{this.onFilterChange}}
      >
        <div data-test-yield>Yield block</div>
      </Clients::Page::Counts>
    `);
  });

  test('it should populate start and end month displays', async function (assert) {
    await this.renderComponent();

    assert.dom(CLIENT_COUNT.dateRange.dateDisplay('start')).hasText('July 2023', 'Start month renders');
    assert.dom(CLIENT_COUNT.dateRange.dateDisplay('end')).hasText('January 2024', 'End month renders');
  });

  test('it should render no data empty state', async function (assert) {
    this.activity = { id: 'no-data' };

    await this.renderComponent();

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No data received from July 2023 to January 2024', 'No data empty state renders');
  });

  test('it should render activity error', async function (assert) {
    this.activity = null;
    this.activityError = { httpStatus: 403 };

    await this.renderComponent();

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('You are not authorized', 'Activity error empty state renders');
  });

  test('it should render config disabled alert', async function (assert) {
    this.config.enabled = 'Off';

    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.counts.configDisabled)
      .hasText('Tracking is disabled', 'Config disabled alert renders');
  });

  test('it should send correct values on start and end date change', async function (assert) {
    assert.expect(3);
    const jan23start = getUnixTime(new Date('2023-01-01T00:00:00Z'));
    const dec23end = getUnixTime(new Date('2023-12-31T00:00:00Z'));
    const jan24end = getUnixTime(new Date('2024-01-31T00:00:00Z'));

    const expected = { start_time: START_TIME, end_time: END_TIME };
    this.onFilterChange = (params) => {
      assert.deepEqual(params, expected, 'Correct values sent on filter change');
      this.set('startTimestamp', params.start_time || START_TIME);
      this.set('endTimestamp', params.end_time || END_TIME);
    };
    // page starts with default billing dates, which are july 23 - jan 24
    await this.renderComponent();

    // First, change only the start date
    expected.start_time = jan23start;
    // the end date which is first set to STATIC_NOW gets recalculated
    // to the end of given month/year on date range change
    expected.end_time = jan24end;
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2023-01');
    await click(GENERAL.saveButton);

    // Then change only the end date
    expected.end_time = dec23end;
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('end'), '2023-12');
    await click(GENERAL.saveButton);

    // Then reset to billing which should reset the params
    expected.start_time = undefined;
    expected.end_time = undefined;
    await click(CLIENT_COUNT.dateRange.edit);
    await click(CLIENT_COUNT.dateRange.reset);
    await click(GENERAL.saveButton);
  });

  test('it should render namespace and auth mount filters', async function (assert) {
    assert.expect(5);

    this.namespace = 'root';
    this.mountPath = 'auth/authid0';

    let assertion = (params) =>
      assert.deepEqual(params, { ns: undefined, mountPath: undefined }, 'Auth mount cleared with namespace');
    this.onFilterChange = (params) => {
      if (assertion) {
        assertion(params);
      }
      const keys = Object.keys(params);
      this.namespace = keys.includes('ns') ? params.ns : this.namespace;
      this.mountPath = keys.includes('mountPath') ? params.mountPath : this.mountPath;
    };

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.counts.namespaces).includesText(this.namespace, 'Selected namespace renders');
    assert.dom(CLIENT_COUNT.counts.mountPaths).includesText(this.mountPath, 'Selected auth mount renders');

    await click(`${CLIENT_COUNT.counts.namespaces} button`);
    // this is only necessary in tests since SearchSelect does not respond to initialValue changes
    // in the app the component is rerender on query param change
    assertion = null;
    await click(`${CLIENT_COUNT.counts.mountPaths} button`);
    assertion = (params) => assert.true(params.ns.includes('ns'), 'Namespace value sent on change');
    await selectChoose(CLIENT_COUNT.counts.namespaces, '.ember-power-select-option', 0);

    assertion = (params) =>
      assert.true(params.mountPath.includes('auth/'), 'Auth mount value sent on change');
    await selectChoose(CLIENT_COUNT.counts.mountPaths, 'auth/authid0');
  });

  test('it should render start time discrepancy alert', async function (assert) {
    this.startTimestamp = getUnixTime(new Date('2022-06-01T00:00:00Z'));

    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.counts.startDiscrepancy)
      .hasText(
        'You requested data from June 2022. We only have data from July 2023, and that is what is being shown here.',
        'Start discrepancy alert renders'
      );
  });

  test('it renders alert if upgrade happened within queried activity', async function (assert) {
    assert.expect(5);
    this.versionHistory = await this.store.findAll('clients/version-history').then((resp) => {
      return resp.map(({ version, previousVersion, timestampInstalled }) => {
        return {
          version,
          previousVersion,
          timestampInstalled,
        };
      });
    });

    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.upgradeWarning)
      .hasTextContaining(
        `Client count data contains 3 upgrades Vault was upgraded during this time period. Keep this in mind while looking at the data. Visit our Client count FAQ for more information.`,
        'it renders title and subtext'
      );

    const [first, second, third] = findAll(`${CLIENT_COUNT.upgradeWarning} li`);
    assert
      .dom(first)
      .hasText(
        `1.9.1 (upgraded on Aug 2, 2023) - We introduced changes to non-entity token and local auth mount logic for client counting in 1.9.`,
        'alert includes 1.9.1 upgrade'
      );
    assert
      .dom(second)
      .hasTextContaining(
        `1.10.1 (upgraded on Sep 2, 2023) - We added monthly breakdowns and mount level attribution starting in 1.10.`,
        'alert includes 1.10.1 upgrade'
      );
    assert
      .dom(third)
      .hasTextContaining(
        `1.17.0 (upgraded on Dec 2, 2023) - We separated ACME clients from non-entity clients starting in 1.17.`,
        'alert includes 1.17.0 upgrade'
      );
    assert
      .dom(`${CLIENT_COUNT.upgradeWarning} ul`)
      .doesNotHaveTextContaining(
        '1.10.3',
        'Warning does not include subsequent patch releases (e.g. 1.10.3) of the same notable upgrade.'
      );
  });

  test('it should render empty state for no start or license start time', async function (assert) {
    this.startTimestamp = null;
    this.config.billingStartTimestamp = null;
    this.activity = {};

    await this.renderComponent();

    assert.dom(GENERAL.emptyStateTitle).hasText('No start date found', 'Empty state renders');
    assert.dom(CLIENT_COUNT.dateRange.set).exists();
  });

  test('it should render catch all empty state', async function (assert) {
    this.activity.total = null;

    await this.renderComponent();

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No data received from July 2023 to January 2024', 'Empty state renders');
  });
});
