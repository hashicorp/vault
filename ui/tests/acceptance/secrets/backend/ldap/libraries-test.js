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
import { click, currentURL } from '@ember/test-helpers';
import { isURL, visitURL } from 'vault/tests/helpers/ldap/ldap-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { LDAP_SELECTORS } from 'vault/tests/helpers/ldap/ldap-selectors';

module('Acceptance | ldap | libraries', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    ldapHandlers(this.server);
    ldapMirageScenario(this.server);
    this.backend = `ldap-test-${uuidv4()}`;
    await authPage.login();
    // mount & configure
    await runCmd([
      mountEngineCmd('ldap', this.backend),
      `write ${this.backend}/config binddn=foo bindpass=bar url=http://localhost:8208`,
    ]);
    return visitURL('libraries', this.backend);
  });

  hooks.afterEach(async function () {
    await runCmd(deleteEngineCmd(this.backend));
  });

  test('it should show libraries on overview page', async function (assert) {
    await visitURL('overview', this.backend);
    assert.dom('[data-test-libraries-count]').hasText('2');
  });

  test('it should transition to create library route on toolbar link click', async function (assert) {
    await click('[data-test-toolbar-action="library"]');
    assert.true(
      isURL('libraries/create', this.backend),
      'Transitions to library create route on toolbar link click'
    );
  });

  test('it should transition to library details route on list item click', async function (assert) {
    await click(LDAP_SELECTORS.libraryItem('test-library'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/ldap/libraries/test-library/details/accounts`,
      'Transitions to library details accounts route on list item click'
    );
    assert.dom('[data-test-account-name]').exists({ count: 2 }, 'lists the accounts');
    assert.dom('[data-test-checked-out-account]').exists({ count: 1 }, 'lists the checked out accounts');
  });

  test('it should transition to library details for hierarchical list items', async function (assert) {
    await click(LDAP_SELECTORS.libraryItem('admin/'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/ldap/libraries/subdirectory/admin/`,
      'Transitions to subdirectory list view'
    );

    await click(LDAP_SELECTORS.libraryItem('admin/test-library'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/ldap/libraries/admin%2Ftest-library/details/accounts`,
      'Transitions to child library details accounts'
    );
    assert.dom('[data-test-account-name]').exists({ count: 2 }, 'lists the accounts');
    assert.dom('[data-test-checked-out-account]').exists({ count: 1 }, 'lists the checked out accounts');
  });

  test('it should transition to routes from list item action menu', async function (assert) {
    assert.expect(2);

    for (const action of ['edit', 'details']) {
      await click('[data-test-popup-menu-trigger]');
      await click(`[data-test-${action}]`);
      const uri = action === 'details' ? 'details/accounts' : action;
      assert.true(
        isURL(`libraries/test-library/${uri}`, this.backend),
        `Transitions to ${action} route on list item action menu click`
      );
      await click('[data-test-breadcrumb="Libraries"] a');
    }
  });

  test('it should transition to details routes from tab links', async function (assert) {
    await click('[data-test-list-item-link] a');
    await click('[data-test-tab="config"]');
    assert.true(
      isURL('libraries/test-library/details/configuration', this.backend),
      'Transitions to configuration route on tab click'
    );

    await click('[data-test-tab="accounts"]');
    assert.true(
      isURL('libraries/test-library/details/accounts', this.backend),
      'Transitions to accounts route on tab click'
    );
  });

  test('it should transition to routes from library details toolbar links', async function (assert) {
    await click('[data-test-list-item-link] a');
    await click('[data-test-edit]');
    assert.true(
      isURL('libraries/test-library/edit', this.backend),
      'Transitions to credentials route from toolbar link'
    );
  });
});
