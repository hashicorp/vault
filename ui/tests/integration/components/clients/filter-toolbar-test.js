/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render, typeIn, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { ClientFilters } from 'core/utils/client-count-utils';

module('Integration | Component | clients/filter-toolbar', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.dataset = [
      { namespace_path: '', mount_type: 'userpass/', mount_path: 'auth/auto/eng/core/auth/core-gh-auth/' },
      { namespace_path: '', mount_type: 'userpass/', mount_path: 'auth/auto/eng/core/auth/core-gh-auth/' },
      { namespace_path: '', mount_type: 'userpass/', mount_path: 'auth/userpass-root/' },
      { namespace_path: 'admin/', mount_type: 'token/', mount_path: 'auth/token/' },
      { namespace_path: 'ns1/', mount_type: 'token/', mount_path: 'auth/token/' },
      { namespace_path: 'ns1/', mount_type: 'ns_token/', mount_path: 'auth/token/' },
    ];
    this.onFilter = sinon.spy();
    this.appliedFilters = { namespace_path: '', mount_path: '', mount_type: '' };

    this.renderComponent = async () => {
      await render(hbs`
    <Clients::FilterToolbar
      @dataset={{this.dataset}}
      @onFilter={{this.onFilter}}
      @appliedFilters={{this.appliedFilters}}
    />`);
    };

    this.presetFilters = () => {
      this.appliedFilters = {
        namespace_path: 'admin/',
        mount_path: 'auth/userpass-root/',
        mount_type: 'token/',
      };
    };

    this.selectFilters = async () => {
      // select namespace
      await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
      await click(FILTERS.dropdownItem('admin/'));
      // select mount path
      await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
      await click(FILTERS.dropdownItem('auth/userpass-root/'));
      // select mount type
      await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
      await click(FILTERS.dropdownItem('token/'));
    };
  });

  test('it renders dropdowns', async function (assert) {
    await this.renderComponent();

    assert.dom(FILTERS.dropdownToggle(ClientFilters.NAMESPACE)).hasText('Namespace');
    assert.dom(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH)).hasText('Mount path');
    assert.dom(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE)).hasText('Mount type');
    assert.dom(GENERAL.button('Apply filters')).exists();
    assert
      .dom(GENERAL.button('Clear filters'))
      .doesNotExist('"Clear filters" button does not render when filters are unset');
  });

  test('it renders dropdown items and does not include duplicates', async function (assert) {
    await this.renderComponent();
    const expectedNamespaces = ['root', 'admin/', 'ns1/'];
    const expectedMountPaths = [
      'auth/auto/eng/core/auth/core-gh-auth/',
      'auth/userpass-root/',
      'auth/token/',
    ];
    const expectedMountTypes = ['userpass/', 'token/', 'ns_token/'];

    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 namespaces');
    findAll('li button').forEach((item, idx) => {
      const ns = expectedNamespaces[idx];
      const msg =
        idx === 0 ? 'it renders empty string as "root" namespace_path' : `it renders namespace_path: ${ns}`;
      assert.dom(item).hasText(ns, msg);
    });

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 mount paths');
    findAll('li button').forEach((item, idx) => {
      const m = expectedMountPaths[idx];
      assert.dom(item).hasText(m, `it renders mount_path: ${m}`);
    });

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('li button').exists({ count: 3 }, 'list renders 3 mount types');
    findAll('li button').forEach((item, idx) => {
      const m = expectedMountTypes[idx];
      assert.dom(item).hasText(m, `it renders mount_type: ${m}`);
    });
  });

  test('it searches dropdown items', async function (assert) {
    await this.renderComponent();

    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.NAMESPACE), 'n');
    let dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 2);
    assert.dom('ul').hasText('admin/ ns1/', 'it renders matching namespaces');

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    await typeIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_PATH), 'eng');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 1);
    assert.dom('ul').hasText('auth/auto/eng/core/auth/core-gh-auth/', 'it renders matching mount paths');

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await typeIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_TYPE), 'token');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 2);
    assert.dom('ul').hasText('token/ ns_token/', 'it renders matching mount types');

    // confirm that search input is cleared and dropdown renders all items again when re-opened
    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('ul').hasText('root admin/ ns1/', 'it resets filter and renders all namespace path');
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert
      .dom('ul')
      .hasText(
        'auth/auto/eng/core/auth/core-gh-auth/ auth/userpass-root/ auth/token/',
        'it resets filter and renders all mount paths'
      );
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('ul').hasText('userpass/ token/ ns_token/', 'it resets filter and renders all mount types');
  });

  test('it searches renders no matches found message', async function (assert) {
    await this.renderComponent();

    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.NAMESPACE), 'no matches');
    let dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching namespaces');

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_PATH), 'no matches');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching mount paths');

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await fillIn(FILTERS.dropdownSearch(ClientFilters.MOUNT_TYPE), 'no matches');
    dropdownItems = findAll('li button');
    await waitUntil(() => dropdownItems.length === 0);
    assert.dom('ul').hasText('No matching mount types');
  });

  test('it renders no items to filter if dropdown is empty', async function (assert) {
    this.dataset = [{ namespace_path: null, mount_type: null, mount_path: null }];
    await this.renderComponent();
    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('ul').hasText('No namespaces to filter');
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom('ul').hasText('No mount paths to filter');
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('ul').hasText('No mount types to filter');
  });

  test('it renders no items to filter if dataset is missing expected keys', async function (assert) {
    this.dataset = [{ foo: null, bar: null, baz: null }];
    await this.renderComponent();
    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom('ul').hasText('No namespaces to filter');
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom('ul').hasText('No mount paths to filter');
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom('ul').hasText('No mount types to filter');
  });

  test('it selects dropdown items', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();

    // dropdown closes when an item is selected, reopen each one to assert the correct item is selected
    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    assert.dom(FILTERS.dropdownItem('admin/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('admin/')} ${GENERAL.icon('check')}`).exists();

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    assert.dom(FILTERS.dropdownItem('auth/userpass-root/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('auth/userpass-root/')} ${GENERAL.icon('check')}`).exists();

    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    assert.dom(FILTERS.dropdownItem('token/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('token/')} ${GENERAL.icon('check')}`).exists();
  });

  test('it applies filters when no filters are set', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();

    await click(GENERAL.button('Apply filters'));
    const [obj] = this.onFilter.lastCall.args;
    assert.strictEqual(
      obj[ClientFilters.NAMESPACE],
      'admin/',
      `onFilter callback has expected "${ClientFilters.NAMESPACE}"`
    );
    assert.strictEqual(
      obj[ClientFilters.MOUNT_PATH],
      'auth/userpass-root/',
      `onFilter callback has expected "${ClientFilters.MOUNT_PATH}"`
    );
    assert.strictEqual(
      obj[ClientFilters.MOUNT_TYPE],
      'token/',
      `onFilter callback has expected "${ClientFilters.MOUNT_TYPE}"`
    );
  });

  test('it applies updated filters when filters are preset', async function (assert) {
    this.appliedFilters = { namespace_path: 'ns1', mount_path: 'auth/token/', mount_type: 'ns_token/' };
    await this.renderComponent();
    // Check initial filters
    await click(GENERAL.button('Apply filters'));
    const [beforeUpdate] = this.onFilter.lastCall.args;
    assert.propEqual(beforeUpdate, this.appliedFilters, 'callback fires with preset filters');
    // Change filters and confirm callback has updated values
    await this.selectFilters();
    await click(GENERAL.button('Apply filters'));
    const [afterUpdate] = this.onFilter.lastCall.args;
    assert.propEqual(
      afterUpdate,
      { namespace_path: 'admin/', mount_path: 'auth/userpass-root/', mount_type: 'token/' },
      'callback fires with updated selection'
    );
  });

  test('it renders a tag for each filter', async function (assert) {
    this.presetFilters();
    await this.renderComponent();

    assert.dom(FILTERS.tag()).exists({ count: 3 }, '3 filter tags render');
    assert.dom(FILTERS.tag(ClientFilters.NAMESPACE, 'admin/')).exists();
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_PATH, 'auth/userpass-root/')).exists();
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_TYPE, 'token/')).exists();
    assert
      .dom(GENERAL.button('Clear filters'))
      .exists('"Clear filters" button renders when filters are present');
  });

  test('it resets all filters', async function (assert) {
    this.presetFilters();
    await this.renderComponent();
    // first check that filters have preset values
    await click(GENERAL.button('Apply filters'));
    const [beforeClear] = this.onFilter.lastCall.args;
    assert.propEqual(
      beforeClear,
      { namespace_path: 'admin/', mount_path: 'auth/userpass-root/', mount_type: 'token/' },
      'callback fires with preset filters'
    );
    // now clear filters and confirm values are cleared
    await click(GENERAL.button('Clear filters'));
    const [afterClear] = this.onFilter.lastCall.args;
    assert.propEqual(
      afterClear,
      { namespace_path: '', mount_path: '', mount_type: '' },
      'onFilter callback has empty values when "Clear filters" is clicked'
    );
  });

  test('it clears individual filters', async function (assert) {
    this.presetFilters();
    await this.renderComponent();
    // first check that filters have preset values
    await click(GENERAL.button('Apply filters'));
    const [beforeClear] = this.onFilter.lastCall.args;
    assert.propEqual(
      beforeClear,
      { namespace_path: 'admin/', mount_path: 'auth/userpass-root/', mount_type: 'token/' },
      'callback fires with preset filters'
    );
    await click(FILTERS.clearTag('admin/'));
    const afterClear = this.onFilter.lastCall.args[0];
    assert.propEqual(
      afterClear,
      { namespace_path: '', mount_path: 'auth/userpass-root/', mount_type: 'token/' },
      'onFilter callback fires with empty namespace_path'
    );
  });

  test('it only renders tags for supported filters', async function (assert) {
    this.appliedFilters = { start_time: '2025-08-31T23:59:59Z' };
    await this.renderComponent();
    assert
      .dom(GENERAL.button('Clear filters'))
      .doesNotExist('"Clear filters" button does not render when filters are unset');
    assert.dom(FILTERS.tag()).doesNotExist();
  });
});
