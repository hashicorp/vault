/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ldapMirageScenario from 'vault/mirage/scenarios/ldap';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { click, fillIn, visit } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { isURL, visitURL } from 'vault/tests/helpers/ldap';

module('Acceptance | ldap | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'ldap';
  });

  hooks.beforeEach(async function () {
    return authPage.login();
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should transition to ldap overview on mount success', async function (assert) {
    await visit('/vault/secrets');
    await click('[data-test-enable-engine]');
    await click('[data-test-mount-type="ldap"]');
    await click('[data-test-mount-next]');
    await fillIn('[data-test-input="path"]', 'ldap-test');
    await click('[data-test-mount-submit]');
    assert.true(isURL('overview'), 'Transitions to ldap overview route on mount success');
  });

  test('it should transition to routes on tab link click', async function (assert) {
    assert.expect(4);

    await visitURL('overview');

    for (const tab of ['roles', 'libraries', 'config', 'overview']) {
      await click(`[data-test-tab="${tab}"]`);
      const route = tab === 'config' ? 'configuration' : tab;
      assert.true(isURL(route), `Transitions to ${route} route on tab link click`);
    }
  });

  test('it should transition to configuration route when engine is not configured', async function (assert) {
    await visitURL('overview');
    await click('[data-test-config-cta] a');
    assert.true(isURL('configure'), 'Transitions to configure route on cta link click');

    await click('[data-test-breadcrumb="ldap-test"]');
    await click('[data-test-toolbar-action="config"]');
    assert.true(isURL('configure'), 'Transitions to configure route on toolbar link click');
  });
  // including a test for the configuration route here since it is the only one needed
  test('it should transition to configuration edit on toolbar link click', async function (assert) {
    ldapMirageScenario(this.server);
    await visitURL('overview');
    await click('[data-test-tab="config"]');
    await click('[data-test-toolbar-config-action]');
    assert.true(isURL('configure'), 'Transitions to configure route on toolbar link click');
  });

  test('it should transition to create role route on card action link click', async function (assert) {
    ldapMirageScenario(this.server);
    await visitURL('overview');
    await click('[data-test-overview-card="Roles"] a');
    assert.true(isURL('roles/create'), 'Transitions to role create route on card action link click');
  });

  test('it should transition to create library route on card action link click', async function (assert) {
    ldapMirageScenario(this.server);
    await visitURL('overview');
    await click('[data-test-overview-card="Libraries"] a');
    assert.true(isURL('libraries/create'), 'Transitions to library create route on card action link click');
  });

  test('it should transition to role credentials route on generate credentials action', async function (assert) {
    ldapMirageScenario(this.server);
    await visitURL('overview');
    await selectChoose('.search-select', 'static-role');
    await click('[data-test-generate-credential-button]');
    assert.true(
      isURL('roles/static/static-role/credentials'),
      'Transitions to role credentials route on generate credentials action'
    );

    await click('[data-test-breadcrumb="ldap-test"]');
    await selectChoose('.search-select', 'dynamic-role');
    await click('[data-test-generate-credential-button]');
    assert.true(
      isURL('roles/dynamic/dynamic-role/credentials'),
      'Transitions to role credentials route on generate credentials action'
    );
  });
});
