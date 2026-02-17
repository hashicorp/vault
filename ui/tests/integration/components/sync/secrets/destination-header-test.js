/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupDataStubs } from 'vault/tests/helpers/sync/setup-hooks';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | sync | Secrets::DestinationHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);
  setupDataStubs(hooks);

  hooks.beforeEach(async function () {
    this.refreshList = sinon.stub();

    this.renderComponent = () =>
      render(
        hbs`<Secrets::DestinationHeader @destination={{this.destination}} @capabilities={{this.capabilities}} @refreshList={{this.refreshList}} />`,
        {
          owner: this.engine,
        }
      );
  });

  test('it should render SyncHeader component', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).includesText('destination-aws', 'SyncHeader component renders');
  });

  test('it should render tabs', async function (assert) {
    await this.renderComponent();
    assert.dom(PAGE.tab('Secrets')).hasText('Secrets', 'Secrets tab renders');
    assert.dom(PAGE.tab('Details')).hasText('Details', 'Details tab renders');
  });

  test('it should render toolbar', async function (assert) {
    await this.renderComponent();
    ['Delete destination', 'Sync secrets', 'Edit destination'].forEach((btn) => {
      assert.dom(PAGE.toolbar(btn)).hasText(btn, `${btn} toolbar action renders`);
    });
  });

  test('it should delete destination', async function (assert) {
    assert.expect(2);

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.server.delete('/sys/sync/destinations/aws-sm/destination-aws', () => {
      assert.ok(true, 'Request made to delete destination');
      return {};
    });

    await this.renderComponent();
    await click(PAGE.toolbar('Delete destination'));
    await fillIn(PAGE.confirmModalInput, 'DELETE');
    await click(PAGE.confirmButton);

    assert.true(
      transitionStub.calledWith('vault.cluster.sync.secrets.overview'),
      'Transition is triggered on delete success'
    );
  });

  test('it should render delete progress banner and hide actions', async function (assert) {
    assert.expect(5);

    this.destination.purge_initiated_at = '2024-01-09T16:54:28.463879';

    await this.renderComponent();
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

    this.destination.purge_initiated_at = '2024-01-09T16:54:28.463879';
    this.destination.purge_error = 'oh no! a problem occurred!';

    await this.renderComponent();
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

    await this.renderComponent();
    await click(PAGE.associations.list.refresh);
    assert.true(this.refreshList.calledOnce, 'Refresh list action is triggered');
  });
});
