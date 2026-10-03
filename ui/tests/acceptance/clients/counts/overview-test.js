/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, {
  STATIC_NOW,
  LICENSE_START,
  UPGRADE_DATE,
  STATIC_PREVIOUS_MONTH,
} from 'vault/mirage/handlers/clients';
import syncHandler from 'vault/mirage/handlers/sync';
import sinon from 'sinon';
import { visit, click, findAll, fillIn, currentURL } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CHARTS, CLIENT_COUNT, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import timestamp from 'core/utils/timestamp';
import {
  ACTIVITY_EXPORT_STUB,
  ACTIVITY_RESPONSE_STUB,
} from 'vault/tests/helpers/clients/client-count-helpers';
import { ClientFilters, flattenMounts } from 'core/utils/client-counts/helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';

module('Acceptance | clients | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    // These tests use the clientsHandler which dynamically generates activity data, used for asserting date querying, etc
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
  });

  test('it should hide secrets sync stats when feature is NOT on license', async function (assert) {
    // mocks endpoint for no additional license modules
    this.server.get('/sys/license/features', () => ({ features: [] }));

    await login();
    await visit('/vault/clients/counts/overview');
    assert.dom(CLIENT_COUNT.statLegendValue('Secret sync clients')).doesNotExist();
    assert.dom(CLIENT_COUNT.statLegendValue('Entity clients')).exists('other stats are still visible');

    await click(GENERAL.inputByAttr('toggle view'));
    assert.dom(CHARTS.legend).hasText('Entity clients Non-entity clients ACME clients');
  });

  test('it should render charts', async function (assert) {
    await login();
    await visit('/vault/clients/counts/overview');
    assert
      .dom(`${GENERAL.flashMessage}.is-info`)
      .includesText(
        'counts returned in this usage period are an estimate',
        'Shows warning from API about client count estimations'
      );
    assert
      .dom(CLIENT_COUNT.dateRange.dateDisplay('start'))
      .hasText('July 2023', 'start month is correctly parsed from license');
    assert
      .dom(CLIENT_COUNT.dateRange.dateDisplay('end'))
      .hasText('January 2024', 'end month is correctly parsed from STATIC_NOW');
    assert
      .dom(CLIENT_COUNT.card('Client usage trends'))
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} ${CHARTS.xAxisLabel}`)
      .hasText('7/23', 'x-axis labels start with billing start date');
    assert.dom(CHARTS.xAxisLabel).exists({ count: 7 }, 'chart months matches query');
  });

  module('community', function (hooks) {
    hooks.beforeEach(async function () {
      this.version.type = 'community';
      await login();
      return await visit('/vault/clients/counts/overview');
    });

    test('it should update charts when querying date ranges', async function (assert) {
      // Use parseAPITimestamp because we want a date string that is timezone agnostic (so it stays in UTC)
      const clientCountingStartDate = parseAPITimestamp(LICENSE_START.toISOString(), 'yyyy-MM');
      const upgradeMonth = parseAPITimestamp(UPGRADE_DATE.toISOString(), 'yyyy-MM');
      const endMonth = parseAPITimestamp(STATIC_PREVIOUS_MONTH.toISOString(), 'yyyy-MM');
      await click(CLIENT_COUNT.dateRange.edit);
      await fillIn(CLIENT_COUNT.dateRange.editDate('start'), clientCountingStartDate);
      await fillIn(CLIENT_COUNT.dateRange.editDate('end'), clientCountingStartDate);
      await click(GENERAL.submitButton);
      assert
        .dom(CLIENT_COUNT.usageStats('Client usage'))
        .exists('running total single month usage stats show');
      assert
        .dom(CLIENT_COUNT.card('Client usage trends'))
        .doesNotExist('running total month over month charts do not show');

      // change to start on month/year of upgrade to 1.10
      await click(CLIENT_COUNT.dateRange.edit);
      await fillIn(CLIENT_COUNT.dateRange.editDate('start'), upgradeMonth);
      await fillIn(CLIENT_COUNT.dateRange.editDate('end'), endMonth);
      await click(GENERAL.submitButton);

      assert
        .dom(CLIENT_COUNT.dateRange.dateDisplay('start'))
        .hasText('September 2023', 'client count start month is correctly parsed from start query');
      assert
        .dom(CLIENT_COUNT.card('Client usage trends'))
        .exists('Shows running totals with monthly breakdown charts');
      assert
        .dom(`${CLIENT_COUNT.card('Client usage trends')} ${CHARTS.xAxisLabel}`)
        .hasText('9/23', 'x-axis labels start with queried start month (upgrade date)');
      assert.dom(CHARTS.xAxisLabel).exists({ count: 4 }, 'chart months matches query');

      // query for single, historical month (upgrade month)
      await click(CLIENT_COUNT.dateRange.edit);
      await fillIn(CLIENT_COUNT.dateRange.editDate('start'), upgradeMonth);
      await fillIn(CLIENT_COUNT.dateRange.editDate('end'), upgradeMonth);
      await click(GENERAL.submitButton);
      assert
        .dom(CLIENT_COUNT.card('Client usage trends'))
        .exists('running total month over month charts show');

      // query historical date range (from September 2023 to December 2023)
      await click(CLIENT_COUNT.dateRange.edit);
      await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2023-09');
      await fillIn(CLIENT_COUNT.dateRange.editDate('end'), '2023-12');
      await click(GENERAL.submitButton);

      assert
        .dom(CLIENT_COUNT.dateRange.dateDisplay('start'))
        .hasText('September 2023', 'it displays correct start time');
      assert
        .dom(CLIENT_COUNT.dateRange.dateDisplay('end'))
        .hasText('December 2023', 'it displays correct end time');
      assert
        .dom(CLIENT_COUNT.card('Client usage trends'))
        .exists('Shows running totals with monthly breakdown charts');

      assert.dom(CHARTS.xAxisLabel).exists({ count: 4 }, 'chart months matches query');
      const xAxisLabels = findAll(CHARTS.xAxisLabel);
      assert
        .dom(xAxisLabels[xAxisLabels.length - 1])
        .hasText('12/23', 'x-axis labels end with queried end month');
    });

    test('it does not render client list links for community versions', async function (assert) {
      assert
        .dom(`${GENERAL.tableData(0, 'clients')} a`)
        .doesNotExist('client counts do not render as hyperlinks');
    });
  });

  test('it does not render client list links for HVD managed clusters', async function (assert) {
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    assert
      .dom(`${GENERAL.tableData(0, 'clients')} a`)
      .doesNotExist('client counts do not render as hyperlinks');
  });

  // * FILTERING ASSERTIONS
  // These tests use the static data from the ACTIVITY_RESPONSE_STUB to assert filtering
  // Filtering tests are split between integration and acceptance tests
  // because changing filters updates the URL query params.
  module('static data', function (hooks) {
    hooks.beforeEach(async function () {
      this.server.get('sys/internal/counters/activity', () => {
        return {
          request_id: 'some-activity-id',
          data: ACTIVITY_RESPONSE_STUB,
        };
      });
      const staticActivity = await this.store.findRecord('clients/activity', 'some-activity-id');
      this.staticMostRecentMonth = staticActivity.byMonth[staticActivity.byMonth.length - 1];
      await login();
      return visit('/vault/clients/counts/overview');
    });

    test('it filters attribution table when filters are applied', async function (assert) {
      const url = '/vault/clients/counts/overview';
      const topMount = flattenMounts(this.staticMostRecentMonth.new_clients.namespaces)[0];
      const timestamp = this.staticMostRecentMonth.timestamp;
      const { namespace_path, mount_type, mount_path } = topMount;
      assert.strictEqual(currentURL(), url, 'URL does not contain query params');
      await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
      await click(FILTERS.dropdownItem(timestamp));
      await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
      await click(FILTERS.dropdownItem(namespace_path));
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
      await click(FILTERS.dropdownItem(mount_path));
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
      await click(FILTERS.dropdownItem(mount_type));
      assert.strictEqual(
        currentURL(),
        `${url}?month=${encodeURIComponent(timestamp)}&mount_path=${encodeURIComponent(
          mount_path
        )}&mount_type=${mount_type}&namespace_path=${namespace_path}`,
        'url query params match filters'
      );
      assert.dom(FILTERS.tag()).exists({ count: 4 }, '4 filter tags render');
      assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'it only renders the filtered table row');
      assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(namespace_path);
      assert.dom(GENERAL.tableData(0, 'mount_type')).hasText(mount_type);
      assert.dom(GENERAL.tableData(0, 'mount_path')).hasText(mount_path);
    });

    test('it updates table when filters are cleared', async function (assert) {
      const url = '/vault/clients/counts/overview';
      const mounts = flattenMounts(this.staticMostRecentMonth.new_clients.namespaces);
      const timestamp = this.staticMostRecentMonth.timestamp;
      const { namespace_path, mount_type, mount_path } = mounts[0];
      await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
      await click(FILTERS.dropdownItem(timestamp));
      await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
      await click(FILTERS.dropdownItem(namespace_path));
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
      await click(FILTERS.dropdownItem(mount_path));
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
      await click(FILTERS.dropdownItem(mount_type));
      assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'it only renders the filtered table row');
      await click(FILTERS.clearTag(namespace_path));
      assert.strictEqual(
        currentURL(),
        `${url}?month=${encodeURIComponent(timestamp)}&mount_path=${encodeURIComponent(
          mount_path
        )}&mount_type=${mount_type}`,
        'url does not have namespace_path query param'
      );
      assert.dom(GENERAL.tableRow()).exists({ count: 2 }, 'it renders 2 data rows that match filters');
      assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText('root');
      assert.dom(GENERAL.tableData(0, 'mount_type')).hasText(mount_type);
      assert.dom(GENERAL.tableData(1, 'namespace_path')).hasText('ns1/');
      assert.dom(GENERAL.tableData(1, 'mount_type')).hasText(mount_type);
      assert.dom(GENERAL.tableData(1, 'mount_path')).hasText(mount_path);
      await click(GENERAL.button('Clear filters'));
      assert.strictEqual(currentURL(), url, 'url does not have any query params');
      assert
        .dom(GENERAL.tableRow())
        .exists({ count: mounts.length }, 'it renders all data when filters are cleared');
    });

    test('it renders client counts for full billing period when month is unselected', async function (assert) {
      const url = '/vault/clients/counts/overview';
      const mounts = flattenMounts(this.staticMostRecentMonth.new_clients.namespaces);
      const timestamp = this.staticMostRecentMonth.timestamp;
      const { namespace_path, mount_type, mount_path } = mounts[0];
      await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
      await click(FILTERS.dropdownItem(timestamp));
      await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
      await click(FILTERS.dropdownItem(namespace_path));
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
      await click(FILTERS.dropdownItem(mount_path));
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
      await click(FILTERS.dropdownItem(mount_type));
      assert.strictEqual(
        currentURL(),
        `${url}?month=${encodeURIComponent(timestamp)}&mount_path=${encodeURIComponent(
          mount_path
        )}&mount_type=${mount_type}&namespace_path=${namespace_path}`,
        'url query params match filters'
      );
      await click(FILTERS.clearTag('September 2023'));
      assert
        .dom(GENERAL.tableData(0, 'clients'))
        .hasText('4003', 'the table renders clients for the full billing period (not September)');
      assert.strictEqual(
        currentURL(),
        `${url}?mount_path=${encodeURIComponent(
          mount_path
        )}&mount_type=${mount_type}&namespace_path=${namespace_path}`,
        'url does not include month'
      );
    });

    test('enterprise: it navigates to the client list page when clicking the client count hyperlink', async function (assert) {
      const mockResponse = {
        status: 200,
        ok: true,
        text: () => Promise.resolve(ACTIVITY_EXPORT_STUB.trim()),
      };
      const adapter = this.store.adapterFor('clients/activity');
      const exportDataStub = sinon.stub(adapter, 'exportData');
      exportDataStub.resolves(mockResponse);
      const timestamp = this.staticMostRecentMonth.timestamp;
      await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
      await click(FILTERS.dropdownItem(timestamp));
      await click(`${GENERAL.tableData(0, 'clients')} a`);
      const url = '/vault/clients/counts/client-list';
      const monthQp = encodeURIComponent(timestamp);
      const ns = encodeURIComponent('ns1/');
      const mPath = encodeURIComponent('auth/userpass/0/');
      const mType = 'userpass';
      assert.strictEqual(
        currentURL(),
        `${url}?month=${monthQp}&mount_path=${mPath}&mount_type=${mType}&namespace_path=${ns}`,
        'url query params match filters'
      );
      exportDataStub.restore();
    });
  });

  module('license includes secrets sync feature', function (hooks) {
    hooks.beforeEach(async function () {
      syncHandler(this.server);
    });

    test('it should show secrets sync stats when the feature is activated', async function (assert) {
      await login();
      await visit('/vault/clients/counts/overview');
      assert
        .dom(CLIENT_COUNT.statLegendValue('Secret sync clients'))
        .exists('shows secret sync data on overview');
      await click(GENERAL.inputByAttr('toggle view'));
      assert
        .dom(CHARTS.legend)
        .hasText(
          'Entity clients Non-entity clients ACME clients Secret sync clients',
          'it renders legend in order that matches the stacked bar data'
        );
    });

    test('it should hide secrets sync stats when feature is NOT activated', async function (assert) {
      this.server.get('/sys/activation-flags', () => {
        return {
          data: { activated: [], unactivated: ['secrets-sync'] },
        };
      });

      await login();
      await visit('/vault/clients/counts/overview');
      assert
        .dom(CLIENT_COUNT.statLegendValue('Secret sync clients'))
        .doesNotExist('stat is hidden because feature is not activated');
      assert.dom(CLIENT_COUNT.statLegendValue('Entity clients')).exists('other stats are still visible');
      await click(GENERAL.inputByAttr('toggle view'));
      assert
        .dom(CHARTS.legend)
        .hasText(
          'Entity clients Non-entity clients ACME clients',
          'it renders legend in order that matches the stacked bar data and does not include secret sync'
        );
    });
  });
});
