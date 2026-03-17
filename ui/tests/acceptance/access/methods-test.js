/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, click, find, findAll, visit } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountAuthCmd, runCmd } from 'vault/tests/helpers/commands';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { sanitizePath } from 'core/utils/sanitize-path';
import localStorage from 'vault/lib/local-storage';

const { searchSelect } = GENERAL;

module('Acceptance | auth-methods list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    await login();
    // dismiss wizard
    localStorage.setItem('dismissed-wizards', ['auth-methods']);
  });

  test('it navigates to auth method', async function (assert) {
    await visit('/vault/access/');
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.dom('[data-test-sidebar-nav-link="Authentication methods"]').hasClass('active');
  });

  test('it filters by name and auth type', async function (assert) {
    assert.expect(4);
    const authPath1 = `userpass-1-${this.uid}`;
    const authPath2 = `userpass-2-${this.uid}`;
    const type = 'userpass';
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath1));
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath2));
    await visit('/vault/access/');

    // filter by auth type
    await clickTrigger('#filter-by-auth-type');
    await click(searchSelect.option(searchSelect.optionIndex(type)));
    let rows = findAll('.list-item-row');
    const rowsUserpass = findAll(GENERAL.button('userpass'));
    assert.strictEqual(rows.length, rowsUserpass.length, 'all rows returned are userpass');

    // filter by name
    await clickTrigger('#filter-by-auth-name');
    await click(searchSelect.option());
    const selectedItem = find(`#filter-by-auth-name ${searchSelect.selectedOption()}`).innerText;
    const singleRow = findAll('.linked-block');
    assert.strictEqual(singleRow.length, 1, 'returns only one row');
    assert.dom(singleRow[0]).includesText(selectedItem, 'shows the filtered by auth name');
    // clear filter by name
    await click(`#filter-by-auth-name ${searchSelect.removeSelected}`);
    rows = findAll('.linked-block');
    assert.true(rows.length > 1, 'filter has been removed');

    // cleanup
    await runCmd(`delete sys/auth/${authPath1}`);
    await runCmd(`delete sys/auth/${authPath2}`);
  });

  test('it should show all methods in list view', async function (assert) {
    const authPayload = {
      'token/': { accessor: 'auth_token_263b8b4e', type: 'token' },
      'userpass/': { accessor: 'auth_userpass_87aca1f8', type: 'userpass' },
    };
    this.server.get('/sys/internal/ui/mounts', () => ({
      data: {
        auth: authPayload,
      },
    }));
    await visit('/vault/access/');
    for (const [key] of Object.entries(authPayload)) {
      assert
        .dom(GENERAL.linkedBlock(sanitizePath(key)))
        .exists({ count: 1 }, `auth method ${key} appears in list view`);
    }
    await visit('/vault/settings/auth/enable');
    await click(GENERAL.navLink('OIDC provider'));
    await visit('/vault/access/');
    for (const [key] of Object.entries(authPayload)) {
      assert
        .dom(GENERAL.linkedBlock(sanitizePath(key)))
        .exists({ count: 1 }, `auth method ${key} appears in list view after navigating from OIDC provider`);
    }
  });

  test('it should disable an auth method', async function (assert) {
    const authPath1 = `userpass-1-${this.uid}`;
    const type = 'userpass';
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath1));
    await visit('/vault/access/');
    await click(`${GENERAL.linkedBlock(authPath1)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.button('Disable auth method'));
    await click(GENERAL.confirmButton);
    assert.dom(GENERAL.linkedBlock(authPath1)).doesNotExist('auth mount disabled');
    await runCmd(`delete sys/auth/${authPath1}`);
  });
});
