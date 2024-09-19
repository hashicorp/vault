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
import { click, fillIn, waitFor } from '@ember/test-helpers';
import { isURL, visitURL } from 'vault/tests/helpers/ldap/ldap-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';

module('Acceptance | ldap | roles', function (hooks) {
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
    return visitURL('roles', this.backend);
  });

  hooks.afterEach(async function () {
    await runCmd(deleteEngineCmd(this.backend));
  });

  test('it should transition to create role route on toolbar link click', async function (assert) {
    await click('[data-test-toolbar-action="role"]');
    assert.true(
      isURL('roles/create', this.backend),
      'Transitions to role create route on toolbar link click'
    );
  });

  test('it should transition to role details route on list item click', async function (assert) {
    await click('[data-test-list-item-link]:nth-of-type(1) a');
    assert.true(
      isURL('roles/dynamic/dynamic-role/details', this.backend),
      'Transitions to role details route on list item click'
    );

    await click('[data-test-breadcrumb="roles"] a');
    await click('[data-test-list-item-link]:nth-of-type(2) a');
    assert.true(
      isURL('roles/static/static-role/details', this.backend),
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
        isURL(`roles/dynamic/dynamic-role/${uri}`, this.backend),
        `Transitions to ${uri} route on list item action menu click`
      );
      await click('[data-test-breadcrumb="roles"] a');
    }
  });

  test('it should transition to routes from role details toolbar links', async function (assert) {
    await click('[data-test-list-item-link]:nth-of-type(1) a');
    await click('[data-test-get-credentials]');
    assert.true(
      isURL('roles/dynamic/dynamic-role/credentials', this.backend),
      'Transitions to credentials route from toolbar link'
    );

    await click('[data-test-breadcrumb="dynamic-role"] a');
    await click('[data-test-edit]');
    assert.true(
      isURL('roles/dynamic/dynamic-role/edit', this.backend),
      'Transitions to edit route from toolbar link'
    );
  });

  test('it should clear roles page filter value on route exit', async function (assert) {
    await fillIn('[data-test-filter-input]', 'foo');
    assert
      .dom('[data-test-filter-input]')
      .hasValue('foo', 'Roles page filter value set after model refresh and rerender');
    await waitFor(GENERAL.emptyStateTitle);
    await click('[data-test-tab="libraries"]');
    await click('[data-test-tab="roles"]');
    assert.dom('[data-test-filter-input]').hasNoValue('Roles page filter value cleared on route exit');
  });
});
