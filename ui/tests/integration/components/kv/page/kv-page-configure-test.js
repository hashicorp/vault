/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import KvForm from 'vault/forms/secrets/kv';

module('Integration | Component | kv-v2 | Page::Configure', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.config = {
      cas_required: true,
      max_versions: 5,
      delete_version_after: '10000s',
    };
    this.form = new KvForm({});
    this.editForm = new KvForm(this.config);
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
      { label: 'Configuration', route: 'configuration', model: this.backend },
      { label: 'Edit' },
    ];
  });

  test('it renders kv configure form', async function (assert) {
    assert.expect(3);

    await render(
      hbs`
      <Page::Configure
        @form={{this.form}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(GENERAL.inputByAttr('max_versions')).exists();
    assert.dom(GENERAL.inputByAttr('cas_required')).exists();
    assert.dom(GENERAL.toggleInput('Automate secret deletion')).exists();
  });

  test('it renders kv configure form with existing config', async function (assert) {
    assert.expect(5);
    this.form = this.editForm;

    await render(
      hbs`
      <Page::Configure
        @form={{this.form}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(GENERAL.inputByAttr('max_versions')).hasValue('5');
    assert.dom(GENERAL.inputByAttr('cas_required')).hasValue('on');
    assert.dom(GENERAL.toggleInput('Automate secret deletion')).isChecked();
    assert.dom(GENERAL.ttl.input('Automate secret deletion')).hasValue('10000');
    assert.dom(GENERAL.selectByAttr('ttl-unit')).hasValue('s');
  });
});
