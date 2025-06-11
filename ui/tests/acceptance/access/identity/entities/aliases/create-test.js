/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { click, currentURL, fillIn, settled, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasDeleteFromForm } from '../../_shared-alias-tests';
import { v4 as uuidv4 } from 'uuid';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | Entities | /access/identity/entities/aliases/add', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    await login();
    return;
  });

  test('Entities: it allows create, list, delete of an entity alias', async function (assert) {
    assert.expect(5);
    const name = `alias-${uuidv4()}`;
    await visit(`/vault/access/identity/entities/create`);
    await fillIn(GENERAL.inputByAttr('name'), name);
    await click(GENERAL.submitButton);
    const entityGeneratedId = document.querySelector(GENERAL.infoRowValue('ID')).innerText;
    // create entity alias
    await visit(`/vault/access/identity/entities/aliases/add/${entityGeneratedId}`);
    await fillIn(GENERAL.inputByAttr('name'), name);
    await click(GENERAL.submitButton);
    assert.true(
      this.flashSuccessSpy.calledWith('Successfully saved Entity alias.'),
      'Entities alias: shows a flash message on create'
    );
    const aliasGeneratedId = document.querySelector(GENERAL.infoRowValue('ID')).innerText;
    assert.strictEqual(
      currentURL(),
      `/vault/access/identity/entities/aliases/${aliasGeneratedId}/details`,
      'navigates to the alias show route after creation'
    );
    assert
      .dom(GENERAL.infoRowValue('Name'))
      .hasText(name, `entities renders the alias name on the alias show page`);

    await visit(`/vault/access/identity/entities/aliases`);
    assert
      .dom(`[data-test-identity-link="${aliasGeneratedId}"]`)
      .exists('entities: lists the entity alias in the entity alias list');

    await click(GENERAL.menuItem(name));
    await click('[data-test-popup-menu="delete"]');
    await click(GENERAL.confirmButton);
    assert.dom(GENERAL.latestFlashContent).includesText(`Successfully deleted`);
  });

  test('it allows delete from the edit form', async function (assert) {
    assert.expect(4);
    const name = `alias-${uuidv4()}`;
    await testAliasDeleteFromForm(name, 'entities', assert);
    await settled();
  });
});
