/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv-v2 | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.config = {
      cas_required: true,
      max_versions: 0,
      delete_version_after: '0s',
      type: 'kv',
      path: 'my-kv',
      accessor: 'kv_80616825',
      running_plugin_version: '2.7.0',
      local: false,
      seal_wrap: false,
      default_lease_ttl: '72h',
      max_lease_ttl: '123h',
      version: '2',
    };
    this.backend = 'my-kv';
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: 'my-kv', route: 'list' },
      { label: 'Configuration' },
    ];
  });

  test('it renders kv configuration details', async function (assert) {
    assert.expect(15);

    await render(
      hbs`
      <Page::Configuration
        @config={{this.config}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.title).includesText('my-kv', 'renders engine path as page title');
    assert.dom(PAGE.secretTab('Secrets')).exists('renders Secrets tab');
    assert.dom(PAGE.secretTab('Configuration')).exists('renders Configuration tab');

    assert.dom(PAGE.infoRowValue('Require check and set')).hasText('Yes');
    assert.dom(PAGE.infoRowValue('Automate secret deletion')).hasText('Never delete');
    assert.dom(PAGE.infoRowValue('Maximum number of versions')).hasText('0');
    assert.dom(PAGE.infoRowValue('Type')).hasText('kv');
    assert.dom(PAGE.infoRowValue('Path')).hasText('my-kv');
    assert.dom(PAGE.infoRowValue('Accessor')).hasText('kv_80616825');
    assert.dom(PAGE.infoRowValue('Running plugin version')).hasText('2.7.0');
    assert.dom(PAGE.infoRowValue('Local')).hasText('No');
    assert.dom(PAGE.infoRowValue('Seal wrap')).hasText('No');
    assert.dom(PAGE.infoRowValue('Default Lease TTL')).hasText('3 days');
    assert.dom(PAGE.infoRowValue('Max Lease TTL')).hasText('5 days 3 hours');
    assert.dom(PAGE.infoRowValue('Version')).hasText('2');
  });
});
