/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

/**
 * The identity edit routes must resolve the policies capability before reading
 * canCreate; reading it off the pending promise always yields undefined, which
 * hides the "create new policy" option in the edit form.
 */
module('Unit | Route | vault/cluster/access/identity edit', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api');
    sinon
      .stub(this.owner.lookup('service:capabilities'), 'for')
      .resolves({ canCreate: true, canDelete: true });
  });

  test('entities edit route returns canCreatePolicy from the resolved capabilities', async function (assert) {
    sinon.stub(this.api.identity, 'entityReadById').resolves({ data: { name: 'my-entity' } });
    const route = this.owner.lookup('route:vault.cluster.access.identity.entities.edit');

    const model = await route.model({ item_id: 'entity-id' });

    assert.true(model.canCreatePolicy, 'canCreatePolicy is true');
  });

  test('groups edit route returns canCreatePolicy from the resolved capabilities', async function (assert) {
    sinon.stub(this.api.identity, 'groupReadById').resolves({ data: { name: 'my-group', type: 'internal' } });
    sinon.stub(this.api.identity, 'groupListById').resolves({ keys: [] });
    sinon.stub(this.api.identity, 'entityListById').resolves({ keys: [] });
    const route = this.owner.lookup('route:vault.cluster.access.identity.groups.edit');

    const model = await route.model({ item_id: 'group-id' });

    assert.true(model.canCreatePolicy, 'canCreatePolicy is true');
  });
});
