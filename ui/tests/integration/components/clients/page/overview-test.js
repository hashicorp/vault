/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, find, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { CHARTS, CLIENT_COUNT, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { ClientFilters, flattenMounts } from 'core/utils/client-counts/helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';

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
    this.filterQueryParams = { namespace_path: '', mount_path: '', mount_type: '', month: '' };
    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::Overview 
        @activity={{this.activity}} 
        @onFilterChange={{this.onFilterChange}} 
        @filterQueryParams={{this.filterQueryParams}} 
      />`);

    this.assertTableData = async (assert, filterKey, filterValue) => {
      const expectedData = flattenMounts(this.activity.byNamespace).filter(
        (d) => d[filterKey] === filterValue
      );
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
    const expectedMonths = this.activity.byMonth
      .map((m) => parseAPITimestamp(m.timestamp, 'MMMM yyyy'))
      .reverse();
    const mounts = flattenMounts(this.activity.byNamespace);
    const expectedNamespaces = [...new Set(mounts.map((m) => m.namespace_path))];
    const expectedMountPaths = [...new Set(mounts.map((m) => m.mount_path))];
    const expectedMountTypes = [...new Set(mounts.map((m) => m.mount_type))];
    await this.renderComponent();

    // Select each filter
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = expectedMonths[idx];
      assert.dom(item).hasText(expected, `month dropdown renders: ${expected}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = expectedNamespaces[idx];
      assert.dom(item).hasText(expected, `namespace dropdown renders: ${expected}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = expectedMountPaths[idx];
      assert.dom(item).hasText(expected, `mount_path dropdown renders: ${expected}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = expectedMountTypes[idx];
      assert.dom(item).hasText(expected, `mount_type dropdown renders: ${expected}`);
    });
  });

  // * FILTERING ASSERTIONS
  // Filtering tests are split between integration and acceptance tests
  // because changing filters updates the URL query params

  test('it shows correct empty state message when selected month has no data', async function (assert) {
    this.filterQueryParams[ClientFilters.MONTH] = '2023-06-01T00:00:00Z';
    await this.renderComponent();
    assert
      .dom(CLIENT_COUNT.card('table empty state'))
      .hasText('No data found Clear or change filters to view client count data. Client count documentation');
  });

  test('it renders NEW monthly clients for self-managed clusters instead of total clients', async function (assert) {
    this.filterQueryParams = {
      month: this.mostRecentMonth.timestamp,
    };
    const topMount = this.mostRecentMonth.new_clients.namespaces
      .find((ns) => ns.label === 'ns1/')
      .mounts.find((m) => m.label === 'auth/userpass/0/');

    await this.renderComponent();
    assert.dom(CHARTS.legend).hasText('New clients');
    assert
      .dom(GENERAL.tableData(0, 'clients'))
      .hasText(`${topMount.clients}`, 'table renders total monthly clients');
  });

  test('it filters data if @filterQueryParams specify a month', async function (assert) {
    const filterKey = 'month';
    const filterValue = this.mostRecentMonth.timestamp;
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    // Drill down to new_clients then grab the first mount
    const sortedMounts = flattenMounts(this.mostRecentMonth.new_clients.namespaces).sort(
      (a, b) => b.clients - a.clients
    );
    const topMount = sortedMounts[0];
    assert.dom(GENERAL.tableData(0, 'namespace_path')).hasText(topMount.namespace_path);
    assert.dom(GENERAL.tableData(0, 'clients')).hasText(`${topMount.clients}`);
    assert.dom(GENERAL.tableData(0, 'mount_path')).hasText(topMount.mount_path);
  });

  test('it filters data if @filterQueryParams specify a namespace_path', async function (assert) {
    const filterKey = 'namespace_path';
    const filterValue = 'ns1/';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTableData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a mount_path', async function (assert) {
    const filterKey = 'mount_path';
    const filterValue = 'acme/pki/0/';
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
      month: this.mostRecentMonth.timestamp,
      namespace_path: 'ns1/',
      mount_path: 'auth/userpass/0/',
      mount_type: 'userpass',
    };

    const { namespace_path, mount_path, mount_type } = this.filterQueryParams;
    await this.renderComponent();
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

  test('it renders TOTAL monthly clients for HVD instead of new clients', async function (assert) {
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    this.filterQueryParams = {
      month: this.mostRecentMonth.timestamp,
    };
    const topMount = this.mostRecentMonth.namespaces
      .find((ns) => ns.label === 'root')
      .mounts.find((m) => m.label === 'acme/pki/0/');

    await this.renderComponent();
    assert.dom(CHARTS.legend).hasText('Clients');
    assert
      .dom(GENERAL.tableData(0, 'clients'))
      .hasText(`${topMount.clients}`, 'table renders total monthly clients');
  });
});
