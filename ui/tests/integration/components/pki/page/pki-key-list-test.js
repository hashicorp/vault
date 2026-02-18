/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { STANDARD_META } from 'vault/tests/helpers/pagination';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_KEYS } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | pki key list page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.keys = [
      {
        key_id: '724862ff-6438-bad0-b598-77a6c7f4e934',
        key_type: 'ec',
        key_name: 'test-key',
      },
      {
        key_id: '9fdddf12-9ce3-0268-6b34-dc1553b00175',
        key_type: 'rsa',
        key_name: 'another-key',
      },
    ];
    this.keys.meta = STANDARD_META;
    this.canImportKeys = true;
    this.canGenerateKeys = true;
    this.keyPermsById = {
      '724862ff-6438-bad0-b598-77a6c7f4e934': {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
      '9fdddf12-9ce3-0268-6b34-dc1553b00175': {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
    };

    this.renderComponent = () =>
      render(
        hbs`
        <Page::PkiKeyList
          @keys={{this.keys}}
          @mountPoint="vault.cluster.secrets.backend.pki"
          @canImportKeys={{this.canImportKeys}}
          @canGenerateKeys={{this.canGenerateKeys}}
          @keyPermsById={{this.keyPermsById}}
        />,
      `,
        { owner: this.engine }
      );
  });

  test('it renders empty state when no keys exist', async function (assert) {
    assert.expect(3);

    this.keys.meta = {
      currentPage: 1,
      total: 0,
      pageSize: 100,
    };

    await this.renderComponent();

    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No keys yet', 'renders empty state that no keys exist');
    assert.dom(PKI_KEYS.importKey).exists('renders toolbar with import action');
    assert.dom(PKI_KEYS.generateKey).exists('renders toolbar with generate action');
  });

  test('it renders list of keys and actions when permission allowed', async function (assert) {
    assert.expect(6);

    await this.renderComponent();

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

    this.canImportKeys = false;
    this.canGenerateKeys = false;
    this.keyPermsById = {
      '724862ff-6438-bad0-b598-77a6c7f4e934': { canRead: false, canUpdate: false },
      '9fdddf12-9ce3-0268-6b34-dc1553b00175': { canRead: false, canUpdate: false },
    };
    await this.renderComponent();

    assert.dom(PKI_KEYS.importKey).doesNotExist('renders import action');
    assert.dom(PKI_KEYS.generateKey).doesNotExist('renders generate action');
    assert.dom(GENERAL.menuTrigger).doesNotExist('does not render popup menu when no permission');
  });
});
