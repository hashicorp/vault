/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

const SELECTORS = {
  trigger: '[data-test-popup-menu-trigger]',
  remove: '.hds-dropdown-list-item button',
  confirm: '[data-test-confirm-button]',
};

// Exercises the three "remove" menus on the identity show tabs. Each must send only the changed
// field to the update-by-id endpoint (Vault treats the request as a partial update), report the
// outcome and refresh the route so the list re-renders.
module('Integration | Component | identity/popups', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api');
    this.router = this.owner.lookup('service:router');
    this.flash = this.owner.lookup('service:flash-messages');
    this.flash.registerTypes(['success', 'danger']);
    this.successSpy = sinon.spy(this.flash, 'success');
    this.dangerSpy = sinon.spy(this.flash, 'danger');
    this.refreshStub = sinon.stub(this.router, 'refresh');
    this.entityUpdate = sinon.stub(this.api.identity, 'entityUpdateById').resolves();
    this.groupUpdate = sinon.stub(this.api.identity, 'groupUpdateById').resolves();

    this.confirmRemoval = async () => {
      await click(SELECTORS.trigger);
      await click(SELECTORS.remove);
      await click(SELECTORS.confirm);
    };
  });

  module('PopupPolicy', function () {
    test('it removes a policy from an entity', async function (assert) {
      this.model = { id: 'e-1', identityType: 'entity', policies: ['default', 'doomed'] };
      await render(hbs`<Identity::PopupPolicy @model={{this.model}} @policyName="doomed" />`);
      await this.confirmRemoval();

      assert.true(this.entityUpdate.calledOnceWithExactly('e-1', { policies: ['default'] }));
      assert.true(this.groupUpdate.notCalled);
      assert.true(
        this.refreshStub.calledOnceWithExactly('vault.cluster.access.identity.entities.show'),
        'only the show route is refreshed so the success flash is not torn down'
      );
      assert.true(this.successSpy.calledOnceWith("Successfully removed 'doomed' policy from e-1"));
    });

    test('it removes a policy from a group and can leave the list empty', async function (assert) {
      this.model = { id: 'g-1', identityType: 'group', policies: ['doomed'] };
      await render(hbs`<Identity::PopupPolicy @model={{this.model}} @policyName="doomed" />`);
      await this.confirmRemoval();

      assert.true(this.groupUpdate.calledOnceWithExactly('g-1', { policies: [] }));
      assert.true(this.entityUpdate.notCalled);
    });

    test('it surfaces an error and does not refresh when the request fails', async function (assert) {
      this.entityUpdate.rejects(new Error('permission denied'));
      this.model = { id: 'e-1', identityType: 'entity', policies: ['default', 'doomed'] };
      await render(hbs`<Identity::PopupPolicy @model={{this.model}} @policyName="doomed" />`);
      await this.confirmRemoval();

      assert.true(this.refreshStub.notCalled);
      assert.true(this.successSpy.notCalled);
      assert.true(this.dangerSpy.calledOnce);
      assert.true(this.dangerSpy.firstCall.args[0].includes("problem removing 'doomed' policy"));
      assert.deepEqual(this.model.policies, ['default', 'doomed'], 'model is not mutated');
    });
  });

  module('PopupMetadata', function () {
    test('it removes a metadata key without mutating the model', async function (assert) {
      this.model = { id: 'e-1', identityType: 'entity', metadata: { keep: 'a', drop: 'b' } };
      await render(hbs`<Identity::PopupMetadata @model={{this.model}} @key="drop" />`);
      await this.confirmRemoval();

      assert.true(this.entityUpdate.calledOnceWithExactly('e-1', { metadata: { keep: 'a' } }));
      assert.deepEqual(this.model.metadata, { keep: 'a', drop: 'b' }, 'model is not mutated');
      assert.true(this.refreshStub.calledOnceWithExactly('vault.cluster.access.identity.entities.show'));
      assert.true(this.successSpy.calledOnceWith("Successfully removed 'drop' from metadata"));
    });

    test('it uses the group endpoint for groups', async function (assert) {
      this.model = { id: 'g-1', identityType: 'group', metadata: { drop: 'b' } };
      await render(hbs`<Identity::PopupMetadata @model={{this.model}} @key="drop" />`);
      await this.confirmRemoval();

      assert.true(this.groupUpdate.calledOnceWithExactly('g-1', { metadata: {} }));
    });

    test('it surfaces an error and does not refresh when the request fails', async function (assert) {
      this.groupUpdate.rejects(new Error('boom'));
      this.model = { id: 'g-1', identityType: 'group', metadata: { drop: 'b' } };
      await render(hbs`<Identity::PopupMetadata @model={{this.model}} @key="drop" />`);
      await this.confirmRemoval();

      assert.true(this.refreshStub.notCalled);
      assert.true(this.dangerSpy.calledOnce);
      assert.deepEqual(this.model.metadata, { drop: 'b' });
    });
  });

  module('PopupMembers', function () {
    test('it removes an entity member from a group', async function (assert) {
      this.model = { id: 'g-1', identityType: 'group', member_entity_ids: ['m-1', 'm-2'] };
      await render(
        hbs`<Identity::PopupMembers @model={{this.model}} @groupArray="member_entity_ids" @memberId="m-2" />`
      );
      await this.confirmRemoval();

      assert.true(this.groupUpdate.calledOnceWithExactly('g-1', { member_entity_ids: ['m-1'] }));
      assert.true(this.refreshStub.calledOnceWithExactly('vault.cluster.access.identity.groups.show'));
      assert.true(this.successSpy.calledOnceWith("Successfully removed 'm-2' from the group"));
    });

    test('it removes a member group from a group', async function (assert) {
      this.model = { id: 'g-1', identityType: 'group', member_group_ids: ['sub-1'] };
      await render(
        hbs`<Identity::PopupMembers @model={{this.model}} @groupArray="member_group_ids" @memberId="sub-1" />`
      );
      await this.confirmRemoval();

      assert.true(this.groupUpdate.calledOnceWithExactly('g-1', { member_group_ids: [] }));
    });

    test('it surfaces an error and does not refresh when the request fails', async function (assert) {
      this.groupUpdate.rejects(new Error('boom'));
      this.model = { id: 'g-1', identityType: 'group', member_entity_ids: ['m-1'] };
      await render(
        hbs`<Identity::PopupMembers @model={{this.model}} @groupArray="member_entity_ids" @memberId="m-1" />`
      );
      await this.confirmRemoval();

      assert.true(this.refreshStub.notCalled);
      assert.true(this.dangerSpy.calledOnce);
      assert.deepEqual(this.model.member_entity_ids, ['m-1']);
    });
  });
});
