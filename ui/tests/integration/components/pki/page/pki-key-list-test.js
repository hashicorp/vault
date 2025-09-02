/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { STANDARD_META } from 'vault/tests/helpers/pagination';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_KEYS } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | pki key list page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.store.pushPayload('pki/key', {
      modelName: 'pki/key',
      key_id: '724862ff-6438-bad0-b598-77a6c7f4e934',
      key_type: 'ec',
      key_name: 'test-key',
    });
    this.store.pushPayload('pki/key', {
      modelName: 'pki/key',
      key_id: '9fdddf12-9ce3-0268-6b34-dc1553b00175',
      key_type: 'rsa',
      key_name: 'another-key',
    });
    const keyModels = this.store.peekAll('pki/key');
    keyModels.meta = STANDARD_META;
    this.keyModels = keyModels;
  });

  test('it renders empty state when no keys exist', async function (assert) {
    assert.expect(3);
    this.keyModels = {
      meta: {
        total: 0,
        currentPage: 1,
        pageSize: 100,
      },
    };
    await render(
      hbs`
        <Page::PkiKeyList
          @keyModels={{this.keyModels}}
          @mountPoint="vault.cluster.secrets.backend.pki"
          @canImportKey={{true}}
          @canGenerateKey={{true}}
        />,
      `,
      { owner: this.engine }
    );
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No keys yet', 'renders empty state that no keys exist');
    assert.dom(PKI_KEYS.importKey).exists('renders toolbar with import action');
    assert.dom(PKI_KEYS.generateKey).exists('renders toolbar with generate action');
  });

  test('it renders list of keys and actions when permission allowed', async function (assert) {
    assert.expect(6);
    await render(
      hbs`
        <Page::PkiKeyList
          @keyModels={{this.keyModels}}
          @mountPoint="vault.cluster.secrets.backend.pki"
          @canImportKey={{true}}
          @canGenerateKey={{true}}
          @canRead={{true}}
          @canEdit={{true}}
        />,
      `,
      { owner: this.engine }
    );
    assert.dom(PKI_KEYS.keyName).hasText('test-key', 'linked block renders key id');
    assert
      .dom(PKI_KEYS.keyId)
      .hasText('724862ff-6438-bad0-b598-77a6c7f4e934', 'linked block renders key name');
    assert.dom(PKI_KEYS.importKey).exists('renders import action');
    assert.dom(PKI_KEYS.generateKey).exists('renders generate action');
    await click(GENERAL.menuTrigger);
    assert.dom(PKI_KEYS.popupMenuDetails).exists('details link exists');
    assert.dom(PKI_KEYS.popupMenuEdit).exists('edit link exists');
  });

  test('it hides actions when permission denied', async function (assert) {
    assert.expect(3);
    await render(
      hbs`
        <Page::PkiKeyList
          @keyModels={{this.keyModels}}
          @mountPoint="vault.cluster.secrets.backend.pki"
          @canImportKey={{false}}
          @canGenerateKey={{false}}
          @canRead={{false}}
          @canEdit={{false}}
        />,
      `,
      { owner: this.engine }
    );
    assert.dom(PKI_KEYS.importKey).doesNotExist('renders import action');
    assert.dom(PKI_KEYS.generateKey).doesNotExist('renders generate action');
    assert.dom(GENERAL.menuTrigger).doesNotExist('does not render popup menu when no permission');
  });
});
