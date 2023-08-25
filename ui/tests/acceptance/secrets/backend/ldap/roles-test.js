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
import { click } from '@ember/test-helpers';
import { isURL, visitURL } from 'vault/tests/helpers/ldap';

module('Acceptance | ldap | roles', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'ldap';
  });

  hooks.beforeEach(async function () {
    ldapMirageScenario(this.server);
    await authPage.login();
    return visitURL('roles');
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should transition to create role route on toolbar link click', async function (assert) {
    await click('[data-test-toolbar-action="role"]');
    assert.true(isURL('roles/create'), 'Transitions to role create route on toolbar link click');
  });

  test('it should transition to role details route on list item click', async function (assert) {
    await click('[data-test-list-item-link]:nth-of-type(1) a');
    assert.true(
      isURL('roles/dynamic/dynamic-role/details'),
      'Transitions to role details route on list item click'
    );

    await click('[data-test-breadcrumb="roles"]');
    await click('[data-test-list-item-link]:nth-of-type(2) a');
    assert.true(
      isURL('roles/static/static-role/details'),
      'Transitions to role details route on list item click'
    );
  });

  test('it should transition to routes from list item action menu', async function (assert) {
    assert.expect(3);

    for (const action of ['edit', 'get-creds', 'details']) {
      await click('[data-test-popup-menu-trigger]');
      await click(`[data-test-${action}]`);
      const uri = action === 'get-creds' ? 'credentials' : action;
      assert.true(
        isURL(`roles/dynamic/dynamic-role/${uri}`),
        `Transitions to ${uri} route on list item action menu click`
      );
      await click('[data-test-breadcrumb="roles"]');
    }
  });

  test('it should transition to routes from role details toolbar links', async function (assert) {
    await click('[data-test-list-item-link]:nth-of-type(1) a');
    await click('[data-test-get-credentials]');
    assert.true(
      isURL('roles/dynamic/dynamic-role/credentials'),
      'Transitions to credentials route from toolbar link'
    );

    await click('[data-test-breadcrumb="dynamic-role"]');
    await click('[data-test-edit]');
    assert.true(isURL('roles/dynamic/dynamic-role/edit'), 'Transitions to edit route from toolbar link');
  });
});
