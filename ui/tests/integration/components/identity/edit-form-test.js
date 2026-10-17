/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import GroupIdentityForm from 'vault/forms/identity/group';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

/**
 * Integration tests for the identity/edit-form component focusing on the
 * SearchSelect fields for member_entity_ids and member_group_ids.
 * Previously, when editing a group that already had member entity UUIDs set,
 * the UUID appeared twice (as both name and id) because the edit route did
 * not pass @model.entities / @model.groups. These tests verify correct name
 * resolution when the option lists are provided.
 */
module('Integration | Component | identity/edit-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    // Stub policy fetches called in the component constructor so they don't fail.
    const api = this.owner.lookup('service:api');
    sinon.stub(api.sys, 'policiesListAclPolicies').resolves({ keys: [] });
    sinon.stub(api.sys, 'systemListPoliciesRgp').resolves({ keys: [] });
  });

  function renderGroupEditForm(ctx, { memberEntityIds = [], entities = [], groups = [] } = {}) {
    const form = new GroupIdentityForm(
      { name: 'my-group', type: 'internal', member_entity_ids: memberEntityIds, member_group_ids: [] },
      { isNew: false }
    );
    ctx.model = {
      form,
      entities,
      groups,
      identityType: 'group',
      itemId: 'test-group-uuid',
      canCreatePolicies: false,
    };
    ctx.onSave = sinon.stub();
    return render(hbs`
      <Identity::EditForm
        @model={{this.model}}
        @mode="edit"
        @onSave={{this.onSave}}
      />
    `);
  }

  test('shows entity name (not UUID twice) in selected list when @model.entities is provided', async function (assert) {
    /**
     * Regression: when editing a group on the edit route, @model.entities was
     * absent, so the SearchSelect fell back to using the UUID as both name and
     * id. With the entities list populated, the name should resolve correctly.
     */
    const entityId = '0cd43cbe-c44d-039e-9412-f4a34280d30d';
    const entities = [{ id: entityId, name: 'my-entity-name' }];

    await renderGroupEditForm(this, { memberEntityIds: [entityId], entities });

    // The selected list item should display the entity name.
    assert
      .dom(GENERAL.searchSelect.selectedOption(0))
      .includesText('my-entity-name', 'selected item shows entity name, not the UUID');

    // The UUID should appear in the secondary (small) id element, not as the primary label.
    assert.dom(`[data-test-smaller-id="0"]`).hasText(entityId, 'UUID is shown in the secondary id element');
  });

  test('shows UUID in selected list when entity is not present in @model.entities', async function (assert) {
    /**
     * When the entity UUID is not found in the options list (e.g. it was deleted),
     * the component should gracefully fall back to displaying the UUID rather than erroring.
     */
    const unknownId = 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee';

    await renderGroupEditForm(this, { memberEntityIds: [unknownId], entities: [] });

    assert
      .dom(GENERAL.searchSelect.selectedOption(0))
      .includesText(unknownId, 'falls back to showing UUID when entity is not in options list');
  });
});
