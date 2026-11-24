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
import { visit, click, fillIn, currentURL, findAll } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT, CHARTS, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import timestamp from 'core/utils/timestamp';
import {
  ACTIVITY_EXPORT_STUB,
  ACTIVITY_RESPONSE_STUB,
} from 'vault/tests/helpers/clients/client-count-helpers';
import { ClientFilters, flattenMounts } from 'core/utils/client-counts/helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { formatByMonths } from 'core/utils/client-counts/serializers';

module('Acceptance | clients | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    // These tests use the clientsHandler which dynamically generates activity data, used for asserting date querying, etc
    clientsHandler(this.server);
    this.version = this.owner.lookup('service:version');
  });

  test('it should hide secrets sync stats when feature is NOT on license', async function (assert) {
    // mocks endpoint for no additional license modules
    this.server.get('/sys/license/features', () => ({ features: [] }));

    await login();
    await visit('/vault/clients/counts/overview');
    const donutLegendLabels = findAll(CHARTS.carbonLegendLabel('Client count and type distribution')).map(
      (el) => el.textContent?.trim()
    );
    assert.notOk(
      donutLegendLabels.includes('Secret sync clients'),
      'donut legend does not include Secret sync clients'
    );
    assert.ok(donutLegendLabels.includes('Entity clients'), 'other stats are still visible');

    await click(GENERAL.inputByAttr('toggle view'));
    assert
      .dom('[data-test-chart="Client usage by month (stacked)"]')
      .exists('stacked chart container renders');
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
      .dom('[data-test-chart="Client usage by month (simple)"]')
      .exists('simple chart container renders for the default view');
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
        .dom('[data-test-chart="Client usage by month (simple)"]')
        .exists('simple chart container renders for the queried date range');

      // query for single, historical month (upgrade month)
      await click(CLIENT_COUNT.dateRange.edit);
      await fillIn(CLIENT_COUNT.dateRange.editDate('start'), upgradeMonth);
      await fillIn(CLIENT_COUNT.dateRange.editDate('end'), upgradeMonth);
      await click(GENERAL.submitButton);
      assert
        .dom(CLIENT_COUNT.card('Client usage trends'))
        .exists('running total month over month charts show');

      // query historical date range (from September 2023 to December 2023)
      const historicalStartMonth = '2023-09';
      const historicalEndMonth = '2023-12';
      await click(CLIENT_COUNT.dateRange.edit);
      await fillIn(CLIENT_COUNT.dateRange.editDate('start'), historicalStartMonth);
      await fillIn(CLIENT_COUNT.dateRange.editDate('end'), historicalEndMonth);
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

      assert
        .dom('[data-test-chart="Client usage by month (simple)"]')
        .exists('simple chart container remains rendered for the historical range');

      const xTickLabels = findAll(CHARTS.carbonXAxisTick('Client usage by month (simple)'))
        .map((el) => el.textContent?.trim())
        .filter(Boolean);
      const expectedStartTick = parseAPITimestamp(`${historicalStartMonth}-01T00:00:00Z`, 'M/yy');
      assert.true(xTickLabels.includes(expectedStartTick), 'x-axis includes queried start month');
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
      const byMonth = formatByMonths(ACTIVITY_RESPONSE_STUB.months);
      this.staticMostRecentMonth = byMonth[byMonth.length - 1];
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
      const api = this.owner.lookup('service:api');
      sinon.stub(api.sys, 'internalClientActivityExportRaw').resolves(mockResponse);
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
    });
  });

  module('license includes secrets sync feature', function (hooks) {
    hooks.beforeEach(async function () {
      syncHandler(this.server);
    });

    test('it should show secrets sync stats when the feature is activated', async function (assert) {
      await login();
      await visit('/vault/clients/counts/overview');
      const donutLegendLabels = findAll(CHARTS.carbonLegendLabel('Client count and type distribution')).map(
        (el) => el.textContent?.trim()
      );
      assert.ok(donutLegendLabels.includes('Secret sync clients'), 'shows secret sync data on overview');
      await click(GENERAL.inputByAttr('toggle view'));
      assert
        .dom('[data-test-chart="Client usage by month (stacked)"]')
        .exists('stacked chart container renders when the feature is activated');
    });

    test('it should hide secrets sync stats when feature is NOT activated', async function (assert) {
      this.server.get('/sys/activation-flags', () => {
        return {
          data: { activated: [], unactivated: ['secrets-sync'] },
        };
      });

      await login();
      await visit('/vault/clients/counts/overview');
      const donutLegendLabels = findAll(CHARTS.carbonLegendLabel('Client count and type distribution')).map(
        (el) => el.textContent?.trim()
      );
      assert.notOk(
        donutLegendLabels.includes('Secret sync clients'),
        'stat is hidden because feature is not activated'
      );
      assert.ok(donutLegendLabels.includes('Entity clients'), 'other stats are still visible');
      await click(GENERAL.inputByAttr('toggle view'));
      assert
        .dom('[data-test-chart="Client usage by month (stacked)"]')
        .exists('stacked chart container renders without secret sync in the stats view');
    });
  });
});
