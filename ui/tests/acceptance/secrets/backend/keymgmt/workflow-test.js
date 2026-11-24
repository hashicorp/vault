/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, settled, visit, waitUntil } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupApplicationTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';

import { create } from 'ember-cli-page-object';
import secretsNavTestHelper from 'vault/tests/acceptance/secrets/secrets-nav-test-helper';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

const consoleComponent = create(consoleClass);

module('Acceptance | Enterprise | keymgmt-configuration-workflow', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    // Mock the plugins pins request for keymgmt
    this.server.get('/sys/plugins/pins/secret/keymgmt', () => ({
      data: null,
    }));

    return login();
  });

  secretsNavTestHelper(test, 'keymgmt');

  for (const tab of ['Keys', 'Providers']) {
    test(`it filters Key Management ${tab.toLowerCase()} by name`, async function (assert) {
      const backend = `keymgmt-${uuidv4()}`;
      const names = ['test-item', 'other-item'];
      const filterInput = '[data-test-component="navigate-input"]';

      if (tab === 'Providers') {
        this.server.get(`/${backend}/kms`, () => ({ data: { keys: names } }));
        for (const name of names) {
          this.server.get(`/${backend}/kms/${name}`, () => ({
            data: { provider: 'azurekeyvault', key_collection: 'test-keyvault' },
          }));
        }
      }

      await runCmd(mountEngineCmd('keymgmt', backend));
      try {
        if (tab === 'Keys') {
          for (const name of names) {
            await runCmd(`write ${backend}/key/${name} type=rsa-2048`);
          }
        }
        await visit(`/vault/secrets-engines/${backend}/list`);
        await click(GENERAL.secretTab(tab));
        assert.dom(SES.secretLink('test-item')).exists('The matching item appears before filtering');
        assert.dom(SES.secretLink('other-item')).exists('The other item appears before filtering');

        await fillIn(filterInput, 'test-item');
        assert.dom(SES.secretLink('test-item')).exists('An exact name match remains visible');
        assert.dom(SES.secretLink('other-item')).doesNotExist('Nonmatching items are hidden');

        await fillIn(filterInput, 'OTHER');
        assert.dom(SES.secretLink('other-item')).exists('Partial matching is case-insensitive');
        assert.dom(SES.secretLink('test-item')).doesNotExist('The previous match is no longer visible');

        await fillIn(filterInput, 'missing-item');
        assert.dom(SES.secretLink()).doesNotExist('No items appear for an unmatched filter');
        assert
          .dom(GENERAL.emptyStateTitle)
          .includesText(`There are no ${tab.toLowerCase()} matching "missing-item".`);

        await fillIn(filterInput, '');
        assert.dom(SES.secretLink('test-item')).exists('Clearing the filter restores the matching item');
        assert.dom(SES.secretLink('other-item')).exists('Clearing the filter restores the other item');
      } finally {
        await runCmd(`delete sys/mounts/${backend}`);
      }
    });
  }

  test('it should display keymgmt configuration and tune keymgmt in the general settings form', async function (assert) {
    await consoleComponent.runCommands([
      // delete any previous mount with same name
      'delete sys/mounts/keymgmt',
    ]);
    const keymgmtType = 'keymgmt';
    await mountSecrets.visit();
    await settled();
    await mountBackend(keymgmtType, keymgmtType, true);
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText(`${keymgmtType} configuration`, 'displays configuration title');
    assert.dom(GENERAL.tab('general-settings')).hasText(`General settings`);
    assert.dom(GENERAL.cardContainer('version')).exists('version card exists');
    assert.dom(GENERAL.infoRowValue('type')).hasText(keymgmtType, 'shows keymgmt engine type');
    assert.dom(GENERAL.cardContainer('metadata')).exists('metadata card exists');
    assert.dom(GENERAL.copySnippet('path')).hasText(`${keymgmtType}/`, 'show path value');
    assert.dom(GENERAL.cardContainer('secrets duration')).exists('secrets duration card exists');
    assert.dom(GENERAL.cardContainer('security')).exists('security card exists');

    // fill in values to tune
    await fillIn(GENERAL.inputByAttr('default_lease_ttl'), 10);
    await fillIn(GENERAL.selectByAttr('default_lease_ttl'), 'm');
    await fillIn(GENERAL.inputByAttr('max_lease_ttl'), 15);
    await fillIn(GENERAL.selectByAttr('max_lease_ttl'), 'd');
    await fillIn(GENERAL.textareaByAttr('description'), 'Some awesome description.');
    await click(GENERAL.submitButton);

    // after submitting go to list and back to configuration
    await visit(`/vault/secrets-engines/${keymgmtType}/list`);
    await visit(`/vault/secrets-engines/${keymgmtType}/configuration/general-settings`);

    // confirm that submitted values were saved and prepopulated with those saved values
    assert
      .dom(GENERAL.textareaByAttr('description'))
      .hasValue('Some awesome description.', 'description was tuned');
    assert.dom(GENERAL.inputByAttr('default_lease_ttl')).hasValue('10', 'default ttl value was tuned');
    assert
      .dom(GENERAL.selectByAttr('default_lease_ttl'))
      .hasValue('m', 'default ttl value was tuned and shows correct unit');
    assert.dom(GENERAL.inputByAttr('max_lease_ttl')).hasValue('15', 'max ttl value was tuned');
    assert
      .dom(GENERAL.selectByAttr('max_lease_ttl'))
      .hasValue('d', 'max ttl value was tuned and shows correct unit');

    // show unsaved changes modal and save
    await fillIn(GENERAL.inputByAttr('default_lease_ttl'), 11);
    await fillIn(GENERAL.selectByAttr('default_lease_ttl'), 'm');
    await fillIn(GENERAL.textareaByAttr('description'), 'Updated awesome description.');
    await click(GENERAL.breadcrumbAtIdx(2));
    assert.dom(GENERAL.modal.container('unsaved-changes')).exists('Unsaved changes exists');
    await click(GENERAL.button('save'));
    await waitUntil(() => currentRouteName() === 'vault.cluster.secrets.backend.list-root');
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.list-root');
    await visit(`/vault/secrets-engines/${keymgmtType}/configuration/general-settings`);

    assert
      .dom(GENERAL.textareaByAttr('description'))
      .hasValue('Updated awesome description.', 'description was tuned from unsaved changes modal');
    assert.dom(GENERAL.inputByAttr('default_lease_ttl')).hasValue('11', 'default ttl value was tuned');
    assert.dom(GENERAL.selectByAttr('default_lease_ttl')).hasValue('m', 'default ttl unit was tuned');

    await visit(`/vault/secrets-engines/${keymgmtType}/configuration/general-settings`);

    // show unsaved changes modal and discard
    await fillIn(GENERAL.inputByAttr('default_lease_ttl'), 12);
    await fillIn(GENERAL.selectByAttr('default_lease_ttl'), 'm');
    await fillIn(GENERAL.textareaByAttr('description'), 'Some awesome description.');

    await click(GENERAL.breadcrumbAtIdx(2));
    assert.dom(GENERAL.modal.container('unsaved-changes')).exists('Unsaved changes exists');
    await click(GENERAL.button('discard'));
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.list-root');
    await visit(`/vault/secrets-engines/${keymgmtType}/configuration/general-settings`);

    assert
      .dom(GENERAL.textareaByAttr('description'))
      .hasValue(
        'Updated awesome description.',
        'description was reset to original values after discarding from unsaved changes modal'
      );
    assert
      .dom(GENERAL.inputByAttr('default_lease_ttl'))
      .hasValue('11', 'default ttl value was reset to original values');
    assert
      .dom(GENERAL.selectByAttr('default_lease_ttl'))
      .hasValue('m', 'default ttl unit was reset to original values');

    // navigate to secrets engines list and delete from the table row icon menu
    await visit('/vault/secrets-engines');
    await fillIn(GENERAL.inputSearch('secret-engine-path'), keymgmtType);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete-engine-path'));
    await fillIn(GENERAL.confirmTextInput, 'delete-engine');
    await click(GENERAL.confirmButton);

    await consoleComponent.runCommands([
      // cleanup after
      'delete sys/mounts/keymgmt',
    ]);
  });
});
