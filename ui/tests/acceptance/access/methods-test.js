/**
 * Copyright (c) HashiCorp, Inc.
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

const { searchSelect } = GENERAL;

module('Acceptance | auth-methods list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it navigates to auth method', async function (assert) {
    await visit('/vault/access/');
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.dom('[data-test-sidebar-nav-link="Authentication Methods"]').hasClass('active');
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
    let rows = findAll('[data-test-auth-backend-link]');
    const rowsUserpass = Array.from(rows).filter((row) => row.innerText.includes('userpass'));

    assert.strictEqual(rows.length, rowsUserpass.length, 'all rows returned are userpass');

    // filter by name
    await clickTrigger('#filter-by-auth-name');
    await click(searchSelect.option());
    const selectedItem = find(`#filter-by-auth-name ${searchSelect.selectedOption()}`).innerText;
    const singleRow = findAll('[data-test-auth-backend-link]');

    assert.strictEqual(singleRow.length, 1, 'returns only one row');
    assert.dom(singleRow[0]).includesText(selectedItem, 'shows the filtered by auth name');
    // clear filter by name
    await click(`#filter-by-auth-name ${searchSelect.removeSelected}`);
    rows = findAll('[data-test-auth-backend-link]');
    assert.true(rows.length > 1, 'filter has been removed');

    // cleanup
    await runCmd(`delete sys/auth/${authPath1}`);
    await runCmd(`delete sys/auth/${authPath2}`);
  });

  test('it should show all methods in list view', async function (assert) {
    this.server.get('/sys/internal/ui/mounts', () => ({
      data: {
        auth: {
          'token/': { accessor: 'auth_token_263b8b4e', type: 'token' },
          'userpass/': { accessor: 'auth_userpass_87aca1f8', type: 'userpass' },
        },
      },
    }));
    await visit('/vault/access/');
    assert.dom('[data-test-auth-backend-link]').exists({ count: 2 }, 'All auth methods appear in list view');

    // verify overflow style exists on auth method name
    assert.dom('[data-test-path]').hasClass('overflow-wrap', 'auth method name has overflow class applied');
    await visit('/vault/settings/auth/enable');
    await click('[data-test-sidebar-nav-link="OIDC Provider"]');
    await visit('/vault/access/');
    assert
      .dom('[data-test-auth-backend-link]')
      .exists({ count: 2 }, 'All auth methods appear in list view after navigating back');
  });
});
