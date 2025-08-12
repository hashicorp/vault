/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | plugin-status-indicator', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders version information', async function (assert) {
    this.set('version', 'v1.12.0');
    this.set('builtin', true);

    await render(hbs`
      <PluginStatusIndicator 
        @version={{this.version}}
        @builtin={{this.builtin}}
      />
    `);

    assert.dom('[data-test-plugin-status-indicator]').exists('Component renders');
    assert
      .dom('[data-test-plugin-version]')
      .containsText('v1.12.0', 'Version is displayed without builtin suffix');
    assert.dom('[data-test-plugin-type-badge]').containsText('Builtin', 'Builtin badge is shown');
  });

  test('it renders external plugin status', async function (assert) {
    this.set('version', 'v1.0.0');
    this.set('builtin', false);

    await render(hbs`
      <PluginStatusIndicator 
        @version={{this.version}}
        @builtin={{this.builtin}}
      />
    `);

    assert.dom('[data-test-plugin-version]').containsText('v1.0.0', 'External plugin version is displayed');
    assert.dom('[data-test-plugin-type-badge]').containsText('External', 'External badge is shown');
  });

  test('it renders deprecation status', async function (assert) {
    this.set('version', 'v1.0.0');
    this.set('builtin', true);
    this.set('deprecationStatus', 'pending-removal');

    await render(hbs`
      <PluginStatusIndicator 
        @version={{this.version}}
        @builtin={{this.builtin}}
        @deprecationStatus={{this.deprecationStatus}}
      />
    `);

    assert
      .dom('[data-test-plugin-deprecation-badge]')
      .containsText('Pending Removal', 'Deprecation badge is shown');
  });

  test('it handles missing plugin information gracefully', async function (assert) {
    await render(hbs`<PluginStatusIndicator />`);

    assert.dom('[data-test-plugin-status-indicator]').exists('Component renders even without data');
    assert.dom('[data-test-plugin-version]').doesNotExist('Version section not shown when no version');
    assert
      .dom('[data-test-plugin-type-badge]')
      .doesNotExist('Type badge not shown when builtin status unknown');
  });

  test('it only shows type badge when builtin status is defined', async function (assert) {
    this.set('version', 'v1.0.0');

    await render(hbs`
      <PluginStatusIndicator 
        @version={{this.version}}
      />
    `);

    assert.dom('[data-test-plugin-version]').containsText('v1.0.0', 'Version is shown');
    assert
      .dom('[data-test-plugin-type-badge]')
      .doesNotExist('Type badge not shown when builtin status is undefined');
  });
});
