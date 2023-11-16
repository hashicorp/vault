/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupModels } from 'vault/tests/helpers/sync/setup-models';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, render } from '@ember/test-helpers';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import sinon from 'sinon';

module('Integration | Component | sync | Secrets::DestinationHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);
  setupModels(hooks);

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    await render(hbs`<Secrets::DestinationHeader @destination={{this.destination}} />`, {
      owner: this.engine,
    });
  });

  test('it should SyncHeader component', async function (assert) {
    assert.dom(PAGE.title).includesText('us-west-1', 'SyncHeader component renders');
  });

  test('it should render tabs', async function (assert) {
    assert.dom(PAGE.tab('Secrets')).hasText('Secrets', 'Secrets tab renders');
    assert.dom(PAGE.tab('Details')).hasText('Details', 'Details tab renders');
  });

  test('it should render toolbar', async function (assert) {
    ['Delete destination', 'Sync new secret', 'Edit destination'].forEach((btn) => {
      assert.dom(PAGE.toolbar(btn)).hasText(btn, `${btn} toolbar action renders`);
    });
  });

  test('it should delete destination', async function (assert) {
    assert.expect(3);

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    const clearDatasetStub = sinon.stub(this.store, 'clearDataset');

    this.server.delete('/sys/sync/destinations/aws-sm/us-west-1', () => {
      assert.ok(true, 'Request made to delete destination');
      return {};
    });

    await click(PAGE.toolbar('Delete destination'));
    await fillIn(PAGE.confirmModalInput, 'DELETE');
    await click(PAGE.confirmButton);

    assert.propEqual(
      transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations'],
      'Transition is triggered on delete success'
    );
    assert.propEqual(
      clearDatasetStub.lastCall.args,
      ['sync/destination'],
      'Store dataset is cleared on delete success'
    );
  });
});
