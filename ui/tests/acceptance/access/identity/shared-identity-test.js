/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { v4 as uuidv4 } from 'uuid';
import { click, fillIn, visit, currentURL } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { capitalize } from '@ember/string';
import { singularize } from 'ember-inflector';

// Helper to create an entity or group
async function createEntityOrGroup(itemType, name) {
  await visit(`/vault/access/identity/${itemType}/create`);
  if (itemType === 'groups') {
    await fillIn(GENERAL.inputByAttr('type'), 'external');
  }
  await fillIn(GENERAL.inputByAttr('name'), name);
  await click(GENERAL.submitButton);
  return document.querySelector(GENERAL.infoRowValue('ID')).innerText;
}

// Helper to create an alias
async function createAlias(itemType, itemGeneratedId, name) {
  await visit(`/vault/access/identity/${itemType}/aliases/add/${itemGeneratedId}`);
  await fillIn(GENERAL.inputByAttr('name'), name);
  await click(GENERAL.submitButton);
  return document.querySelector(GENERAL.infoRowValue('ID')).innerText;
}

// This module covers both groups and entities, so the module name differs from the route path.
// Creation of an Entity or Group is inherently tested as part of the alias flow, so no separate test is needed.
module('Acceptance | Create groups and entities alias test', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    await login();
  });

  hooks.afterEach(function () {
    this.flashSuccessSpy.restore();
  });

  for (const itemType of ['groups', 'entities']) {
    test(`${itemType}: it allows create, list, delete of an entity alias`, async function (assert) {
      const name = `${itemType}-${uuidv4()}`;
      const itemGeneratedId = await createEntityOrGroup(itemType, name);

      assert.true(
        this.flashSuccessSpy.calledWith(`Successfully saved ${singularize(capitalize(itemType))}.`),
        `${itemType}: shows a flash message on create`
      );

      const aliasGeneratedId = await createAlias(itemType, itemGeneratedId, name);

      assert.true(
        this.flashSuccessSpy.calledWith(`Successfully saved ${singularize(capitalize(itemType))} alias.`),
        `${itemType}: shows a flash message when creating an alias.`
      );

      assert.strictEqual(
        currentURL(),
        `/vault/access/identity/${itemType}/aliases/${aliasGeneratedId}/details`,
        'navigates to the alias show route after creation'
      );
      assert
        .dom(GENERAL.infoRowValue('Name'))
        .hasText(name, `${itemType}: renders the alias name on the alias show page`);

      await visit(`/vault/access/identity/${itemType}/aliases`);
      assert
        .dom(`[data-test-identity-link="${aliasGeneratedId}"]`)
        .exists(`${itemType}: lists the entity alias`);

      await click(GENERAL.menuItem(name));
      await click('[data-test-popup-menu="delete"]');
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted');
    });

    test(`${itemType}: it allows delete from the edit form`, async function (assert) {
      const name = `${itemType}-${uuidv4()}`;
      const itemGeneratedId = await createEntityOrGroup(itemType, name);
      const aliasGeneratedId = await createAlias(itemType, itemGeneratedId, name);

      await click('[data-test-alias-edit-link]');
      assert.strictEqual(
        currentURL(),
        `/vault/access/identity/${itemType}/aliases/edit/${aliasGeneratedId}`,
        `${itemType}: correctly navigates to edit`
      );

      await click(GENERAL.confirmTrigger); // click the Delete entity-alias trigger button
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted');
      assert.strictEqual(
        currentURL(),
        `/vault/access/identity/${itemType}/aliases`,
        `${itemType}: navigates to the list page after deletion`
      );
    });

    test(`${itemType}: it allows you to delete the ${itemType} from the list view`, async function (assert) {
      const name = `${itemType}-${uuidv4()}`;
      await createEntityOrGroup(itemType, name);
      await visit(`/vault/access/identity/${itemType}`);

      const rowSelector = `[data-test-identity-row="${name}"]`;
      const menuTriggerSelector = `${rowSelector} ${GENERAL.menuTrigger}`;

      assert.dom(rowSelector).exists(`${itemType}: is in the list view`);

      await click(menuTriggerSelector);
      await click(GENERAL.menuItem('delete'));
      await click(GENERAL.confirmButton);

      assert.dom(rowSelector).doesNotExist(`${itemType}: is NOT in the list view`);
    });
  }
});
