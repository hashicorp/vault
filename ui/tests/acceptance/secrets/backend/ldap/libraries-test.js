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

module('Acceptance | ldap | libraries', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'ldap';
  });

  hooks.beforeEach(async function () {
    ldapMirageScenario(this.server);
    await authPage.login();
    return visitURL('libraries');
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should transition to create library route on toolbar link click', async function (assert) {
    await click('[data-test-toolbar-action="library"]');
    assert.true(isURL('libraries/create'), 'Transitions to library create route on toolbar link click');
  });

  test('it should transition to library details route on list item click', async function (assert) {
    await click('[data-test-list-item-link] a');
    assert.true(
      isURL('libraries/test-library/details/accounts'),
      'Transitions to library details accounts route on list item click'
    );
  });

  test('it should transition to routes from list item action menu', async function (assert) {
    assert.expect(2);

    for (const action of ['edit', 'details']) {
      await click('[data-test-popup-menu-trigger]');
      await click(`[data-test-${action}]`);
      const uri = action === 'details' ? 'details/accounts' : action;
      assert.true(
        isURL(`libraries/test-library/${uri}`),
        `Transitions to ${action} route on list item action menu click`
      );
      await click('[data-test-breadcrumb="libraries"]');
    }
  });

  test('it should transition to details routes from tab links', async function (assert) {
    await click('[data-test-list-item-link] a');
    await click('[data-test-tab="config"]');
    assert.true(
      isURL('libraries/test-library/details/configuration'),
      'Transitions to configuration route on tab click'
    );

    await click('[data-test-tab="accounts"]');
    assert.true(
      isURL('libraries/test-library/details/accounts'),
      'Transitions to accounts route on tab click'
    );
  });

  test('it should transition to routes from library details toolbar links', async function (assert) {
    await click('[data-test-list-item-link] a');
    await click('[data-test-edit]');
    assert.true(isURL('libraries/test-library/edit'), 'Transitions to credentials route from toolbar link');
  });
});
