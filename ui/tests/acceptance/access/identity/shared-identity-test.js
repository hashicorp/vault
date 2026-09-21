/**
 * Copyright IBM Corp. 2016, 2025
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
import { setupMirage } from 'ember-cli-mirage/test-support';
import { runCmd, mountAuthCmd } from 'vault/tests/helpers/commands';

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
  setupMirage(hooks);

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
        this.flashSuccessSpy.calledWith(`Successfully saved ${singularize(capitalize(itemType))}: ${name}.`),
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

    test(`${itemType}: it shows alias management options on the aliases tab of the item details page`, async function (assert) {
      const name = `${itemType}-${uuidv4()}`;
      const itemGeneratedId = await createEntityOrGroup(itemType, name);
      const aliasName = `alias-${uuidv4()}`;
      const aliasGeneratedId = await createAlias(itemType, itemGeneratedId, aliasName);

      await visit(`/vault/access/identity/${itemType}/${itemGeneratedId}/aliases`);
      await click(`[data-test-popup-menu="${aliasName}"]`);

      assert
        .dom('[data-test-alias-actions-menu]')
        .hasText('Details Edit Remove', `${itemType}: alias actions menu shows all options`);

      await click(`a[href*="/aliases/edit/"]`);

      assert.strictEqual(
        currentURL(),
        `/vault/access/identity/${itemType}/aliases/edit/${aliasGeneratedId}`,
        `${itemType}: navigates to the alias edit page`
      );
    });

    test(`${itemType}: it allows editing an alias's auth backend`, async function (assert) {
      const name = `${itemType}-${uuidv4()}`;
      const itemGeneratedId = await createEntityOrGroup(itemType, name);
      const aliasName = `alias-${uuidv4()}`;
      const aliasGeneratedId = await createAlias(itemType, itemGeneratedId, aliasName);

      const authPath = `userpass-${uuidv4().slice(0, 8)}`;
      await runCmd(mountAuthCmd('userpass', authPath));

      await visit(`/vault/access/identity/${itemType}/aliases/edit/${aliasGeneratedId}`);

      const select = document.querySelector('[data-test-mount-accessor-select]');
      const option = [...select.options].find((opt) => opt.text.includes(authPath));
      await fillIn('[data-test-mount-accessor-select]', option.value);
      await click(GENERAL.submitButton);

      assert.dom(GENERAL.messageError).doesNotExist(`${itemType}: no raw API error is shown`);
      assert
        .dom(GENERAL.latestFlashContent)
        .includesText('Successfully saved', `${itemType}: shows a flash message on save`);
      assert
        .dom(GENERAL.infoRowValue('Mount'))
        .includesText(authPath, `${itemType}: alias is saved with the newly selected auth backend`);
    });

    test(`${itemType}: it allows delete from the edit form`, async function (assert) {
      assert.expect(3);
      const itemId = uuidv4();
      const name = `${itemType}-${itemId}`;
      this.server.get(`/identity/${singularize(itemType)}/id/${itemId}`, () => {
        return { data: { id: itemId, name } };
      });
      this.server.delete(`/identity/${singularize(itemType)}/id/${itemId}`, () => {
        assert.true(true, `request made to delete ${name}`);
      });
      await visit(`/vault/access/identity/${itemType}/edit/${itemId}`);
      await click(GENERAL.confirmTrigger); // click the Delete entity-alias trigger button
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted');
      assert.strictEqual(
        currentURL(),
        `/vault/access/identity/${itemType}`,
        `${itemType}: navigates to the list page after deletion`
      );
    });

    test(`${itemType}: it allows you to delete the ${itemType} from the list view`, async function (assert) {
      assert.expect(3);
      const itemId = uuidv4();
      const name = `${itemType}-${itemId}`;

      this.server.get(`/identity/${singularize(itemType)}/id`, () => {
        return {
          data: {
            key_info: { [itemId]: { name } },
            keys: [itemId],
          },
        };
      });

      this.server.get(`/identity/${singularize(itemType)}/id/${itemId}`, () => {
        return { data: { id: itemId, name } };
      });

      this.server.delete(`/identity/${singularize(itemType)}/id/${itemId}`, () => {
        assert.true(true, `request made to delete ${name}`);
      });

      await visit(`/vault/access/identity/${itemType}`);

      const rowSelector =
        itemType === 'groups' ? `[data-test-identity-link="${name}"]` : GENERAL.listItem(name);
      const menuTriggerSelector =
        itemType === 'groups'
          ? `[data-test-popup-menu-trigger="${name}"]`
          : `${rowSelector} ${GENERAL.menuTrigger}`;
      const deleteItemSelector =
        itemType === 'groups' ? GENERAL.menuItem('delete-group') : GENERAL.menuItem('delete');

      assert.dom(rowSelector).exists(`${itemType}: is in the list view`);

      await click(menuTriggerSelector);
      await click(deleteItemSelector);
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted');
    });
  }
});
