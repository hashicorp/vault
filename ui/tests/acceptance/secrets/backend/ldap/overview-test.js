/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';
import ldapMirageScenario from 'vault/mirage/scenarios/ldap';
import ldapHandlers from 'vault/mirage/handlers/ldap';
import authPage from 'vault/tests/pages/auth';
import { click, visit } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { isURL, visitURL } from 'vault/tests/helpers/ldap/ldap-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';

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
    return authPage.login();
  });

  test('it should transition to ldap overview on mount success', async function (assert) {
    const backend = 'ldap-test-mount';
    await visit('/vault/secrets');
    await click('[data-test-enable-engine]');
    await mountBackend('ldap', backend);
    assert.true(isURL('overview', backend), 'Transitions to ldap overview route on mount success');
    assert.dom('[data-test-header-title]').hasText(backend);
    // cleanup mounted engine
    await visit('/vault/secrets');
    await runCmd(deleteEngineCmd(backend));
  });

  test('it should transition to routes on tab link click', async function (assert) {
    assert.expect(4);
    await this.mountAndConfig(this.backend);

    await visitURL('overview', this.backend);

    for (const tab of ['roles', 'libraries', 'config', 'overview']) {
      await click(`[data-test-tab="${tab}"]`);
      const route = tab === 'config' ? 'configuration' : tab;
      assert.true(isURL(route, this.backend), `Transitions to ${route} route on tab link click`);
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
  // including a test for the configuration route here since it is the only one needed
  test('it should transition to configuration edit on toolbar link click', async function (assert) {
    ldapMirageScenario(this.server);
    await this.mountAndConfig(this.backend);
    await visitURL('overview', this.backend);
    await click('[data-test-tab="config"]');
    await click('[data-test-toolbar-config-action]');
    assert.true(isURL('configure', this.backend), 'Transitions to configure route on toolbar link click');
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
