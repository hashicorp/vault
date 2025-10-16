/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/index';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { v4 as uuidv4 } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';

const SELECTORS = {
  listItem: (name) => `[data-test-identity-row="${name}"]`,
};

module('Acceptance | /access/identity/entities', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return login();
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

    await click(`${SELECTORS.listItem(name)} ${GENERAL.menuTrigger}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Create alias Edit Disable Delete', 'all actions render for entities');
    await click(`${SELECTORS.listItem(name)} ${GENERAL.menuItem('delete')}`);
    await click(GENERAL.confirmButton);
  });

  test('it renders popup menu for external groups', async function (assert) {
    const name = `external-${uuidv4()}`;
    await runCmd(`vault write identity/group name="${name}" policies="default" type="external"`);
    await visit('/vault/access/identity/groups');
    assert.strictEqual(currentURL(), '/vault/access/identity/groups', 'navigates to the groups tab');

    await click(`${SELECTORS.listItem(name)} ${GENERAL.menuTrigger}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Create alias Edit Delete', 'all actions render for external groups');
    await click(`${SELECTORS.listItem(name)} ${GENERAL.menuItem('delete')}`);
    await click(GENERAL.confirmButton);
  });

  test('it renders popup menu for external groups with alias', async function (assert) {
    const groupId = '44b2f1d1-699a-4a79-3a7b-37e53e17e7b2';
    const groupName = 'external-hasalias';
    // only relevant response keys are stubbed to simplify testing (more data is actually returned by both endpoints)
    this.server.get('/identity/group/id', () => {
      return {
        data: {
          key_info: { [groupId]: { name: groupName } },
          keys: [groupId],
        },
      };
    });

    this.server.get(`/identity/group/id/${groupId}`, () => {
      return {
        data: {
          alias: { id: '15bac764-d690-b72a-9cbc-b1fdeac1af9e', name: 'alias-test' },
          type: 'external',
        },
      };
    });

    await visit('/vault/access/identity/groups');
    await click(`${SELECTORS.listItem(groupName)} ${GENERAL.menuTrigger}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Edit Delete', 'no "Create alias" option for external groups with an alias');
  });

  test('it renders popup menu for internal groups', async function (assert) {
    const name = `internal-${uuidv4()}`;
    await runCmd(`vault write identity/group name="${name}" policies="default" type="internal"`);
    await visit('/vault/access/identity/groups');
    await click(`${SELECTORS.listItem(name)} ${GENERAL.menuTrigger}`);
    assert
      .dom('.hds-dropdown ul')
      .hasText('Details Edit Delete', 'no "Create alias" option for internal groups');
    await click(`${SELECTORS.listItem(name)} ${GENERAL.menuItem('delete')}`);
    await click(GENERAL.confirmButton);
  });
});
