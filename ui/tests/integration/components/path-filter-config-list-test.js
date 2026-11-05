/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, findAll } from '@ember/test-helpers';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupEngine } from 'ember-engines/test-support';
import Service from '@ember/service';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);

const MOUNTS_RESPONSE = {
  data: {
    secret: {},
    auth: {
      'userpass/': { type: 'userpass', accessor: 'userpass' },
    },
  },
};
const NAMESPACE_MOUNTS_RESPONSE = {
  data: {
    secret: {
      'namespace-kv/': { type: 'kv', accessor: 'kv' },
    },
    auth: {},
  },
};

module('Integration | Component | path filter config list', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'replication');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    const ajaxStub = sinon.stub();
    ajaxStub.withArgs('/v1/sys/internal/ui/mounts', 'GET').resolves(MOUNTS_RESPONSE);
    ajaxStub
      .withArgs('/v1/sys/internal/ui/mounts', 'GET', { namespace: 'ns1' })
      .resolves(NAMESPACE_MOUNTS_RESPONSE);
    this.set('ajaxStub', ajaxStub);
    const namespaceServiceStub = Service.extend({
      init() {
        this._super(...arguments);
        this.set('accessibleNamespaces', ['ns1']);
      },
    });

    const storeServiceStub = Service.extend({
      adapterFor() {
        return {
          ajax: ajaxStub,
        };
      },
    });
    this.engine.register('service:namespace', namespaceServiceStub);
    this.engine.register('service:store', storeServiceStub);
  });

  test('it renders', async function (assert) {
    this.set('config', { mode: null, paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} @paths={{this.paths}} />`, this.context);

    assert.dom('[data-test-component=path-filter-config]').exists();
    assert.dom('[data-test-filter-mode="none"]').isChecked('Include everything is selected by default');
    assert.dom('[data-test-component="search-select"]').doesNotExist('No paths are needed without filtering');
    assert.dom('[data-test-remove-warning]').doesNotExist('A new configuration has no removal warning');
  });

  test.each(
    'it selects Include everything for an unset mode',
    [undefined, ''],
    async function (assert, mode) {
      this.set('config', { mode, paths: [] });
      await render(hbs`<PathFilterConfigList @config={{this.config}} />`, this.context);

      assert.dom('[data-test-filter-mode="none"]').isChecked('Falsy modes display Include everything');
      assert.dom('[data-test-remove-warning]').doesNotExist('An unset mode has no filter to remove');
    }
  );

  test.each('it displays an existing filter mode', ['allow', 'deny'], async function (assert, mode) {
    this.set('config', { mode, paths: ['auth/userpass/'] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`, this.context);

    assert.dom(`[data-test-filter-mode="${mode}"]`).isChecked('The saved filter mode is selected');
    assert.dom('[data-test-field-label]').hasText(`Paths in this ${mode}`, 'Paths label reflects the mode');
    assert.dom('[data-test-selected-option="0"]').hasText('auth/userpass/', 'Saved paths remain selected');
    assert.dom('[data-test-remove-warning]').doesNotExist('The existing filter is not being removed');
  });

  test('it sets config.paths', async function (assert) {
    this.set('config', { mode: 'allow', paths: [] });
    this.set('paths', []);
    await render(hbs`<PathFilterConfigList @config={{this.config}} @paths={{this.paths}} />`, this.context);

    await clickTrigger();
    await typeInSearch('auth');
    await searchSelect.options.objectAt(1).click();
    assert.ok(this.config.paths.includes('auth/userpass/'), 'adds to paths');

    await clickTrigger();
    await assert.strictEqual(searchSelect.options.length, 1, 'has one option left');

    await searchSelect.deleteButtons.objectAt(0).click();
    assert.strictEqual(this.config.paths.length, 0, 'removes from paths');
    await clickTrigger();
    await assert.strictEqual(searchSelect.options.length, 2, 'has both options');
  });

  test('it sets config.mode', async function (assert) {
    this.set('config', { mode: 'allow', paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`, this.context);
    await click('[data-test-filter-mode="deny"]');
    assert.strictEqual(this.config.mode, 'deny');
    assert.dom('[data-test-filter-mode="deny"]').isChecked('Deny is selected after changing mode');
    assert.dom('[data-test-filter-mode="allow"]').isNotChecked('Allow is deselected');
    assert.dom('[data-test-field-label]').hasText('Paths in this deny', 'Paths label updates to deny');
    await click('[data-test-filter-mode="none"]');
    assert.strictEqual(this.config.mode, null);
    assert.dom('[data-test-filter-mode="none"]').isChecked('Include everything is selected');
    assert.dom('[data-test-component="search-select"]').doesNotExist('Path selection is hidden');
    await click('[data-test-filter-mode="allow"]');
    assert.strictEqual(this.config.mode, 'allow', 'Allow can be reselected');
    assert.dom('[data-test-filter-mode="allow"]').isChecked('Allow is selected again');
    assert.dom('[data-test-field-label]').hasText('Paths in this allow', 'Paths label updates to allow');
  });

  test.each('it warns when removing an existing filter', ['allow', 'deny'], async function (assert, mode) {
    this.set('config', { mode, paths: ['auth/userpass/'] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`, this.context);
    await click('[data-test-filter-mode="none"]');
    assert.strictEqual(this.config.mode, null, 'Removing a filter sets mode to null, not a string');
    assert
      .dom('[data-test-remove-warning]')
      .hasTextContaining(
        "Saving with 'Include everything' will remove any existing filter configuration",
        'Warns that saving removes the existing filter'
      );
    assert.deepEqual(this.config.paths, ['auth/userpass/'], 'Changing mode does not discard selected paths');
    await click(`[data-test-filter-mode="${mode}"]`);
    assert.dom('[data-test-remove-warning]').doesNotExist('Reselecting a filter dismisses the warning');
    assert.dom('[data-test-selected-option="0"]').hasText('auth/userpass/', 'Selected paths are restored');
  });

  test('it does not warn when changing a new configuration back to Include everything', async function (assert) {
    this.set('config', { mode: null, paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`, this.context);

    await click('[data-test-filter-mode="allow"]');
    assert.strictEqual(this.config.mode, 'allow', 'A new configuration can select filtering');
    await click('[data-test-filter-mode="none"]');
    assert.strictEqual(this.config.mode, null, 'Include everything clears the mode');
    assert.dom('[data-test-remove-warning]').doesNotExist('There is no saved filter to remove');
  });

  test('it fetches mounts from a namespace when namespace name is entered', async function (assert) {
    this.set('config', { mode: 'allow', paths: [] });
    this.set('paths', []);
    await render(hbs`<PathFilterConfigList @config={{this.config}} @paths={{this.paths}} />`, this.context);

    await clickTrigger();
    assert.strictEqual(searchSelect.options.length, 2, 'shows userpass and namespace as an option');
    // type the namespace to trigger an ajax request
    await typeInSearch('ns1');
    assert.strictEqual(searchSelect.options.length, 2, 'has ns and ns mount in the list');
    await searchSelect.options.objectAt(1).click();
    assert.ok(this.config.paths.includes('ns1/namespace-kv/'), 'adds namespace mount to paths');
  });

  test('it selects mounts from different groups, and puts discarded option back within group', async function (assert) {
    this.set('config', { mode: 'allow', paths: [] });
    this.set('paths', []);
    await render(hbs`<PathFilterConfigList @config={{this.config}} @paths={{this.paths}} />`, this.context);
    await clickTrigger();
    await searchSelect.options.objectAt(1).click();
    await clickTrigger();
    await typeInSearch('ns1');
    await searchSelect.options.objectAt(1).click();
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    assert.dom('[data-test-selected-option="0"]').hasText('auth/userpass/', 'renders first mount selected');
    assert
      .dom('[data-test-selected-option="1"]')
      .hasText('ns1/namespace-kv/', 'renders second mount selected');
    assert.dom('[data-test-selected-option="2"]').hasText('ns1', 'renders third mount selected');
    assert.propEqual(
      this.config.paths,
      ['auth/userpass/', 'ns1/namespace-kv/', 'ns1'],
      'adds all selections to paths'
    );
    await searchSelect.deleteButtons.objectAt(0).click();
    await clickTrigger();
    assert
      .dom('.ember-power-select-group')
      .hasText('Auth Methods auth/userpass/', 'puts auth method back within group');
    await clickTrigger();
    await searchSelect.deleteButtons.objectAt(1).click();
    await clickTrigger();
    await typeInSearch('ns1');
    assert.dom('.ember-power-select-group').hasText('Namespaces ns1', 'puts ns back within group');
    await clickTrigger();
  });

  test('it renders previously set config.paths when editing the mount config', async function (assert) {
    this.set('config', { mode: 'allow', paths: ['auth/userpass/'] });
    this.set('paths', []);
    await render(hbs`<PathFilterConfigList @config={{this.config}} @paths={{this.paths}} />`, this.context);
    assert.strictEqual(
      searchSelect.selectedOptions.objectAt(0).text,
      'auth/userpass/',
      'renders config.path as selected on init'
    );
    await clickTrigger();
    assert.strictEqual(findAll('.ember-power-select-group').length, 1, 'renders only remaining group');
    await searchSelect.deleteButtons.objectAt(0).click();
    await clickTrigger();
    assert.strictEqual(findAll('.ember-power-select-group').length, 2, 'renders two groups');
  });
});
