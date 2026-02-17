/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render, typeIn, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { ClientFilters } from 'core/utils/client-counts/helpers';
import { parseAPITimestamp } from 'core/utils/date-formatters';

module('Integration | Component | clients/filter-toolbar', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.generateData = ({ withTimestamps = false }) => {
      const timestamps = [
        '2025-04-27T07:36:21Z',
        '2025-04-01T00:00:00Z',
        '2025-03-21T05:36:21Z',
        '2025-03-21T07:26:21Z',
        '2025-02-06T03:36:21Z',
        '2025-01-29T01:36:21Z',
      ];
      const data = [
        { namespace_path: '', mount_type: 'userpass/', mount_path: 'auth/auto/eng/core/auth/core-gh-auth/' },
        { namespace_path: '', mount_type: 'userpass/', mount_path: 'auth/auto/eng/core/auth/core-gh-auth/' },
        { namespace_path: '', mount_type: 'userpass/', mount_path: 'auth/userpass-root/' },
        { namespace_path: 'admin/', mount_type: 'token/', mount_path: 'auth/token/' },
        { namespace_path: 'ns1/', mount_type: 'token/', mount_path: 'auth/token/' },
        { namespace_path: 'ns1/', mount_type: 'ns_token/', mount_path: 'auth/token/' },
      ];
      // Only activity export data from Vault versions 1.21 or later will have a `client_first_used_time`
      return withTimestamps
        ? data.map((d, idx) => ({ ...d, client_first_used_time: timestamps[idx] }))
        : data;
    };

    this.onFilter = sinon.spy();
    this.filterQueryParams = { namespace_path: '', mount_path: '', mount_type: '', month: '' };
    this.dataset = undefined;
    this.dropdownMonths = undefined;
    this.isExportData = undefined;
    this.renderComponent = async () => {
      await render(hbs`
    <Clients::FilterToolbar
      @dataset={{this.dataset}}
      @onFilter={{this.onFilter}}
      @filterQueryParams={{this.filterQueryParams}}
      @isExportData={{this.isExportData}}
      @dropdownMonths={{this.dropdownMonths}}
    />`);
    };

    this.presetFilters = () => {
      this.filterQueryParams = {
        namespace_path: 'admin/',
        mount_path: 'auth/userpass-root/',
        mount_type: 'token/',
        month: '2025-04-01T00:00:00Z',
      };
    };

    this.selectFilters = async () => {
      // select namespace
      await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
      await click(FILTERS.dropdownItem('admin/'));
      // select mount path
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
      await click(FILTERS.dropdownItem('auth/userpass-root/'));
      // select mount type
      await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
      await click(FILTERS.dropdownItem('token/'));
      // select month
      await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
      await click(FILTERS.dropdownItem('2025-04-01T00:00:00Z'));
    };
  });

  test('it renders dropdowns when there is no data', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.dropdownToggle(ClientFilters.NAMESPACE)).hasText('Namespace');
    assert.dom(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH)).hasText('Mount path');
    assert.dom(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE)).hasText('Mount type');
    assert.dom(GENERAL.dropdownToggle(ClientFilters.MONTH)).hasText('Month');
    assert.dom(FILTERS.tagContainer).hasText('Filters applied: None');
  });

  test('it renders dropdown items and does not include duplicates', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    await this.renderComponent();
    const expectedNamespaces = ['root', 'admin/', 'ns1/'];
    const expectedMountPaths = [
      'auth/auto/eng/core/auth/core-gh-auth/',
      'auth/userpass-root/',
      'auth/token/',
    ];
    const expectedMountTypes = ['userpass/', 'token/', 'ns_token/'];
    // The component normalizes timestamps to the first of the month
    const expectedMonths = [
      '2025-04-01T00:00:00Z',
      '2025-03-01T00:00:00Z',
      '2025-02-01T00:00:00Z',
      '2025-01-01T00:00:00Z',
    ];

    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 namespaces');
    findAll('li button').forEach((item, idx) => {
      const ns = expectedNamespaces[idx];
      const msg =
        idx === 0 ? 'it renders empty string as "root" namespace_path' : `it renders namespace_path: ${ns}`;
      assert.dom(item).hasText(ns, msg);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 mount paths');
    findAll('li button').forEach((item, idx) => {
      const m = expectedMountPaths[idx];
      assert.dom(item).hasText(m, `it renders mount_path: ${m}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 mount types');
    findAll('li button').forEach((item, idx) => {
      const m = expectedMountTypes[idx];
      assert.dom(item).hasText(m, `it renders mount_type: ${m}`);
    });

    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom('li button').exists({ count: 4 }, 'list renders 4 months');
    findAll('li button').forEach((item, idx) => {
      const m = expectedMonths[idx];
      const display = parseAPITimestamp(m, 'MMMM yyyy');
      assert.dom(item).hasText(display, `it renders month: ${m}`);
    });
  });

  test('it renders months passed in as an arg instead of from dataset', async function (assert) {
    // Include timestamps in the dataset AND pass in months to ensure @dropdownMonths overrides the timestamps in dataset
    this.dataset = this.generateData({ withTimestamps: true });
    this.dropdownMonths = ['2025-10-01T07:36:21Z', '2025-09-01T02:38:21Z', '2025-08-01T03:56:21Z'];
    await this.renderComponent();

    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 months');
    findAll('li button').forEach((item, idx) => {
      const m = this.dropdownMonths[idx];
      const display = parseAPITimestamp(m, 'MMMM yyyy');
      assert.dom(item).hasText(display, `it renders month: ${m}`);
    });
  });

  test('it searches dropdown items', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    await this.renderComponent();
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.NAMESPACE), 'n');
    let dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 2);
    assert.dom('ul').hasText('admin/ ns1/', 'it renders matching namespaces');

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    await typeIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_PATH), 'eng');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 1);
    assert.dom('ul').hasText('auth/auto/eng/core/auth/core-gh-auth/', 'it renders matching mount paths');

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await typeIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_TYPE), 'token');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 2);
    assert.dom('ul').hasText('token/ ns_token/', 'it renders matching mount types');

    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    await typeIn(FILTERS.dropdownSearch(ClientFilters.MONTH), 'y');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 2);
    assert.dom('ul').hasText('February 2025 January 2025', 'it renders matching months');
    // Months can be searched by the ISO timestamp or the display value
    await fillIn(FILTERS.dropdownSearch(ClientFilters.MONTH), '4');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 1);
    assert.dom('ul').hasText('April 2025', 'it renders matching months');

    // Re-open each dropdown to confirm search input and dropdown reset after close
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('ul').hasText('root admin/ ns1/', 'namespace dropdown resets on close');
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert
      .dom('ul')
      .hasText(
        'auth/auto/eng/core/auth/core-gh-auth/ auth/userpass-root/ auth/token/',
        'mount path dropdown resets on close'
      );
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('ul').hasText('userpass/ token/ ns_token/', 'mount types dropdown resets on close');
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert
      .dom('ul')
      .hasText('April 2025 March 2025 February 2025 January 2025', 'months dropdown resets on close');
  });

  test('it searches and renders no matches found message', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    await this.renderComponent();

    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.NAMESPACE), 'no matches');
    let dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching namespaces');

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_PATH), 'no matches');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching mount paths');

    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_TYPE), 'no matches');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching mount types');

    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.MONTH), 'no matches');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching months');
  });

  test('it renders no items to filter if dropdown is empty', async function (assert) {
    this.dataset = [{ namespace_path: null, mount_type: null, mount_path: null, months: null }];
    await this.renderComponent();
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('ul').hasText('No namespaces to filter');
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom('ul').hasText('No mount paths to filter');
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('ul').hasText('No mount types to filter');
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom('ul').hasText('No months to filter');
  });

  test('it renders version message when no month data exists and @isExportData is true', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: false });
    this.isExportData = true;
    await this.renderComponent();
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert
      .dom('ul')
      .hasText(
        'Filtering by month is only available for clients initially used after upgrading to version 1.21.'
      );
  });

  test('it renders no months to filter message when data has no client_first_used_time', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: false });
    await this.renderComponent();
    assert.dom(GENERAL.dropdownToggle(ClientFilters.MONTH)).hasText('Month');
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom('ul').hasText('No months to filter');
  });

  test('it renders no months to filter message when @dropdownMonths is empty', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    this.dropdownMonths = [];
    await this.renderComponent();
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom('ul').hasText('No months to filter');
  });

  test('it renders no items to filter if dataset does not contain expected keys', async function (assert) {
    this.dataset = [{ foo: null, bar: null, baz: null }];
    await this.renderComponent();
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('ul').hasText('No namespaces to filter');
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom('ul').hasText('No mount paths to filter');
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('ul').hasText('No mount types to filter');
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom('ul').hasText('No months to filter');
  });

  test('it selects dropdown items and renders a filter tag', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    await this.renderComponent();

    // select namespace
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    await click(FILTERS.dropdownItem('admin/'));
    assert.dom(FILTERS.tag(ClientFilters.NAMESPACE, 'admin/')).exists();
    assert.dom(FILTERS.tag()).exists({ count: 1 }, '1 filter tag renders');
    // dropdown should close after an item is selected, reopen to assert the correct item is selected
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom(FILTERS.dropdownItem('admin/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('admin/')} ${GENERAL.icon('check')}`).exists();

    // select mount path
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    await click(FILTERS.dropdownItem('auth/userpass-root/'));
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_PATH, 'auth/userpass-root/')).exists();
    assert.dom(FILTERS.tag()).exists({ count: 2 }, '2 filter tags render');
    // dropdown should close after an item is selected, reopen to assert the correct item is selected
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom(FILTERS.dropdownItem('auth/userpass-root/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('auth/userpass-root/')} ${GENERAL.icon('check')}`).exists();

    // select mount type
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await click(FILTERS.dropdownItem('token/'));
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_TYPE, 'token/')).exists();
    assert.dom(FILTERS.tag()).exists({ count: 3 }, '3 filter tags render');
    // dropdown should close after an item is selected, reopen to assert the correct item is selected
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom(FILTERS.dropdownItem('token/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('token/')} ${GENERAL.icon('check')}`).exists();

    // select month
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    await click(FILTERS.dropdownItem('2025-02-01T00:00:00Z'));
    assert.dom(FILTERS.tag(ClientFilters.MONTH, '2025-02-01T00:00:00Z')).exists();
    assert.dom(FILTERS.tag()).exists({ count: 4 }, '4 filter tags render');
    // dropdown should close after an item is selected, reopen to assert the correct item is selected
    await click(GENERAL.dropdownToggle(ClientFilters.MONTH));
    assert.dom(FILTERS.dropdownItem('2025-02-01T00:00:00Z')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('2025-02-01T00:00:00Z')} ${GENERAL.icon('check')}`).exists();
  });

  test('it fires callback when a filter is selected', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    await this.renderComponent();

    // select namespace
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    await click(FILTERS.dropdownItem('admin/'));
    let lastCall = this.onFilter.lastCall.args[0];
    // this.filterQueryParams has empty values for each filter type
    let expectedObject = { ...this.filterQueryParams, [ClientFilters.NAMESPACE]: 'admin/' };
    assert.propEqual(lastCall, expectedObject, `callback includes value for ${ClientFilters.NAMESPACE}`);

    // select mount path
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    await click(FILTERS.dropdownItem('auth/userpass-root/'));
    lastCall = this.onFilter.lastCall.args[0];
    expectedObject = { ...expectedObject, [ClientFilters.MOUNT_PATH]: 'auth/userpass-root/' };
    assert.propEqual(lastCall, expectedObject, `callback includes value for ${ClientFilters.MOUNT_PATH}`);

    // select mount type
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await click(FILTERS.dropdownItem('token/'));
    lastCall = this.onFilter.lastCall.args[0];
    expectedObject = { ...expectedObject, [ClientFilters.MOUNT_TYPE]: 'token/' };
    assert.propEqual(lastCall, expectedObject, `callback includes value for ${ClientFilters.MOUNT_TYPE}`);
  });

  test('it renders filter tags when initialized with @filterQueryParams', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    this.presetFilters();
    await this.renderComponent();

    assert.dom(FILTERS.tag()).exists({ count: 4 }, '4 filter tags render');
    assert.dom(FILTERS.tag(ClientFilters.NAMESPACE, 'admin/')).exists();
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_PATH, 'auth/userpass-root/')).exists();
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_TYPE, 'token/')).exists();
    assert.dom(FILTERS.tag(ClientFilters.MONTH, '2025-04-01T00:00:00Z')).exists();
  });

  test('it updates filters tags when initialized with @filterQueryParams', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    this.filterQueryParams = {
      namespace_path: 'ns1/',
      mount_path: 'auth/token/',
      mount_type: 'ns_token/',
      month: '2025-03-01T00:00:00Z',
    };
    await this.renderComponent();
    // Check initial filters
    assert
      .dom(FILTERS.tagContainer)
      .hasText('Filters applied: ns1/ auth/token/ ns_token/ March 2025 Clear filters');
    // Change filters and confirm callback has updated values
    await this.selectFilters();
    const [afterUpdate] = this.onFilter.lastCall.args;
    assert.propEqual(
      afterUpdate,
      {
        namespace_path: 'admin/',
        mount_path: 'auth/userpass-root/',
        mount_type: 'token/',
        month: '2025-04-01T00:00:00Z',
      },
      'callback fires with updated selection'
    );
    assert
      .dom(FILTERS.tagContainer)
      .hasText('Filters applied: admin/ auth/userpass-root/ token/ April 2025 Clear filters');
  });

  test('it clears all filters', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    this.presetFilters();
    await this.renderComponent();
    await click(GENERAL.button('Clear filters'));
    const [afterClear] = this.onFilter.lastCall.args;
    assert.propEqual(
      afterClear,
      { namespace_path: '', mount_path: '', mount_type: '', month: '' },
      'onFilter callback has empty values when "Clear filters" is clicked'
    );
    assert.dom(FILTERS.tagContainer).hasText('Filters applied: None');
  });

  test('it clears individual filters', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    this.presetFilters();
    await this.renderComponent();
    await click(FILTERS.clearTag('admin/'));
    const afterClear = this.onFilter.lastCall.args[0];
    assert.propEqual(
      afterClear,
      {
        namespace_path: '',
        mount_path: 'auth/userpass-root/',
        mount_type: 'token/',
        month: '2025-04-01T00:00:00Z',
      },
      'onFilter callback fires with empty namespace_path'
    );
  });

  test('it renders an alert when initialized with @filterQueryParams that are not present in the dropdown', async function (assert) {
    this.dataset = this.generateData({ withTimestamps: true });
    this.filterQueryParams = { namespace_path: 'admin/', mount_path: '', mount_type: 'banana', month: '' };
    await this.renderComponent();
    assert.dom(FILTERS.tag()).exists({ count: 2 }, '2 filter tags render');
    assert.dom(FILTERS.tagContainer).hasText('Filters applied: admin/ banana Clear filters');
    assert.dom(GENERAL.inlineAlert).hasText(`Mount type "banana" not found in the current data.`);
  });
});
