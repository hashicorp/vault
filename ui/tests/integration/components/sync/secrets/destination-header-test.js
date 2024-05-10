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
import { click, fillIn, render, settled } from '@ember/test-helpers';
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

  test('it should render SyncHeader component', async function (assert) {
    assert.dom(PAGE.title).includesText('us-west-1', 'SyncHeader component renders');
  });

  test('it should render tabs', async function (assert) {
    assert.dom(PAGE.tab('Secrets')).hasText('Secrets', 'Secrets tab renders');
    assert.dom(PAGE.tab('Details')).hasText('Details', 'Details tab renders');
  });

  test('it should render toolbar', async function (assert) {
    ['Delete destination', 'Sync secrets', 'Edit destination'].forEach((btn) => {
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
      ['vault.cluster.sync.secrets.overview'],
      'Transition is triggered on delete success'
    );
    assert.propEqual(
      clearDatasetStub.lastCall.args,
      ['sync/destination'],
      'Store dataset is cleared on delete success'
    );
  });

  test('it should render delete progress banner and hide actions', async function (assert) {
    assert.expect(5);
    this.destination.set('purgeInitiatedAt', '2024-01-09T16:54:28.463879');
    await settled();
    assert
      .dom(PAGE.destinations.deleteBanner)
      .hasText(
        'Deletion in progress Purge initiated on Jan 09, 2024 at 04:54:28 pm. This process may take some time depending on how many secrets must be un-synced from this destination.'
      );
    assert
      .dom(`${PAGE.destinations.deleteBanner} ${PAGE.icon('loading-static')}`)
      .exists('banner renders loading icon');
    assert.dom(PAGE.toolbar('Sync secrets')).doesNotExist('Sync action is hidden');
    assert.dom(PAGE.toolbar('Edit destination')).doesNotExist('Edit action is hidden');
    assert.dom('.toolbar-separator').doesNotExist('Divider is hidden when only delete action is available');
  });

  test('it should render delete error banner', async function (assert) {
    assert.expect(2);
    this.destination.set('purgeInitiatedAt', '2024-01-09T16:54:28.463879');
    this.destination.set('purgeError', 'oh no! a problem occurred!');
    await settled();
    assert
      .dom(PAGE.destinations.deleteBanner)
      .hasText(
        'Deletion failed There was a problem with the delete purge initiated at Jan 09, 2024 at 04:54:28 pm. oh no! a problem occurred!',
        'banner renders error message'
      );
    assert
      .dom(`${PAGE.destinations.deleteBanner} ${PAGE.icon('alert-diamond')}`)
      .exists('banner renders critical icon');
  });

  test('it should render refresh list button', async function (assert) {
    assert.expect(1);

    this.refreshList = () => assert.ok(true, 'Refresh list callback fires');

    await render(
      hbs`<Secrets::DestinationHeader @destination={{this.destination}} @refreshList={{this.refreshList}} />`,
      {
        owner: this.engine,
      }
    );

    await click(PAGE.associations.list.refresh);
  });
});
