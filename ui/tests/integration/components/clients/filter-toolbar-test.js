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

module('Integration | Component | clients/filter-toolbar', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.namespaces = ['root', 'admin/', 'ns1/'];
    this.mountPaths = ['auth/token/', 'auth/auto/eng/core/auth/core-gh-auth/', 'auth/userpass-root/'];
    this.mountTypes = ['token/', 'userpass/', 'ns_token/'];
    this.onFilter = sinon.spy();
    this.renderComponent = async () => {
      await render(hbs`
    <Clients::FilterToolbar
      @namespaces={{this.namespaces}}
      @mountPaths={{this.mountPaths}}
      @mountTypes={{this.mountTypes}}
      @onFilter={{this.onFilter}}
    />`);

      this.selectFilters = async () => {
        // select namespace_path
        await click(FILTERS.dropdownToggle('namespace_path'));
        await click(FILTERS.dropdownItem('admin/'));
        // select mount_path
        await click(FILTERS.dropdownToggle('mount_path'));
        await click(FILTERS.dropdownItem('auth/userpass-root/'));
        // select mount_type
        await click(FILTERS.dropdownToggle('mount_type'));
        await click(FILTERS.dropdownItem('token/'));
      };
    };
  });

  test('it renders dropdowns', async function (assert) {
    await this.renderComponent();

    assert.dom(FILTERS.dropdownToggle('namespace_path')).hasText('Namespace');
    assert.dom(FILTERS.dropdownToggle('mount_path')).hasText('Mount path');
    assert.dom(FILTERS.dropdownToggle('mount_type')).hasText('Mount type');
    assert.dom(GENERAL.button('Apply filters')).exists();
    assert
      .dom(GENERAL.button('Clear filters'))
      .doesNotExist('"Clear filters" button does not render by default');
  });

  test('it renders dropdown items', async function (assert) {
    await this.renderComponent();

    await click(FILTERS.dropdownToggle('namespace_path'));
    findAll('li button').forEach((item, idx) => {
      assert.dom(item).hasText(this.namespaces[idx]);
    });
    await click(FILTERS.dropdownToggle('mount_path'));
    findAll('li button').forEach((item, idx) => {
      assert.dom(item).hasText(this.mountPaths[idx]);
    });
    await click(FILTERS.dropdownToggle('mount_type'));
    findAll('li button').forEach((item, idx) => {
      assert.dom(item).hasText(this.mountTypes[idx]);
    });
  });

  test('it selects filters and renders a tag for each', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();

    // dropdown closes when an item is selected
    // reopen each one to assert the correct item is selected
    await click(FILTERS.dropdownToggle('namespace_path'));
    assert.dom(FILTERS.dropdownItem('admin/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('admin/')} ${GENERAL.icon('check')}`).exists();
    // select mount_path
    await click(FILTERS.dropdownToggle('mount_path'));
    assert.dom(FILTERS.dropdownItem('auth/userpass-root/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('auth/userpass-root/')} ${GENERAL.icon('check')}`).exists();
    // select mount_type
    await click(FILTERS.dropdownToggle('mount_type'));
    assert.dom(FILTERS.dropdownItem('token/')).hasAttribute('aria-selected', 'true');
    assert.dom(`${FILTERS.dropdownItem('token/')} ${GENERAL.icon('check')}`).exists();
    // Confirm filter tags render for each item
    assert.dom(FILTERS.tag()).exists({ count: 3 }, '3 filter tags render');
    assert.dom(FILTERS.tag('namespace_path', 'admin/')).exists();
    assert.dom(FILTERS.tag('mount_path', 'auth/userpass-root/')).exists();
    assert.dom(FILTERS.tag('mount_type', 'token/')).exists();
    assert
      .dom(GENERAL.button('Clear filters'))
      .exists('"Clear filters" button renders when filters are present');
  });

  test('it resets all filters', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();

    assert.dom(FILTERS.tag()).exists({ count: 3 }, '3 filter tags render');
    await click(GENERAL.button('Clear filters'));
    assert.dom(FILTERS.tag()).doesNotExist('tag filters disappear when "Clear filters" is clicked');
    assert
      .dom(GENERAL.button('Clear filters'))
      .doesNotExist('"Clear filters" button disappears when all filters are cleared');
    await click(GENERAL.button('Apply filters'));
    const [obj] = this.onFilter.lastCall.args;
    assert.propEqual(
      obj,
      { mount_path: '', mount_type: '', namespace_path: '' },
      'onFilter callback has empty values when filters are cleared'
    );
  });

  test('it clears individual filters', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();

    assert.dom(FILTERS.tag()).exists({ count: 3 }, '3 filter tags render');
    // Remove two of the filters
    await click(FILTERS.clearTag('admin/'));
    assert.dom(FILTERS.tag('namespace_path', 'admin/')).doesNotExist();
    await click(FILTERS.clearTag('auth/userpass-root/'));
    assert.dom(FILTERS.tag('mount_path', 'auth/userpass-root/')).doesNotExist();
    assert.dom(FILTERS.tag()).exists({ count: 1 }, '1 filter tags render');
    await click(GENERAL.button('Apply filters'));
    const [obj] = this.onFilter.lastCall.args;
    assert.propEqual(
      obj,
      { mount_path: '', mount_type: 'token/', namespace_path: '' },
      'onFilter callback has empty values for cleared filters'
    );
  });

  test('it applies filters', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();

    await click(GENERAL.button('Apply filters'));
    const [obj] = this.onFilter.lastCall.args;
    assert.strictEqual(obj.namespace_path, 'admin/', 'onFilter callback has expected "namespace_path"');
    assert.strictEqual(obj.mount_path, 'auth/userpass-root/', 'onFilter callback has expected "mount_path"');
    assert.strictEqual(obj.mount_type, 'token/', 'onFilter callback has expected "mount_type"');
  });

  test('it updates filters', async function (assert) {
    await this.renderComponent();
    await this.selectFilters();
    assert.dom(FILTERS.tag('mount_path', 'auth/userpass-root/')).exists();
    // selected a different mount path
    await click(FILTERS.dropdownToggle('mount_path'));
    await click(FILTERS.dropdownItem('auth/token/'));
    assert
      .dom(FILTERS.tag('mount_path', 'auth/userpass-root/'))
      .doesNotExist('"auth/userpass-root/" tag no longer exists');
    assert.dom(FILTERS.tag('mount_path', 'auth/token/')).exists('"auth/token/" tag renders');
    await click(GENERAL.button('Apply filters'));
    const [obj] = this.onFilter.lastCall.args;
    assert.strictEqual(obj.mount_path, 'auth/token/', 'onFilter callback has expected "mount_path"');
  });
});
