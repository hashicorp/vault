/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';

module('Integration | Component | sync | Page::Destinations', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:version').version = '1.16.0+ent';

    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    const store = this.owner.lookup('service:store');
    const modelName = 'sync/destinations/aws-sm';
    const destination = this.server.create('sync-destination', 'aws-sm');
    store.pushPayload(modelName, {
      modelName,
      ...destination,
      id: destination.name,
    });

    this.destinations = store.peekAll(modelName).toArray();
    this.destinations.meta = {
      filteredTotal: this.destinations.length,
      currentPage: 1,
      pageSize: 5,
    };

    this.renderComponent = () => {
      return render(
        hbs`<Secrets::Page::Destinations
          @destinations={{this.destinations}}
          @typeFilter={{this.typeFilter}}
          @nameFilter={{this.nameFilter}}
        />`,
        { owner: this.engine }
      );
    };

    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.clearDatasetStub = sinon.stub(store, 'clearDataset');
  });

  test('it should render header and tabs', async function (assert) {
    await this.renderComponent();
    assert.dom('[data-test-crumb="0"]').includesText('Secrets sync', 'Breadcrumb renders');
    assert.dom('[data-test-page-title]').hasText('Secrets sync', 'Page title renders');
    assert.dom('[data-test-tab="Overview"]').exists('Overview tab renders');
    assert.dom('[data-test-tab="Destinations"]').exists('Destinations tab renders');
  });

  test('it should render toolbar filters and actions', async function (assert) {
    assert.expect(4);

    this.typeFilter = 'aws-sm';
    await this.renderComponent();

    assert
      .dom('[data-test-create-destination]')
      .includesText('Create new destination', 'Create action renders');
    assert
      .dom('[data-test-filter="type"]')
      .includesText('AWS Secrets Manager', 'Filter is populated for correct initial value');
    await click('[data-test-selected-list-button="delete"]');

    for (const filter of ['type', 'name']) {
      await click(`[data-test-filter="${filter}"] .ember-basic-dropdown-trigger`);
      await click('[data-option-index="0"]');

      const value = filter === 'type' ? 'aws-sm' : 'destination-aws';
      assert.deepEqual(
        this.transitionStub.lastCall.args,
        ['vault.cluster.sync.secrets.destinations', { queryParams: { [filter]: value } }],
        `${filter} filter triggered transition with correct query params`
      );
    }
  });

  test('it should render empty state when there are no filtered results', async function (assert) {
    this.destinations = [];
    this.typeFilter = 'foo';
    this.nameFilter = 'bar';

    await this.renderComponent();

    const selector = '[data-test-empty-state-title]';
    assert
      .dom(selector)
      .hasText(
        'There are no foo destinations matching "bar".',
        'Renders correct empty state when both type and name filters are defined'
      );

    this.set('nameFilter', undefined);
    assert
      .dom(selector)
      .hasText('There are no foo destinations.', 'Renders correct empty state when type filter is defined');

    this.setProperties({
      typeFilter: undefined,
      nameFilter: 'bar',
    });
    assert.dom(selector).hasText('There are no destinations matching "bar".');
  });

  test('it should render destination list items', async function (assert) {
    assert.expect(6);

    this.server.delete('/sys/sync/destinations/aws-sm/destination-aws', () => {
      assert.ok('Request made to delete destination');
      return {};
    });

    await this.renderComponent();

    assert.dom('[data-test-destination-icon]').hasClass('flight-icon-aws-color', 'Correct icon renders');
    assert.dom('[data-test-destination-name]').hasText('destination-aws', 'Name renders');
    assert.dom('[data-test-destination-type]').hasText('AWS Secrets Manager', 'Type renders');

    await click('[data-test-popup-menu-trigger]');

    await click('[data-test-delete]');
    await click('[data-test-confirm-button]');

    assert.propEqual(
      this.transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations'],
      'Transition is triggered on delete success'
    );
    assert.propEqual(
      this.clearDatasetStub.lastCall.args,
      ['sync/destinations/aws-sm'],
      'Store dataset is cleared on delete success'
    );
  });
});
