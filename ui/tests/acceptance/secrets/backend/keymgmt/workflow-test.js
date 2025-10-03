/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { settled, click, fillIn, visit, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SELECTORS } from 'vault/tests/helpers/secret-engine/general-settings-selectors';
import { create } from 'ember-cli-page-object';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

module('Acceptance | Enterprise | keymgmt-configuration-workflow', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('it should display keymgmt configuration and tune keymgmt in the general settings form', async function (assert) {
    await consoleComponent.runCommands([
      // delete any previous mount with same name
      'delete sys/mounts/keymgmt',
    ]);
    const keymgmtType = 'keymgmt';
    await mountSecrets.visit();
    await settled();
    await mountBackend(keymgmtType, keymgmtType);
    await click(SELECTORS.manageDropdown);
    await click(SELECTORS.manageDropdownItem('Configure'));
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText(`${keymgmtType} configuration`, 'displays configuration title');
    assert.dom(GENERAL.tab('general-settings')).hasText(`General settings`);
    assert.dom(GENERAL.cardContainer('version')).exists('version card exists');
    assert.dom(SELECTORS.versionCard.engineType).hasText(keymgmtType, 'shows keymgmt engine type');
    assert.dom(GENERAL.cardContainer('metadata')).exists('metadata card exists');
    assert.dom(GENERAL.inputByAttr('path')).hasValue(`${keymgmtType}/`, 'show path value');
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
    await visit(`/vault/secrets/${keymgmtType}/list`);
    await visit(`/vault/secrets/${keymgmtType}/configuration`);

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

    // navigate back to keymgmt list view to delete the engine from the manage dropdown
    await visit(`/vault/secrets/${keymgmtType}/list`);
    await click(SELECTORS.manageDropdown);
    await click(SELECTORS.manageDropdownItem('Delete'));
    await click(GENERAL.confirmButton);
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backends');
    await consoleComponent.runCommands([
      // cleanup after
      'delete sys/mounts/keymgmt',
    ]);
  });
});
