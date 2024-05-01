/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { fillIn, click, currentRouteName, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/index';
import authPage from 'vault/tests/pages/auth';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { v4 as uuidv4 } from 'uuid';

const SELECTORS = {
  listItem: (name) => `[data-test-identity-row="${name}"]`,
  menu: `[data-test-popup-menu-trigger]`,
  menuItem: (element) => `[data-test-popup-menu="${element}"]`,
  submit: '[data-test-identity-submit]',
  confirm: '[data-test-confirm-button]',
};
module('Acceptance | /access/identity/entities', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it renders the entities page', async function (assert) {
    await page.visit({ item_type: 'entities' });
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.identity.index',
      'navigates to the correct route'
    );
  });

  test('it renders the groups page', async function (assert) {
    await page.visit({ item_type: 'groups' });
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.identity.index',
      'navigates to the correct route'
    );
  });

  test('it renders popup menu for entities', async function (assert) {
    const name = `entity-${uuidv4()}`;
    await runCmd(`vault write identity/entity name="${name}" policies="default"`);
    await visit('/vault/access/identity/entities');
    assert.strictEqual(currentURL(), '/vault/access/identity/entities', 'navigates to entities tab');

    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menu}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Create alias Edit Disable Delete', 'all actions render for entities');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menuItem('delete')}`);
    await click(SELECTORS.confirm);
  });

  test('it renders popup menu for external groups', async function (assert) {
    const name = `external-${uuidv4()}`;
    await runCmd(`vault write identity/group name="${name}" policies="default" type="external"`);
    await visit('/vault/access/identity/groups');
    assert.strictEqual(currentURL(), '/vault/access/identity/groups', 'navigates to the groups tab');

    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menu}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Create alias Edit Delete', 'all actions render for external groups');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menuItem('delete')}`);
    await click(SELECTORS.confirm);
  });

  test('it renders popup menu for external groups with alias', async function (assert) {
    const name = `external-hasalias-${uuidv4()}`;
    await runCmd(`vault write identity/group name="${name}" policies="default" type="external"`);
    await visit('/vault/access/identity/groups');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menu}`);
    await click(SELECTORS.menuItem('create alias'));
    await fillIn(GENERAL.inputByAttr('name'), 'alias-test');
    await click(SELECTORS.submit);

    await visit('/vault/access/identity/groups');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menu}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Edit Delete', 'no "Create alias" option for external groups with an alias');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menuItem('delete')}`);
    await click(SELECTORS.confirm);
  });

  test('it renders popup menu for internal groups', async function (assert) {
    const name = `internal-${uuidv4()}`;
    await runCmd(`vault write identity/group name="${name}" policies="default" type="internal"`);
    await visit('/vault/access/identity/groups');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menu}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Edit Delete', 'no "Create alias" option for internal groups');
    await click(`${SELECTORS.listItem(name)} ${SELECTORS.menuItem('delete')}`);
    await click(SELECTORS.confirm);
  });
});
