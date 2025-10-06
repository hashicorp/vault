/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const keyManagementMockModel = {
  secretsEngine: {
    accessor: 'keymgmt_accessor',
    config: {
      default_lease_ttl: 2073600,
      force_no_cache: false,
      listing_visibility: 'hidden',
      max_lease_ttl: 4320000,
    },
    description: 'hello',
    external_entropy_access: false,
    local: true,
    options: {},
    path: 'keymgmt/',
    plugin_version: '',
    running_plugin_version: 'v0.17.1+builtin',
    running_sha256: '',
    seal_wrap: false,
    type: 'keymgmt',
    uuid: '4ea92618-5b52-f89a-9cbe-b65dc7e65689',
    id: 'keymgmt',
    backendConfigurationLink: `vault.cluster.secrets.backend.configuration`,
  },
  versions: ['v0.17.1+builtin'],
};

module('Integration | Component | SecretEngine::Page::GeneralSettings', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyManagementMockModel;
  });

  test('it shows general settings form', async function (assert) {
    assert.expect(4);

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} />
    `);
    assert.dom(GENERAL.cardContainer('secrets duration')).exists(`Lease duration card exists`);
    assert.dom(GENERAL.cardContainer('security')).exists(`Security card exists`);
    assert.dom(GENERAL.cardContainer('version')).exists(`Version card exists`);
    assert.dom(GENERAL.cardContainer('metadata')).exists(`Metadata card exists`);
  });

  test('it shows unsaved changes modal', async function (assert) {
    assert.expect(3);

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} />
    `);
    await fillIn(GENERAL.textareaByAttr('description'), 'Some awesome description');
    await click(GENERAL.cancelButton);

    assert.dom(GENERAL.modal.container('unsaved-changes')).exists('Unsaved changes exists');
    assert.dom(GENERAL.modal.header('unsaved-changes')).hasText('Unsaved changes');
    assert
      .dom(GENERAL.modal.body('unsaved-changes'))
      .hasText(
        `You've made changes to the following ${this.model.secretsEngine.id} settings: Description Would you like to apply them?`
      );
  });

  test('it does not show unsaved changes modal when there are no unsaved changes', async function (assert) {
    assert.expect(1);

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} />
    `);
    await click(GENERAL.cancelButton);
    assert.dom(GENERAL.modal.container('unsaved-changes')).doesNotExist('Unsaved changes does not show');
  });
});
