/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, find, findAll, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { CLIENT_COUNT, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { ClientFilters, flattenMounts } from 'core/utils/client-count-utils';

module('Integration | Component | clients/page/overview', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });

    this.store = this.owner.lookup('service:store');
    this.activity = await this.store.queryRecord('clients/activity', {});
    this.mostRecentMonth = this.activity.byMonth[this.activity.byMonth.length - 1];
    this.onFilterChange = sinon.spy();
    this.filterQueryParams = { namespace_path: '', mount_path: '', mount_type: '' };
    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Overview 
        @activity={{this.activity}} 
        @onFilterChange={{this.onFilterChange}} 
        @filterQueryParams={{this.filterQueryParams}} 
      />`);

    this.assertTableData = async (assert, filterKey, filterValue) => {
      const expectedData = flattenMounts(this.mostRecentMonth.new_clients.namespaces).filter(
        (d) => d[filterKey] === filterValue
      );
      await fillIn(GENERAL.selectByAttr('attribution-month'), this.mostRecentMonth.timestamp);
      assert.dom(GENERAL.tableRow()).exists({ count: expectedData.length });
      // Find all rendered rows and assert they satisfy the filter value and table data matches expected values
      const rows = findAll(GENERAL.tableRow());
      rows.forEach((_, idx) => {
        assert.dom(GENERAL.tableData(idx, filterKey)).hasText(filterValue);
        // Get namespace and mount paths to find original data in expectedData
        const rowMountPath = find(GENERAL.tableData(idx, 'mount_path')).innerText;
        const rowNsPath = find(GENERAL.tableData(idx, 'namespace_path')).innerText;
        // find the expected clients from the response and assert the table matches
        const { clients: expectedClients } = expectedData.find(
          (d) => d.mount_path === rowMountPath && d.namespace_path === rowNsPath
        );
        assert.dom(GENERAL.tableData(idx, 'clients')).hasText(`${expectedClients}`);
      });
    };
  });

  test('it hides attribution when there is no data', async function (assert) {
    // Stub activity response when there's no activity data
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: {
          by_namespace: [],
          end_time: '2024-08-31T23:59:59Z',
          months: [],
          start_time: '2024-01-01T00:00:00Z',
          total: {
            distinct_entities: 0,
            entity_clients: 0,
            non_entity_tokens: 0,
            non_entity_clients: 0,
            clients: 0,
            secret_syncs: 0,
          },
        },
      };
    });
    this.activity = await this.store.queryRecord('clients/activity', {});
    await this.renderComponent();
    assert.dom(CLIENT_COUNT.card('Client attribution')).doesNotExist('it does not render attribution card');
    assert.dom(GENERAL.selectByAttr('attribution-month')).doesNotExist('it hides months dropdown');
  });

  test('it shows correct state message when selected month has no data', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '2023-06-01T00:00:00Z');

    assert
      .dom(CLIENT_COUNT.card('table empty state'))
      .hasText('No data found Clear or change filters to view client count data. Client count documentation');
  });

  test('it shows table when month selection has data', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '9/23');

    assert.dom(CLIENT_COUNT.card('table empty state')).doesNotExist('does not show card when table has data');
    assert.dom(GENERAL.table('attribution')).exists('shows table');
    assert.dom(GENERAL.paginationInfo).hasText('1–6 of 6', 'shows correct pagination info');
    assert.dom(GENERAL.paginationSizeSelector).hasValue('10', 'page size selector defaults to "10"');
  });

  test('it shows correct month options for billing period', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.selectByAttr('attribution-month')).exists('shows month selection dropdown');
    await fillIn(GENERAL.selectByAttr('attribution-month'), '');
    await triggerEvent(GENERAL.selectByAttr('attribution-month'), 'change');

    // assert that months options in select are those of selected billing period
    // '' represents default state of 'Select month'
    const expectedOptions = ['', ...this.activity.byMonth.reverse().map((m) => m.timestamp)];
    const actualOptions = findAll(`${GENERAL.selectByAttr('attribution-month')} option`).map(
      (option) => option.value
    );
    assert.deepEqual(actualOptions, expectedOptions, 'All <option> values match expected list');
  });

  test('it initially renders attribution with by_namespace data', async function (assert) {
    await this.renderComponent();
    const topNamespace = this.activity.byNamespace[0];
    const topMount = topNamespace.mounts[0];
    // Assert table renders namespace with the highest counts at the top
    assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(topNamespace.label);
    assert.dom(GENERAL.tableData(0, 'clients')).hasText(`${topMount.clients}`);
  });

  test('it renders dropdown lists from activity response to filter table data', async function (assert) {
    const mounts = flattenMounts(this.mostRecentMonth.new_clients.namespaces);
    const expectedNamespaces = [...new Set(mounts.map((m) => m.namespace_path))];
    const expectedMountPaths = [...new Set(mounts.map((m) => m.mount_path))];
    const expectedMountTypes = [...new Set(mounts.map((m) => m.mount_type))];
    await this.renderComponent();
    // Assert filters do not exist until month is selected
    assert.dom(FILTERS.dropdownToggle(ClientFilters.NAMESPACE)).doesNotExist();
    assert.dom(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH)).doesNotExist();
    assert.dom(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE)).doesNotExist();
    // Select month
    await fillIn(GENERAL.selectByAttr('attribution-month'), this.mostRecentMonth.timestamp);
    // Select each filter
    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    findAll(`${FILTERS.dropdown(ClientFilters.NAMESPACE)} li button`).forEach((item, idx) => {
      const expected = expectedNamespaces[idx];
      assert.dom(item).hasText(expected, `namespace dropdown renders: ${expected}`);
    });

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    findAll(`${FILTERS.dropdown(ClientFilters.MOUNT_PATH)} li button`).forEach((item, idx) => {
      const expected = expectedMountPaths[idx];
      assert.dom(item).hasText(expected, `mount_path dropdown renders: ${expected}`);
    });

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    findAll(`${FILTERS.dropdown(ClientFilters.MOUNT_TYPE)} li button`).forEach((item, idx) => {
      const expected = expectedMountTypes[idx];
      assert.dom(item).hasText(expected, `mount_type dropdown renders: ${expected}`);
    });
  });

  // * FILTERING ASSERTIONS
  // Filtering tests are split between integration and acceptance tests
  // because changing filters updates the URL query params/
  test('it filters attribution table by month', async function (assert) {
    await this.renderComponent();
    const mostRecentMonth = this.mostRecentMonth;
    await fillIn(GENERAL.selectByAttr('attribution-month'), mostRecentMonth.timestamp);
    // Drill down to new_clients then grab the first mount
    const sortedMounts = flattenMounts(mostRecentMonth.new_clients.namespaces).sort(
      (a, b) => b.clients - a.clients
    );
    const topMount = sortedMounts[0];
    assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(topMount.namespace_path);
    assert.dom(GENERAL.tableData(0, 'clients')).hasText(`${topMount.clients}`);
    assert.dom(GENERAL.tableData(0, 'mount_path')).hasText(topMount.mount_path);
  });

  test('it resets pagination when a month is selected change', async function (assert) {
    const attributionByMount = flattenMounts(this.activity.byNamespace);
    await this.renderComponent();
    // Decrease page size for test so we don't have to seed more data
    await fillIn(GENERAL.paginationSizeSelector, '5');
    assert.dom(GENERAL.paginationInfo).hasText(`1–5 of ${attributionByMount.length}`);
    // Change pages because we should go back to page 1 when a month is selected
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, '1 row render');
    assert.dom(GENERAL.paginationInfo).hasText(`6–6 of ${attributionByMount.length}`);
    // Select a month and assert table resets to page 1
    await fillIn(GENERAL.selectByAttr('attribution-month'), this.mostRecentMonth.timestamp);
    const monthMounts = flattenMounts(this.mostRecentMonth.new_clients.namespaces);
    assert
      .dom(GENERAL.paginationInfo)
      .hasText(`1–5 of ${monthMounts.length}`, 'pagination resets to page one');
    assert.dom(GENERAL.tableRow()).exists({ count: 5 }, '5 rows render');
    assert.dom(GENERAL.paginationSizeSelector).hasValue('5', 'size selector does not reset to 10');
  });

  test('it filters data if @filterQueryParams specify a namespace_path', async function (assert) {
    const filterKey = 'namespace_path';
    const filterValue = 'ns1';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTableData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a mount_path', async function (assert) {
    const filterKey = 'mount_path';
    const filterValue = 'acme/pki/0';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTableData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a mount_type', async function (assert) {
    const filterKey = 'mount_type';
    const filterValue = 'kv';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTableData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a multiple filters', async function (assert) {
    this.filterQueryParams = {
      namespace_path: 'ns1',
      mount_path: 'auth/userpass/0',
      mount_type: 'userpass',
    };

    const { namespace_path, mount_path, mount_type } = this.filterQueryParams;
    await this.renderComponent();
    await fillIn(GENERAL.selectByAttr('attribution-month'), this.mostRecentMonth.timestamp);
    const expectedData = flattenMounts(this.mostRecentMonth.new_clients.namespaces).find(
      (d) => d.namespace_path === namespace_path && d.mount_path === mount_path && d.mount_type === mount_type
    );
    assert.dom(GENERAL.tableRow()).exists({ count: 1 });
    assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(expectedData.namespace_path);
    assert.dom(GENERAL.tableData(0, 'mount_path')).hasText(expectedData.mount_path);
    assert.dom(GENERAL.tableData(0, 'mount_type')).hasText(expectedData.mount_type);
    assert.dom(GENERAL.tableData(0, 'clients')).hasText(`${expectedData.clients}`);
  });

  test('it renders empty state message when filter selections yield no results', async function (assert) {
    this.filterQueryParams = { namespace_path: 'dev/', mount_path: 'pluto/', mount_type: 'banana' };
    await this.renderComponent();
    assert
      .dom(CLIENT_COUNT.card('table empty state'))
      .hasText('No data found Clear or change filters to view client count data. Client count documentation');
  });
});
