/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const emptyModel = () => {
  const model = [];
  Object.defineProperty(model, 'meta', { value: { total: 0 }, writable: false });
  return model;
};

const populatedModel = () => {
  const model = [{ id: 'entity-1', name: 'test-entity' }];
  Object.defineProperty(model, 'meta', { value: { total: 1 }, writable: false });
  return model;
};

const noMetaModel = () => [];

module('Integration | Component | identity/page/entities', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the page header with title, description, and breadcrumbs', async function (assert) {
    this.model = emptyModel();
    await render(hbs`<Identity::Page::Entities @model={{this.model}} />`);

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Entities', 'renders Entities as the page title');
    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .containsText(
        'Create and manage unique identities',
        'renders the entity description in the page header'
      );
    assert.dom(GENERAL.breadcrumbs).exists('renders breadcrumbs');
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists('renders Vault breadcrumb link');
    assert.dom(GENERAL.currentBreadcrumb('Entities')).exists('renders Entities as current breadcrumb');
  });

  test('it renders the Merge entities and Create new entity header actions', async function (assert) {
    this.model = emptyModel();
    await render(hbs`<Identity::Page::Entities @model={{this.model}} />`);

    assert
      .dom(GENERAL.button('entity-merge-link'))
      .hasText('Merge entities', 'renders the Merge entities button with correct label');
    assert
      .dom(GENERAL.button('entity-create-link'))
      .hasText('Create new entity', 'renders the Create new entity button with correct label');
  });

  test('it renders the empty state when there are no entities', async function (assert) {
    this.model = emptyModel();
    await render(hbs`<Identity::Page::Entities @model={{this.model}} />`);

    assert.dom(GENERAL.emptyStateTitle).hasText('No entities yet', 'renders the empty state title');
    assert
      .dom(GENERAL.emptyStateMessage)
      .containsText('Create your first entity to get started', 'renders the empty state message');
    assert.dom(GENERAL.emptyStateActions).exists('renders empty state actions');
  });

  test('it renders the empty state when the model has no meta (404 case)', async function (assert) {
    this.model = noMetaModel();
    await render(hbs`<Identity::Page::Entities @model={{this.model}} />`);

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No entities yet', 'renders empty state when model has no meta property');
  });

  test('it does not render the empty state when entities exist', async function (assert) {
    this.model = populatedModel();
    await render(hbs`<Identity::Page::Entities @model={{this.model}} />`);

    assert
      .dom(GENERAL.emptyStateTitle)
      .doesNotExist('does not render the empty state when entities are present');
  });
});
