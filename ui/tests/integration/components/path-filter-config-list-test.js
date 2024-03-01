/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll } from '@ember/test-helpers';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupEngine } from 'ember-engines/test-support';
import Service from '@ember/service';
import sinon from 'sinon';
import { Promise } from 'rsvp';
import { create } from 'ember-cli-page-object';
import ss from 'vault/tests/pages/components/search-select';
import { setRunOptions } from 'ember-a11y-testing/test-support';

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
    const ajaxStub = sinon.stub().usingPromise(Promise);
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
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
        // TODO: Fix groupname rendering on SearchSelect component
        'aria-required-parent': { enabled: false },
      },
    });
  });

  test('it renders', async function (assert) {
    this.set('config', { mode: null, paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} @paths={{this.paths}} />`, this.context);

    assert.dom('[data-test-component=path-filter-config]').exists();
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
    await click('#deny');
    assert.strictEqual(this.config.mode, 'deny');
    await click('#no-filtering');
    assert.strictEqual(this.config.mode, null);
  });

  test('it shows a warning when going from a mode to allow all', async function (assert) {
    this.set('config', { mode: 'allow', paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`, this.context);
    await click('#no-filtering');
    assert.dom('[data-test-remove-warning]').exists('shows removal warning');
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
