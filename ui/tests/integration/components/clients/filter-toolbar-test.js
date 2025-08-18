/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { ClientFilters } from 'core/utils/client-count-utils';

module('Integration | Component | clients/filter-toolbar', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.namespaces = ['root', 'admin/', 'ns1/'];
    this.mountPaths = ['auth/token/', 'auth/auto/eng/core/auth/core-gh-auth/', 'auth/userpass-root/'];
    this.mountTypes = ['token/', 'userpass/', 'ns_token/'];
    this.onFilter = sinon.spy();
    this.appliedFilters = { nsLabel: '', mountPath: '', mountType: '' };

    this.renderComponent = async () => {
      await render(hbs`
    <Clients::FilterToolbar
      @namespaces={{this.namespaces}}
      @mountPaths={{this.mountPaths}}
      @mountTypes={{this.mountTypes}}
      @onFilter={{this.onFilter}}
      @appliedFilters={{this.appliedFilters}}
    />`);
    };

    this.presetFilters = () => {
      this.appliedFilters = { nsLabel: 'admin/', mountPath: 'auth/userpass-root/', mountType: 'token/' };
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

  test('it renders dropdown items', async function (assert) {
    await this.renderComponent();

    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    findAll('li button').forEach((item, idx) => {
      assert.dom(item).hasText(this.namespaces[idx]);
    });
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    findAll('li button').forEach((item, idx) => {
      assert.dom(item).hasText(this.mountPaths[idx]);
    });
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    findAll('li button').forEach((item, idx) => {
      assert.dom(item).hasText(this.mountTypes[idx]);
    });
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
    this.appliedFilters = { mountPath: 'auth/token/', mountType: 'ns_token/', nsLabel: 'ns1' };
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
      { mountPath: 'auth/userpass-root/', mountType: 'token/', nsLabel: 'admin/' },
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
      { mountPath: 'auth/userpass-root/', mountType: 'token/', nsLabel: 'admin/' },
      'callback fires with preset filters'
    );
    // now clear filters and confirm values are cleared
    await click(GENERAL.button('Clear filters'));
    const [afterClear] = this.onFilter.lastCall.args;
    assert.propEqual(
      afterClear,
      { mountPath: '', mountType: '', nsLabel: '' },
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
      { mountPath: 'auth/userpass-root/', mountType: 'token/', nsLabel: 'admin/' },
      'callback fires with preset filters'
    );
    await click(FILTERS.clearTag('admin/'));
    const afterClear = this.onFilter.lastCall.args[0];
    assert.propEqual(
      afterClear,
      { mountPath: 'auth/userpass-root/', mountType: 'token/', nsLabel: '' },
      'onFilter callback fires with empty nsLabel'
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
