/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';

const selectors = {
  toggle: '[data-test-mount-config-toggle]',
  field: '[data-test-mount-config-field]',
};

module('Integration | Component | secrets-engine-mount-config', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.secretsEngine = new SecretsEngineResource({
      path: 'ldap-test/',
      type: 'ldap',
      accessor: 'ldap_7e838627',
      local: false,
      sealWrap: true,
      config: {
        id: 'foo',
        defaultLeaseTtl: 0,
        maxLeaseTtl: 10000,
      },
    });
  });

  test('it should toggle config fields visibility', async function (assert) {
    await render(hbs`<SecretsEngineMountConfig @secretsEngine={{this.secretsEngine}} />`);

    assert
      .dom(selectors.toggle)
      .hasText('Show mount configuration', 'Correct toggle copy renders when closed');
    assert.dom(selectors.field).doesNotExist('Mount config fields are hidden');

    await click(selectors.toggle);

    assert.dom(selectors.toggle).hasText('Hide mount configuration', 'Correct toggle copy renders when open');
    assert.dom(selectors.field).exists('Mount config fields are visible');
  });

  test('it should render correct config fields', async function (assert) {
    await render(hbs`<SecretsEngineMountConfig @secretsEngine={{this.secretsEngine}} />`);
    await click(selectors.toggle);

    assert
      .dom(GENERAL.infoRowValue('Secret engine type'))
      .hasText(this.secretsEngine.engineType, 'Secret engine type renders');
    assert.dom(GENERAL.infoRowValue('Path')).hasText(this.secretsEngine.path, 'Path renders');
    assert.dom(GENERAL.infoRowValue('Accessor')).hasText(this.secretsEngine.accessor, 'Accessor renders');
    assert.dom(GENERAL.infoRowValue('Local')).includesText('No', 'Local renders');
    assert.dom(GENERAL.infoRowValue('Seal wrap')).includesText('Yes', 'Seal wrap renders');
    assert.dom(GENERAL.infoRowValue('Default Lease TTL')).includesText('0', 'Default Lease TTL renders');
    assert
      .dom(GENERAL.infoRowValue('Max Lease TTL'))
      .includesText('2 hours 46 minutes 40 seconds', 'Max Lease TTL renders');
  });

  test('it should yield block for additional fields', async function (assert) {
    await render(hbs`
      <SecretsEngineMountConfig @secretsEngine={{this.secretsEngine}}>
        <span data-test-yield>It Yields!</span>
      </SecretsEngineMountConfig>
    `);

    await click(selectors.toggle);
    assert.dom('[data-test-yield]').hasText('It Yields!', 'Component yields block for additional fields');
  });
});
