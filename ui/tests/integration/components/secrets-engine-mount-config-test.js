/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const selectors = {
  toggle: '[data-test-mount-config-toggle]',
  field: '[data-test-mount-config-field]',
};

module('Integration | Component | secrets-engine-mount-config', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const store = this.owner.lookup('service:store');
    store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        path: 'ldap-test/',
        type: 'ldap',
        accessor: 'ldap_7e838627',
        local: false,
        seal_wrap: true,
        config: {
          id: 'foo',
          default_lease_ttl: 0,
          max_lease_ttl: 10000,
        },
      },
    });
    this.model = store.peekRecord('secret-engine', 'ldap-test');
  });

  test('it should toggle config fields visibility', async function (assert) {
    await render(hbs`<SecretsEngineMountConfig @model={{this.model}} />`);

    assert
      .dom(selectors.toggle)
      .hasText('Show mount configuration', 'Correct toggle copy renders when closed');
    assert.dom(selectors.field).doesNotExist('Mount config fields are hidden');

    await click(selectors.toggle);

    assert.dom(selectors.toggle).hasText('Hide mount configuration', 'Correct toggle copy renders when open');
    assert.dom(selectors.field).exists('Mount config fields are visible');
  });

  test('it should render correct config fields', async function (assert) {
    await render(hbs`<SecretsEngineMountConfig @model={{this.model}} />`);
    await click(selectors.toggle);

    assert
      .dom(GENERAL.infoRowValue('Secret Engine Type'))
      .hasText(this.model.engineType, 'Secret engine type renders');
    assert.dom(GENERAL.infoRowValue('Path')).hasText(this.model.path, 'Path renders');
    assert.dom(GENERAL.infoRowValue('Accessor')).hasText(this.model.accessor, 'Accessor renders');
    assert.dom(GENERAL.infoRowValue('Local')).includesText('No', 'Local renders');
    assert.dom(GENERAL.infoRowValue('Seal Wrap')).includesText('Yes', 'Seal wrap renders');
    assert.dom(GENERAL.infoRowValue('Default Lease TTL')).includesText('0', 'Default Lease TTL renders');
    assert
      .dom(GENERAL.infoRowValue('Max Lease TTL'))
      .includesText('2 hours 46 minutes 40 seconds', 'Max Lease TTL renders');
  });

  test('it should yield block for additional fields', async function (assert) {
    await render(hbs`
      <SecretsEngineMountConfig @model={{this.model}}>
        <span data-test-yield>It Yields!</span>
      </SecretsEngineMountConfig>
    `);

    await click(selectors.toggle);
    assert.dom('[data-test-yield]').hasText('It Yields!', 'Component yields block for additional fields');
  });
});
