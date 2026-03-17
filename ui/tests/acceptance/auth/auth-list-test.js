/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, settled, visit, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login, loginNs } from 'vault/tests/helpers/auth/auth-helpers';
import { MANAGED_AUTH_BACKENDS } from 'vault/helpers/supported-managed-auth-backends';
import { deleteAuthCmd, mountAuthCmd, runCmd, createNS } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import localStorage from 'vault/lib/local-storage';

const SELECTORS = {
  createUser: '[data-test-entity-create-link="user"]',
  methods: '[data-test-access-methods] a',
  listItem: '[data-test-list-item-content]',
};

module('Acceptance | auth backend list', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    // dismiss wizard
    localStorage.setItem('dismissed-wizards', ['auth-methods']);
  });

  test('userpass secret backend', async function (assert) {
    this.path1 = `userpass-${uuidv4()}`;
    this.path2 = `userpass-${uuidv4()}`;
    this.user1 = 'user1';
    this.user2 = 'user2';

    await runCmd([mountAuthCmd('userpass', this.path1), mountAuthCmd('userpass', this.path2)], false);
    // helper function to create a user in the specified backend
    async function createUser(backendPath, username) {
      await click(GENERAL.linkedBlock(backendPath));
      assert.dom(GENERAL.emptyStateTitle).exists('shows empty state');
      await click(SELECTORS.createUser);
      await fillIn(GENERAL.inputByAttr('username'), username);
      await fillIn(GENERAL.inputByAttr('password'), username);
      await click(GENERAL.submitButton);
      assert.strictEqual(currentURL(), `/vault/access/${backendPath}/item/user`);
    }
    // visit access page and enable the first user in the first userpass backend
    await visit('/vault/access');
    await createUser(this.path1, this.user1);

    // navigate back to the methods list
    await click(GENERAL.breadcrumbAtIdx(0));
    assert.strictEqual(currentURL(), '/vault/access');

    // enable a second user in the second userpass backend
    await createUser(this.path2, this.user2);

    // verify the second user is listed after creation
    assert.dom(SELECTORS.listItem).hasText(this.user2, 'user2 exists in the list');

    // check that switching back to the first auth method shows the first user
    await click(GENERAL.breadcrumbAtIdx(0));
    await click(GENERAL.linkedBlock(this.path1));
    assert.dom(SELECTORS.listItem).hasText(this.user1, 'user1 exists in the list');

    await login();
    await runCmd([deleteAuthCmd(this.path1), deleteAuthCmd(this.path2)], false);
  });

  module('auth methods are linkable and link to correct view', function (hooks) {
    hooks.beforeEach(async function () {
      this.uid = uuidv4();
      await visit('/vault/access');
    });

    // Test all auth methods, not just those you can log in with
    filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: false })
      .map((backend) => backend.type)
      .forEach((type) => {
        test(`${type} auth method`, async function (assert) {
          const supportManaged = MANAGED_AUTH_BACKENDS;
          const isTokenType = type === 'token';
          const path = isTokenType ? 'token' : `auth-list-${type}-${this.uid}`;

          // Enable auth if the backend is not type token
          if (!isTokenType) {
            await visit('/vault/settings/auth/enable');
            await click(GENERAL.cardContainer(type));
            await fillIn(GENERAL.inputByAttr('path'), path);
            await click(GENERAL.submitButton);
          }

          await visit('/vault/access');

          // check popup menu for auth method
          const itemCount = isTokenType ? 2 : 3;
          const triggerSelector = `${GENERAL.linkedBlock(path)} ${GENERAL.menuTrigger}`;
          const itemSelector = `${GENERAL.linkedBlock(path)} .hds-dropdown-list-item`;

          await click(triggerSelector);

          assert
            .dom(itemSelector)
            .exists({ count: itemCount }, `shows ${itemCount} dropdown items for ${type}`);

          // check that auth methods are linkable
          await click(GENERAL.linkedBlock(path));

          if (!supportManaged.includes(type)) {
            assert.dom(GENERAL.linkTo('auth-tab')).exists({ count: 1 });
            assert
              .dom(GENERAL.linkTo('auth-tab'))
              .hasText('Configuration', `only shows configuration tab for ${type} auth method`);
            assert.dom(GENERAL.docLinkByAttr(path)).exists(`includes doc link for ${type} auth method`);
          } else {
            // determine expected number of managed auth tabs
            const expectedTabs = ['ldap', 'okta'].includes(type) ? 3 : 2;
            assert
              .dom(GENERAL.linkTo('generated-tab'))
              .exists({ count: expectedTabs }, `has management tabs for ${type} auth method`);
          }
          // cleanup method
          if (!isTokenType) {
            await runCmd(deleteAuthCmd(path));
          }
        });
      });
  });

  module('enterprise', function () {
    test('ent-only auth methods are linkable and link to correct view', async function (assert) {
      assert.expect(3);
      const uid = uuidv4();
      await visit('/vault/access');

      // Only SAML is enterprise-only for now
      const type = 'saml';
      const path = `auth-list-${type}-${uid}`;
      await runCmd([mountAuthCmd(type, path), 'refresh']);
      await settled();
      await visit('/vault/access');

      // all auth methods should be linkable
      await click(GENERAL.linkedBlock(path));
      assert.dom(GENERAL.linkTo('auth-tab')).exists({ count: 1 });
      assert
        .dom(GENERAL.linkTo('auth-tab'))
        .hasText('Configuration', `only shows configuration tab for ${type} auth method`);
      assert.dom(GENERAL.docLinkByAttr(path)).exists(`includes doc link for ${type} auth method`);
      await runCmd(deleteAuthCmd(path));
    });

    test('token config within namespace', async function (assert) {
      const ns = 'ns-wxyz';
      await runCmd(createNS(ns), false);
      await settled();
      await loginNs(ns);
      localStorage.setItem('dismissed-wizards', ['auth-methods']);

      // go directly to token configure route
      await visit(`/vault/settings/auth/configure/token/options?namespace=${ns}`);
      await fillIn(GENERAL.inputByAttr('description'), 'My custom description');
      await click(GENERAL.submitButton);
      assert.strictEqual(
        currentURL(),
        `/vault/access?namespace=${ns}`,
        'successfully saves and navigates away'
      );
      await click(GENERAL.linkedBlock('token'));
      assert
        .dom(GENERAL.infoRowValue('Description'))
        .hasText('My custom description', 'description was saved');
      await login();
      await runCmd(`delete sys/namespaces/${ns}`);
    });
  });
});
