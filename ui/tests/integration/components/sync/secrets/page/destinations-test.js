/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';

import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';

const { title, tab, filter, searchSelect, emptyStateTitle, destinations, menuTrigger, confirmButton } = PAGE;

module('Integration | Component | sync | Page::Destinations', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    this.version.features = ['Secrets Sync'];
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    const store = this.owner.lookup('service:store');
    const modelName = 'sync/destinations/aws-sm';
    const destination = this.server.create('sync-destination', 'aws-sm');
    store.pushPayload(modelName, {
      modelName,
      ...destination,
      id: destination.name,
    });

    // mimic what happens in lazyPaginatedQuery
    this.destinations = store.peekAll(modelName);
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
    assert.dom(title).hasText('Secrets Sync', 'Page title renders');
    assert.dom(tab('Overview')).exists('Overview tab renders');
    assert.dom(tab('Destinations')).exists('Destinations tab renders');
  });

  test('it should render toolbar filters and actions', async function (assert) {
    assert.expect(4);

    this.typeFilter = 'aws-sm';
    await this.renderComponent();

    assert.dom(destinations.list.create).includesText('Create new destination', 'Create action renders');
    assert
      .dom(filter('type'))
      .includesText('AWS Secrets Manager', 'Filter is populated for correct initial value');
    await click(searchSelect.removeSelected);

    // TYPE FILTER
    await click(`${filter('type')} .ember-basic-dropdown-trigger`);
    await click(searchSelect.option(searchSelect.optionIndex('AWS Secrets Manager')));
    assert.deepEqual(
      this.transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations', { queryParams: { type: 'aws-sm' } }],
      'type filter triggered transition with correct query params'
    );

    // NAME FILTER
    await fillIn(filter('name'), 'destination-aws');
    assert.deepEqual(
      this.transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations', { queryParams: { name: 'destination-aws' } }],
      'name filter triggered transition with correct query params'
    );
  });

  test('it should render empty state when there are no filtered results', async function (assert) {
    this.destinations = [];
    this.typeFilter = 'foo';
    this.nameFilter = 'bar';

    await this.renderComponent();

    assert
      .dom(emptyStateTitle)
      .hasText(
        'There are no foo destinations matching "bar".',
        'Renders correct empty state when both type and name filters are defined'
      );

    this.set('nameFilter', undefined);
    assert
      .dom(emptyStateTitle)
      .hasText('There are no foo destinations.', 'Renders correct empty state when type filter is defined');

    this.setProperties({
      typeFilter: undefined,
      nameFilter: 'bar',
    });
    assert.dom(emptyStateTitle).hasText('There are no destinations matching "bar".');
  });

  test('it should render destination list items', async function (assert) {
    assert.expect(6);

    this.server.delete('/sys/sync/destinations/aws-sm/destination-aws', () => {
      assert.ok('Request made to delete destination');
      return {};
    });

    await this.renderComponent();

    const { icon, name, type, deleteAction } = destinations.list;

    assert.dom(icon).hasClass('flight-icon-aws-color', 'Correct icon renders');
    assert.dom(name).hasText('destination-aws', 'Name renders');
    assert.dom(type).hasText('AWS Secrets Manager', 'Type renders');

    await click(menuTrigger);

    await click(deleteAction);
    await click(confirmButton);

    assert.propEqual(
      this.transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.overview'],
      'Transition is triggered on delete success'
    );
    assert.propEqual(
      this.clearDatasetStub.lastCall.args,
      ['sync/destination'],
      'Store dataset is cleared on delete success'
    );
  });
});
