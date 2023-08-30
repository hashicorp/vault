/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import page from 'vault/tests/pages/access/methods';
import authEnable from 'vault/tests/pages/settings/auth/enable';
import authPage from 'vault/tests/pages/auth';
import ss from 'vault/tests/pages/components/search-select';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

import { v4 as uuidv4 } from 'uuid';

const consoleComponent = create(consoleClass);
const searchSelect = create(ss);

module('Acceptance | auth-methods list view', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it navigates to auth method', async function (assert) {
    await page.visit();
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.ok(page.methodsLink.isActive, 'the first link is active');
    assert.strictEqual(page.methodsLink.text, 'Authentication Methods');
  });

  test('it filters by name and auth type', async function (assert) {
    assert.expect(4);
    const authPath1 = `userpass-1-${this.uid}`;
    const authPath2 = `userpass-2-${this.uid}`;
    const type = 'userpass';
    await authEnable.visit();
    await authEnable.enable(type, authPath1);
    await authEnable.visit();
    await authEnable.enable(type, authPath2);
    await page.visit();
    // filter by auth type

    await clickTrigger('#filter-by-auth-type');
    await searchSelect.options.objectAt(0).click();

    const rows = document.querySelectorAll('[data-test-auth-backend-link]');
    const rowsUserpass = Array.from(rows).filter((row) => row.innerText.includes('userpass'));

    assert.strictEqual(rows.length, rowsUserpass.length, 'all rows returned are userpass');

    // filter by name
    await clickTrigger('#filter-by-auth-name');
    const firstItemToSelect = searchSelect.options.objectAt(0).text;
    await searchSelect.options.objectAt(0).click();
    const singleRow = document.querySelectorAll('[data-test-auth-backend-link]');

    assert.strictEqual(singleRow.length, 1, 'returns only one row');
    assert.dom(singleRow[0]).includesText(firstItemToSelect, 'shows the filtered by auth name');
    // clear filter by engine name
    await searchSelect.deleteButtons.objectAt(1).click();
    const rowsAgain = document.querySelectorAll('[data-test-auth-backend-link]');
    assert.ok(rowsAgain.length > 1, 'filter has been removed');

    // cleanup
    await consoleComponent.runCommands([`delete sys/auth/${authPath1}`]);
    await consoleComponent.runCommands([`delete sys/auth/${authPath2}`]);
  });
});
