/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, find, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { ACTIVITY_EXPORT_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { CLIENT_COUNT, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { ClientFilters } from 'core/utils/client-counts/helpers';

const EXPORT_TAB_TO_TYPE = {
  Entity: 'entity',
  'Non-entity': 'non-entity-token',
  ACME: 'pki-acme',
  'Secret sync': 'secret-sync',
};

module('Integration | Component | clients/page/client-list', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.exportData = ACTIVITY_EXPORT_STUB.trim()
      .split('\n')
      .map((line) => JSON.parse(line));
    this.onFilterChange = sinon.spy();
    this.filterQueryParams = { namespace_path: '', mount_path: '', mount_type: '' };

    this.expectedData = (type, { key, value } = {}) =>
      this.exportData.filter((d) => {
        const isClientType = d.client_type === type;
        return key && value ? isClientType && d[key] === value : isClientType;
      });

    this.expectedOptions = (type) => [...new Set(this.exportData.map((m) => m[type]))];
    this.expectedNamespaces = this.expectedOptions('namespace_path');
    this.expectedMountPaths = this.expectedOptions('mount_path');
    this.expectedMountTypes = this.expectedOptions('mount_type');

    this.renderComponent = () =>
      render(hbs`
      <Clients::Page::ClientList
        @exportData={{this.exportData}}
        @onFilterChange={{this.onFilterChange}}
        @filterQueryParams={{this.filterQueryParams}}
      />`);

    // Filter key is one of ClientFilterTypes
    this.assertTabData = async (assert, filterKey, filterValue) => {
      // Iterate over each tab and assert rendered table data
      for (const [tabName, clientType] of Object.entries(EXPORT_TAB_TO_TYPE)) {
        const expectedData = this.expectedData(clientType, { key: filterKey, value: filterValue });
        const length = expectedData.length;
        await click(GENERAL.hdsTab(tabName));
        assert
          .dom(GENERAL.hdsTab(tabName))
          .hasText(`${tabName} ${length}`, `${tabName} tab counts match dataset length`);
        const noun = length === 1 ? 'client' : 'clients';
        const verb = length === 1 ? 'matches' : 'match';
        assert
          .dom(CLIENT_COUNT.tableSummary(tabName))
          .hasText(`Summary: ${length} ${noun} ${verb} the filter criteria.`);
        assert
          .dom(GENERAL.hdsTab(tabName))
          .hasAttribute('aria-selected', 'true', `it selects the tab: ${tabName}`);
        assert.dom(GENERAL.tableRow()).exists({ count: length });

        // Find all rendered rows and assert they satisfy the filter value and client IDs match
        const rows = findAll(GENERAL.tableRow());
        rows.forEach((_, idx) => {
          assert.dom(GENERAL.tableData(idx, filterKey)).hasText(filterValue);
          const clientId = find(GENERAL.tableData(idx, 'client_id')).innerText;
          // Make sure the rendered client id exists in the expected data
          const isValid = expectedData.find((d) => d.client_id === clientId);
          assert.true(!!isValid, `client_id: ${clientId} exists in expected dataset`);
        });
      }
    };
  });

  test('it renders export data by client type in tabs organized by client type', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.hdsTab('Entity')).hasAttribute('aria-selected', 'true', 'the first tab is selected');

    for (const [tabName, clientType] of Object.entries(EXPORT_TAB_TO_TYPE)) {
      const expectedData = this.expectedData(clientType);
      await click(GENERAL.hdsTab(tabName));
      assert
        .dom(GENERAL.hdsTab(tabName))
        .hasText(`${tabName} ${expectedData.length}`, `${tabName} tab counts match dataset length`);
      assert
        .dom(GENERAL.hdsTab(tabName))
        .hasAttribute('aria-selected', 'true', `it selects the tab: ${tabName}`);

      // Find all rendered rows and assert they match the client type tab
      const rows = findAll(GENERAL.tableRow());
      rows.forEach((_, idx) => {
        assert
          .dom(GENERAL.tableData(idx, 'client_type'))
          .hasText(clientType, `it renders ${clientType} data when ${tabName} is selected`);
      });
    }
  });

  test('it renders expected columns for each client type', async function (assert) {
    const expectedColumns = (isEntity = false) => {
      const base = [
        { label: 'Client ID' },
        { label: 'Client type' },
        { label: 'Namespace path' },
        { label: 'Namespace ID' },
        { label: 'Initial usage More information for' }, // renders a tooltip which is why "More information for" is included
        { label: 'Mount path' },
        { label: 'Mount type' },
        { label: 'Mount accessor' },
      ];
      const entityOnly = [
        { label: 'Entity name More information for' }, // renders a tooltip which is why "More information for" is included
        { label: 'Entity alias name' },
        { label: 'Local entity alias' },
        { label: 'Policies' },
        { label: 'Entity metadata' },
        { label: 'Entity alias metadata' },
        { label: 'Entity alias custom metadata' },
        { label: 'Entity group IDs' },
      ];
      return isEntity ? [...base, ...entityOnly] : base;
    };
    await this.renderComponent();

    for (const tabName of Object.keys(EXPORT_TAB_TO_TYPE)) {
      await click(GENERAL.hdsTab(tabName));
      expectedColumns(tabName === 'Entity').forEach((col, idx) => {
        assert
          .dom(GENERAL.tableColumnHeader(idx + 1, { isAdvanced: true }))
          .hasText(col.label, `${tabName} renders ${col.label} column`);
      });
    }
  });

  test('it renders dropdown lists from activity response to filter table data', async function (assert) {
    await this.renderComponent();
    // Select each filter
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = this.expectedNamespaces[idx] === '' ? 'root' : this.expectedNamespaces[idx];
      assert.dom(item).hasText(expected, `namespace dropdown renders: ${expected}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = this.expectedMountPaths[idx];
      assert.dom(item).hasText(expected, `mount_path dropdown renders: ${expected}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    findAll(FILTERS.dropdownItem()).forEach((item, idx) => {
      const expected = this.expectedMountTypes[idx];
      assert.dom(item).hasText(expected, `mount_type dropdown renders: ${expected}`);
    });
  });

  test('it fires @onFilterChange when filters are selected', async function (assert) {
    const ns = 'root';
    const { mount_path, mount_type } = this.exportData[0];
    await this.renderComponent();

    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    await click(FILTERS.dropdownItem(ns));
    // select mount path
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    await click(FILTERS.dropdownItem(mount_path));
    // select mount type
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await click(FILTERS.dropdownItem(mount_type));

    const [actual] = this.onFilterChange.lastCall.args;
    assert.strictEqual(actual.namespace_path, ns, `@onFilterChange called with: ${ns}`);
    assert.strictEqual(actual.mount_path, mount_path, `@onFilterChange called with: ${mount_path}`);
    assert.strictEqual(actual.mount_type, mount_type, `@onFilterChange called with: ${mount_type}`);
  });

  // *FILTERING TESTS
  test('it filters data if @filterQueryParams specify a namespace_path', async function (assert) {
    const filterKey = 'namespace_path';
    const filterValue = 'ns2/';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTabData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a mount_path', async function (assert) {
    const filterKey = 'mount_path';
    const filterValue = 'auth/token/';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTabData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a mount_type', async function (assert) {
    const filterKey = 'mount_type';
    const filterValue = 'auth/ns_token/';
    this.filterQueryParams[filterKey] = filterValue;
    await this.renderComponent();
    await this.assertTabData(assert, filterKey, filterValue);
  });

  test('it filters data if @filterQueryParams specify a multiple filters', async function (assert) {
    this.filterQueryParams = { namespace_path: 'ns5/', mount_path: 'auth/token/', mount_type: 'ns_token' };
    const { namespace_path, mount_path, mount_type } = this.filterQueryParams;
    await this.renderComponent();

    for (const [tabName, clientType] of Object.entries(EXPORT_TAB_TO_TYPE)) {
      const expectedData = this.expectedData(clientType).filter(
        (d) =>
          d.namespace_path == namespace_path && d.mount_path === mount_path && d.mount_type === mount_type
      );
      const length = expectedData.length;
      await click(GENERAL.hdsTab(tabName));
      assert
        .dom(GENERAL.hdsTab(tabName))
        .hasText(`${tabName} ${length}`, `${tabName} tab counts match dataset length`);
      const noun = length === 1 ? 'client' : 'clients';
      const verb = length === 1 ? 'matches' : 'match';
      assert
        .dom(CLIENT_COUNT.tableSummary(tabName))
        .hasText(`Summary: ${length} ${noun} ${verb} the filter criteria.`);
      assert
        .dom(GENERAL.hdsTab(tabName))
        .hasAttribute('aria-selected', 'true', `it selects the tab: ${tabName}`);
      assert.dom(GENERAL.tableRow()).exists({ count: length });

      // Find all rendered rows and assert they satisfy the filter value and client IDs match
      const rows = findAll(GENERAL.tableRow());
      rows.forEach((_, idx) => {
        assert.dom(GENERAL.tableData(idx, 'namespace_path')).hasText('ns5/');
        assert.dom(GENERAL.tableData(idx, 'mount_path')).hasText('auth/token/');
        assert.dom(GENERAL.tableData(idx, 'mount_type')).hasText('ns_token');
        // client_id is the unique identifier for each row
        const clientId = find(GENERAL.tableData(idx, 'client_id')).innerText;
        // Make sure the rendered client id exists in the expected data
        const isValid = expectedData.find((d) => d.client_id === clientId);
        assert.true(!!isValid, `client_id: ${clientId} exists in expected dataset`);
      });
    }
  });

  test('it renders empty state message when filter selections yield no results', async function (assert) {
    const flags = this.owner.lookup('service:flags');
    flags.activatedFlags = ['secrets-sync'];
    this.filterQueryParams = { namespace_path: 'dev/', mount_path: 'pluto/', mount_type: 'banana' };
    await this.renderComponent();

    for (const tabName of Object.keys(EXPORT_TAB_TO_TYPE)) {
      await click(GENERAL.hdsTab(tabName));
      assert
        .dom(CLIENT_COUNT.card('table empty state'))
        .hasText('No data found Select another client type or update filters to view client count data.');
    }
  });

  test('it renders empty state message when secret sync is not activated', async function (assert) {
    this.filterQueryParams = { namespace_path: 'dev/', mount_path: 'pluto/', mount_type: 'banana' };
    await this.renderComponent();
    await click(GENERAL.hdsTab('Secret sync'));
    assert
      .dom(CLIENT_COUNT.card('table empty state'))
      .hasText(
        'No secret sync clients No data is available because Secrets Sync has not been activated. Activate Secrets Sync'
      );
  });
});
