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
import { Response } from 'miragejs';

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

      const rowSelector = `[data-test-identity-row="${name}"]`;
      const menuTriggerSelector = `${rowSelector} ${GENERAL.menuTrigger}`;

      assert.dom(rowSelector).exists(`${itemType}: is in the list view`);

      await click(menuTriggerSelector);
      await click(GENERAL.menuItem('delete'));
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted');
    });
  }

  test('entity groups: it displays group names and keeps IDs in tooltips', async function (assert) {
    const entityId = 'entity-id';
    const directGroupId = 'direct-group-id';
    const inheritedGroupId = 'inherited-group-id';
    let groupListRequests = 0;

    this.server.get(`/identity/entity/id/${entityId}`, () => ({
      data: {
        id: entityId,
        name: 'example entity',
        direct_group_ids: [directGroupId],
        inherited_group_ids: [inheritedGroupId],
      },
    }));
    this.server.get('/identity/group/id', () => {
      groupListRequests++;
      return {
        data: {
          keys: [directGroupId, inheritedGroupId],
          key_info: {
            [directGroupId]: { name: 'direct group' },
            [inheritedGroupId]: { name: 'inherited group' },
          },
        },
      };
    });

    await visit(`/vault/access/identity/entities/${entityId}/groups`);

    assert
      .dom(`[data-test-identity-item-name="${directGroupId}"]`)
      .hasText('direct group')
      .hasAttribute('title', directGroupId);
    assert
      .dom(`[data-test-identity-item-name="${inheritedGroupId}"]`)
      .hasText('inherited group')
      .hasAttribute('title', inheritedGroupId);
    assert.strictEqual(groupListRequests, 1, 'loads all group names with one list request');
  });

  test('group members: it displays names for member groups and entities', async function (assert) {
    const groupId = 'group-id';
    const memberGroupId = 'member-group-id';
    const memberEntityId = 'member-entity-id';
    let groupListRequests = 0;
    let entityListRequests = 0;

    this.server.get(`/identity/group/id/${groupId}`, () => ({
      data: {
        id: groupId,
        name: 'example group',
        type: 'internal',
        member_group_ids: [memberGroupId],
        member_entity_ids: [memberEntityId],
      },
    }));
    this.server.get('/identity/group/id', () => {
      groupListRequests++;
      return {
        data: {
          keys: [memberGroupId],
          key_info: { [memberGroupId]: { name: 'member group' } },
        },
      };
    });
    this.server.get('/identity/entity/id', () => {
      entityListRequests++;
      return {
        data: {
          keys: [memberEntityId],
          key_info: { [memberEntityId]: { name: 'member entity' } },
        },
      };
    });

    await visit(`/vault/access/identity/groups/${groupId}/members`);

    assert
      .dom(`[data-test-identity-item-name="${memberGroupId}"]`)
      .hasText('member group')
      .hasAttribute('title', memberGroupId);
    assert
      .dom(`[data-test-identity-item-name="${memberEntityId}"]`)
      .hasText('member entity')
      .hasAttribute('title', memberEntityId);
    assert.strictEqual(groupListRequests, 1, 'loads all member group names with one list request');
    assert.strictEqual(entityListRequests, 1, 'loads all member entity names with one list request');
  });

  test('related identity names: it falls back to IDs when the list request fails', async function (assert) {
    const entityId = 'entity-id';
    const groupId = 'group-id';

    this.server.get(`/identity/entity/id/${entityId}`, () => ({
      data: {
        id: entityId,
        name: 'example entity',
        direct_group_ids: [groupId],
      },
    }));
    this.server.get('/identity/group/id', () => new Response(403, {}, { errors: ['permission denied'] }));

    await visit(`/vault/access/identity/entities/${entityId}/groups`);

    assert.dom(`[data-test-identity-item-name="${groupId}"]`).hasText(groupId);
  });
});
