/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, click, findAll, visit, fillIn } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountAuthCmd, runCmd } from 'vault/tests/helpers/commands';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

module('Acceptance | auth-methods list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    await login();
    // dismiss wizard
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.authMethods);
  });

  test('it navigates to auth method', async function (assert) {
    await visit('/vault/access/');
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.dom('[data-test-sidebar-nav-link="Authentication methods"]').hasClass('active');
  });

  test('it filters by name', async function (assert) {
    assert.expect(2);
    const authPath1 = `userpass-1-${this.uid}`;
    const type = 'userpass';
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath1));
    await visit('/vault/access/');

    // filter by name
    await fillIn(GENERAL.inputSearch('auth-method-name'), authPath1);
    assert.dom(GENERAL.tableData(0, 'path')).hasText(`${authPath1}/`);

    // clear filter by name
    await fillIn(GENERAL.inputSearch('auth-method-name'), '');
    const rows = findAll(GENERAL.tableRow());
    assert.true(rows.length > 1, 'filter has been removed');

    // cleanup
    await runCmd(`delete sys/auth/${authPath1}`);
  });

  test('it filters by auth type', async function (assert) {
    assert.expect(2);
    const authPath1 = `userpass-1-${this.uid}`;
    const authPath2 = `userpass-2-${this.uid}`;
    const type = 'userpass';
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath1));
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath2));
    await visit('/vault/access/');

    // filter by auth type
    await click(GENERAL.toggleInput('filter-by-auth-type'));
    await click(GENERAL.checkboxByAttr('userpass'));

    const rows = document.querySelectorAll(GENERAL.tableRow());
    const rowsAws = Array.from(rows).filter((row) => row.innerText.includes('userpass'));
    assert.strictEqual(rows.length, rowsAws.length, 'all rows returned are userpass');

    // clear filter
    await click(GENERAL.button('Clear all'));
    const rowsAgain = document.querySelectorAll(GENERAL.tableRow());
    assert.true(rowsAgain.length > 2, 'filter has been removed');

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
      assert.dom(GENERAL.listItem(`${key}`)).exists({ count: 1 }, `auth method ${key} appears in list view`);
    }
    await visit('/vault/settings/auth/enable');
    await click(GENERAL.navLink('OIDC provider'));
    await visit('/vault/access/');
    for (const [key] of Object.entries(authPayload)) {
      assert
        .dom(GENERAL.listItem(`${key}`))
        .exists({ count: 1 }, `auth method ${key} appears in list view after navigating from OIDC provider`);
    }
  });

  test('it should disable an auth method', async function (assert) {
    const authPath1 = `userpass-1-${this.uid}`;
    const type = 'userpass';
    await visit('/vault/settings/auth/enable');
    await runCmd(mountAuthCmd(type, authPath1));
    await visit('/vault/access/');
    await click(`${GENERAL.listItem(`${authPath1}/`)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.button('Disable auth method'));
    await click(GENERAL.confirmButton);
    assert.dom(GENERAL.listItem(`${authPath1}/`)).doesNotExist('auth mount disabled');
    await runCmd(`delete sys/auth/${authPath1}`);
  });
});
