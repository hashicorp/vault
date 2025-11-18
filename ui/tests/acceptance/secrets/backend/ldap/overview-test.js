/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';
import ldapMirageScenario from 'vault/mirage/scenarios/ldap';
import ldapHandlers from 'vault/mirage/handlers/ldap';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { click, visit, currentURL } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { isURL, visitURL } from 'vault/tests/helpers/ldap/ldap-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import secretsNavTestHelper from 'vault/tests/acceptance/secrets/secrets-nav-test-helper';

module('Acceptance | ldap | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    ldapHandlers(this.server);
    this.backend = `ldap-test-${uuidv4()}`;
    this.mountAndConfig = (backend) => {
      return runCmd([
        mountEngineCmd('ldap', backend),
        `write ${backend}/config binddn=foo bindpass=bar url=http://localhost:8208`,
      ]);
    };
    return login();
  });

  secretsNavTestHelper(test, 'ldap');

  test('it should transition to ldap overview on mount success', async function (assert) {
    const backend = 'ldap-test-mount';
    await visit('/vault/secrets-engines');
    await click('[data-test-enable-engine]');
    await mountBackend('ldap', backend);
    assert.true(isURL('overview', backend), 'Transitions to ldap overview route on mount success');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(backend);
    // cleanup mounted engine
    await visit('/vault/secrets-engines');
    await runCmd(deleteEngineCmd(backend));
  });

  test('it should transition to routes on tab link click', async function (assert) {
    assert.expect(3);
    await this.mountAndConfig(this.backend);

    await visitURL('overview', this.backend);

    for (const tab of ['roles', 'libraries', 'overview']) {
      await click(`[data-test-tab="${tab}"]`);
      assert.true(isURL(tab, this.backend), `Transitions to ${tab} route on tab link click`);
    }
  });

  test('it should transition to configuration route when engine is not configured', async function (assert) {
    await runCmd(mountEngineCmd('ldap', this.backend));
    await visitURL('overview', this.backend);
    await click('[data-test-config-cta] a');
    assert.true(isURL('configure', this.backend), 'Transitions to configure route on cta link click');

    await click(`[data-test-breadcrumb="${this.backend}"] a`);
    await click('[data-test-toolbar-action="config"]');
    assert.true(isURL('configure', this.backend), 'Transitions to configure route on toolbar link click');
  });

  test('it should transition to configuration edit on empty state click', async function (assert) {
    ldapMirageScenario(this.server);
    await runCmd(mountEngineCmd('ldap', this.backend));
    await visitURL('overview', this.backend);
    await click('.hds-link-standalone');
    assert.true(isURL('configure', this.backend), 'Transitions to configure route on empty state click');
  });

  test('it should delete the ldap engine on delete action', async function (assert) {
    ldapMirageScenario(this.server);
    await runCmd(mountEngineCmd('ldap', this.backend));
    await visitURL('overview', this.backend);
    await click(GENERAL.manageDropdown);
    await click(GENERAL.manageDropdownItem('Delete'));
    assert.dom('[data-test-confirm-modal]').exists('Confirm delete modal renders');
    await click('[data-test-confirm-button]');
    assert.strictEqual(
      currentURL(),
      '/vault/secrets-engines',
      'navigates back to the secrets engines list after engine deletion'
    );
  });

  test('it should transition to create role route on card action link click', async function (assert) {
    ldapMirageScenario(this.server);
    await this.mountAndConfig(this.backend);
    await visitURL('overview', this.backend);
    await click('[data-test-overview-card="Roles"] a');
    assert.true(
      isURL('roles/create', this.backend),
      'Transitions to role create route on card action link click'
    );
  });

  test('it should transition to create library route on card action link click', async function (assert) {
    ldapMirageScenario(this.server);
    await this.mountAndConfig(this.backend);
    await visitURL('overview', this.backend);
    await click('[data-test-overview-card="Libraries"] a');
    assert.true(
      isURL('libraries/create', this.backend),
      'Transitions to library create route on card action link click'
    );
  });

  test('it should transition to role credentials route on generate credentials action', async function (assert) {
    ldapMirageScenario(this.server);
    await this.mountAndConfig(this.backend);
    await visitURL('overview', this.backend);
    await selectChoose('.search-select', 'static-role');
    await click('[data-test-generate-credential-button]');
    assert.true(
      isURL('roles/static/static-role/credentials', this.backend),
      'Transitions to role credentials route on generate credentials action'
    );

    await click(`[data-test-breadcrumb="${this.backend}"] a`);
    await selectChoose('.search-select', 'dynamic-role');
    await click('[data-test-generate-credential-button]');
    assert.true(
      isURL('roles/dynamic/dynamic-role/credentials', this.backend),
      'Transitions to role credentials route on generate credentials action'
    );
  });
});
