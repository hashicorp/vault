/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint qunit/no-conditional-assertions: "warn" */
import { click, fillIn, settled, visit, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { supportedManagedAuthBackends } from 'vault/helpers/supported-managed-auth-backends';
import { deleteAuthCmd, mountAuthCmd, runCmd } from 'vault/tests/helpers/commands';

const SELECTORS = {
  backendLink: (path) => `[data-test-auth-backend-link="${path}"]`,
  createUser: '[data-test-entity-create-link="user"]',
  input: (attr) => `[data-test-input="${attr}"]`,
  password: '[data-test-textarea]',
  saveBtn: '[data-test-save-config]',
  methods: '[data-test-access-methods]',
  listItem: '[data-test-list-item-content]',
};
module('Acceptance | auth backend list', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    this.path1 = `userpass-${uuidv4()}`;
    this.path2 = `userpass-${uuidv4()}`;
    this.user1 = 'user1';
    this.user2 = 'user2';

    await runCmd(mountAuthCmd('userpass', this.path1));
    await runCmd(mountAuthCmd('userpass', this.path2));
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd(deleteAuthCmd(this.path1));
    await runCmd(deleteAuthCmd(this.path2));
    return;
  });

  test('userpass secret backend', async function (assert) {
    assert.expect(5);
    // enable a user in first userpass backend
    await visit('/vault/access');
    await click(SELECTORS.backendLink(this.path1));
    await click(SELECTORS.createUser);
    await fillIn(SELECTORS.input('username'), this.user1);
    await fillIn(SELECTORS.password, this.user1);
    await click(SELECTORS.saveBtn);
    assert.strictEqual(currentURL(), `/vault/access/${this.path1}/item/user`);

    await click(SELECTORS.methods);
    assert.strictEqual(currentURL(), '/vault/access');

    // enable a user in second userpass backend
    await click(SELECTORS.backendLink(this.path2));
    await click(SELECTORS.createUser);
    await fillIn(SELECTORS.input('username'), this.user2);
    await fillIn(SELECTORS.password, this.user2);
    await click(SELECTORS.saveBtn);
    assert.strictEqual(currentURL(), `/vault/access/${this.path2}/item/user`);

    // Confirm that the user was created. There was a bug where the apiPath was not being updated when toggling between auth routes.
    assert.dom(SELECTORS.listItem).hasText(this.user2, 'user2 exists in the list');

    // Confirm that the auth method 1 shows user1. There was a bug where the user was not listed when toggling between auth routes.
    await click(SELECTORS.methods);
    await click(SELECTORS.backendLink(this.path1));
    assert.dom(SELECTORS.listItem).hasText(this.user1, 'user1 exists in the list');
  });

  test('auth methods are linkable and link to correct view', async function (assert) {
    assert.expect(16);
    const uid = uuidv4();
    await visit('/vault/access');

    const supportManaged = supportedManagedAuthBackends();
    const backends = supportedAuthBackends();
    for (const backend of backends) {
      const { type } = backend;
      const path = `auth-list-${type}-${uid}`;
      if (type !== 'token') {
        await enablePage.enable(type, path);
      }
      await settled();
      await visit('/vault/access');

      // all auth methods should be linkable
      await click(`[data-test-auth-backend-link="${type === 'token' ? type : path}"]`);
      if (!supportManaged.includes(type)) {
        assert.dom('[data-test-auth-section-tab]').exists({ count: 1 });
        assert
          .dom('[data-test-auth-section-tab]')
          .hasText('Configuration', `only shows configuration tab for ${type} auth method`);
        assert.dom('[data-test-doc-link] .doc-link').exists(`includes doc link for ${type} auth method`);
      } else {
        let expectedTabs = 2;
        if (type == 'ldap' || type === 'okta') {
          expectedTabs = 3;
        }
        assert
          .dom('[data-test-auth-section-tab]')
          .exists({ count: expectedTabs }, `has management tabs for ${type} auth method`);
        // cleanup method
        await runCmd(deleteAuthCmd(path));
      }
    }
  });
});
