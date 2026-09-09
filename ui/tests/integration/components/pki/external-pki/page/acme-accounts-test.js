/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { CreationMethod } from 'vault/utils/constants/snippet';

module('Integration | Component | pki | external-pki | ExternalPki::Page::AcmeAccounts', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.setupModel = (acmeAccounts = []) => {
      return {
        engine: new SecretsEngineResource({
          accessor: 'pki-external-ca_e158c567',
          type: 'pki-external-ca',
          path: 'my-pki-external-ca/',
        }),
        acmeAccounts,
      };
    };
    this.renderComponent = () =>
      render(hbs`<ExternalPki::Page::AcmeAccounts @model={{this.model}} />`, { owner: this.engine });
  });

  test('it renders empty state when no ACME accounts exist', async function (assert) {
    this.model = this.setupModel([]);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No ACME accounts exist yet');
    // Implementation select should be visible
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)).exists();
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.APICLI)).exists();
    assert.dom(GENERAL.textDisplay('0')).hasText('Configure an ACME account');
  });

  module('with ACME accounts', function (hooks) {
    hooks.beforeEach(function () {
      this.acmeAccounts = [
        {
          name: 'letsencrypt-prod',
          directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
          email_contacts: ['admin@example.com'],
          active_key_version: 2, // Last key is always active key
          account_keys: {
            0: {
              key_version: 0,
              key_type: 'rsa-2048',
              key_creation_date: '2024-01-15T10:30:00Z',
            },
            1: {
              key_version: 1,
              key_type: 'ec-256',
              key_creation_date: '2025-02-25T11:25:00Z',
            },
            2: {
              key_version: 2,
              key_type: 'ec-521',
              key_creation_date: '2026-07-07T16:35:12Z',
            },
          },
        },
        {
          name: 'dev-account',
          directory_url: 'https://acme-dev-v02.api.letsencrypt.org/directory',
          email_contacts: ['admin@example.com'],
          active_key_version: 1, // Last key is always active key
          account_keys: {
            0: {
              key_version: 0,
              key_type: 'rsa-2048',
              key_creation_date: '2024-01-15T10:30:00Z',
            },
            1: {
              key_version: 1,
              key_type: 'ec-256',
              key_creation_date: '2025-02-25T11:25:00Z',
            },
          },
        },
        {
          name: 'letsencrypt-staging',
          directory_url: 'https://acme-staging-v02.api.letsencrypt.org/directory',
          email_contacts: ['staging@example.com'],
          active_key_version: 0,
          account_keys: {
            0: {
              key_version: 0,
              key_type: 'rsa-2048',
              key_creation_date: '2026-07-07T16:35:12Z',
            },
          },
        },
      ];
      this.model = this.setupModel(this.acmeAccounts);
    });

    test('it renders list of ACME accounts', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.emptyStateTitle).doesNotExist();
      assert.dom(GENERAL.radioCardByAttr()).doesNotExist();

      assert.dom(GENERAL.cardContainer()).exists({ count: 3 });
      const accounts = ['letsencrypt-prod', 'dev-account', 'letsencrypt-staging'];
      accounts.forEach((a) => {
        assert.dom(`${GENERAL.cardContainer(a)} ${GENERAL.textDisplay()}`).hasText(a);
      });
      assert.dom(GENERAL.button('View key history 0')).exists();
      assert.dom(GENERAL.button('View key history 1')).exists();
      assert
        .dom(GENERAL.button('View key history 2'))
        .doesNotExist('it does not render button when no key history');
    });

    test('it opens flyout when "View key history" is clicked', async function (assert) {
      await this.renderComponent();
      // Flyout should not be visible initially
      assert.dom('#key-history-flyout').doesNotExist();
      // Click first account's "View key history" button
      await click(GENERAL.button('View key history 0'));
      // Flyout should be visible
      assert.dom('#key-history-flyout').exists();
      assert.dom('#key-history-flyout h1').hasText(`letsencrypt-prod key history`);
      // Active key details should display
      assert.dom(GENERAL.infoRowValue('Key type')).hasText('ec-521');
      assert.dom(GENERAL.infoRowValue('Key version')).hasText('2');
      assert.dom(GENERAL.infoRowValue('Creation date')).hasTextContaining('Jul 7, 2026');
      // Table should render inactive keys
      assert.dom(GENERAL.tableRow()).exists({ count: 2 });
      const flyoutAccount = this.acmeAccounts['0'];
      const expectedDates = ['Jan 15, 2024', 'Feb 25, 2025'];
      for (const keyIdx in flyoutAccount.account_keys) {
        const value = flyoutAccount.account_keys[keyIdx];
        // Active key should not appear in table
        if (keyIdx === '2') continue;
        assert.dom(GENERAL.tableData(keyIdx, 'key_creation_date')).hasTextContaining(expectedDates[keyIdx]);
        assert.dom(GENERAL.tableData(keyIdx, 'key_type')).hasText(value.key_type);
        assert.dom(GENERAL.tableData(keyIdx, 'key_version')).hasText(String(value.key_version));
      }
      assert.dom(GENERAL.tableRow('2')).doesNotExist('active key is NOT in table');
      await click('button[aria-label="Dismiss"]');
      assert.dom('#key-history-flyout').doesNotExist('flyout disappears on dismiss');
    });

    test('it switches between different account flyouts', async function (assert) {
      await this.renderComponent();
      // Click first account's "View key history" button
      await click(GENERAL.button('View key history 0'));
      // Flyout renders first account
      assert.dom('#key-history-flyout h1').hasText(`letsencrypt-prod key history`);
      assert.dom(GENERAL.infoRowValue('Key version')).hasText('2');
      assert.dom(GENERAL.infoRowValue('Key type')).hasText('ec-521');
      assert.dom(GENERAL.tableRow()).exists({ count: 2 });
      await click('button[aria-label="Dismiss"]');
      // Flyout renders updated account
      await click(GENERAL.button('View key history 1'));
      assert.dom('#key-history-flyout h1').hasText(`dev-account key history`);
      assert.dom(GENERAL.infoRowValue('Key version')).hasText('1');
      assert.dom(GENERAL.infoRowValue('Key type')).hasText('ec-256');
      assert.dom(GENERAL.tableRow()).exists({ count: 1 });
      await click('button[aria-label="Dismiss"]');
    });
  });
});
