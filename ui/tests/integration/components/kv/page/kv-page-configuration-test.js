/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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
    this.backend = {
      accessor: 'kv_05319fa9',
      config: {
        default_lease_ttl: 2764800,
        force_no_cache: false,
        listing_visibility: 'hidden',
        max_lease_ttl: 2764800,
      },
      description: '',
      external_entropy_access: false,
      local: false,
      options: {
        version: '2',
      },
      path: 'my-kv/',
      plugin_version: '',
      running_plugin_version: 'v0.25.0+builtin',
      running_sha256: '',
      seal_wrap: false,
      type: 'kv',
      uuid: '0cd6346f-c93a-ecfa-b01d-6b690a745c8e',
      id: 'my-kv',
    };
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: 'my-kv', route: 'list' },
      { label: 'Configuration' },
    ];
  });

  test('it renders kv configuration details', async function (assert) {
    assert.expect(6);

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

    assert.dom(PAGE.title).includesText('my-kv configuration', 'renders engine path as page title');
    assert.dom(GENERAL.tab('general-settings')).exists('renders general settings tab');
    assert.dom(GENERAL.tab('plugin-settings')).exists('renders kv settings tab');

    await click(GENERAL.tab('plugin-settings'));
    assert.dom(PAGE.infoRowValue('Require check and set')).hasText('Yes');
    assert.dom(PAGE.infoRowValue('Automate secret deletion')).hasText('Never delete');
    assert.dom(PAGE.infoRowValue('Maximum number of versions')).hasText('0');
  });
});
