import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const recordsData = Array.apply(null, Array(100)).map(function (x, i) {
  return {
    id: `paginated-${i}`,
    type: 'secret',
    attributes: {
      backend: 'test',
    },
  };
});

const SELECTORS = {
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  pagination: '[data-test-client-pagination-control]',
  pageSizeSelector: '.hds-pagination-size-selector .hds-form-select',
};

module('Integration | Component | client-pagination', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.store = await this.owner.lookup('service:store');
    this.store.push({
      data: recordsData,
    });
    this.recordsList = this.store.peekAll('secret');
    this.emptyList = this.store.peekAll('node');
  });

  test('it renders empty state with default messages', async function (assert) {
    this.set('noun', '');
    await render(hbs`<ClientPagination @items={{this.emptyList}} @itemNoun={{this.noun}} />`);

    assert.dom(SELECTORS.emptyStateTitle).hasText('No items yet');
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText('Your items will be listed here. Add your first item to get started.');

    this.set('noun', 'node');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No nodes yet');
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText('Your nodes will be listed here. Add your first node to get started.');
    assert.dom(SELECTORS.pagination).doesNotExist('Pagination not shown');
  });

  test('it renders emptyActions within empty state', async function (assert) {
    await render(hbs`
      <ClientPagination @items={{this.emptyList}}>
        <:emptyActions>
          <div data-test-my-action>
            Action rendered here
          </div>
        </:emptyActions>
      </ClientPagination>
    `);

    assert.dom(`${SELECTORS.emptyStateActions} [data-test-my-action]`).hasText('Action rendered here');
    assert.dom(SELECTORS.pagination).doesNotExist('Pagination not shown');
  });

  test('it does not render pagination if record count <= min page size', async function (assert) {
    this.set('items', this.recordsList.slice(0, 10));
    await render(hbs`
      <ClientPagination @items={{this.items}}>
        <:item as |item|>
          <div data-test-item={{item.id}}>
            {{item.id}}
          </div>
        </:item>
      </ClientPagination>
  `);

    assert.dom(SELECTORS.emptyStateTitle).doesNotExist('No empty state');
    assert.dom(SELECTORS.pagination).doesNotExist('Pagination is not rendered');
    assert.dom('[data-test-item]').exists({ count: 10 }, `10 items are rendered`);
  });

  test('it renders the correct number of items on the page', async function (assert) {
    this.set('items', this.recordsList);
    await render(hbs`
      <ClientPagination @items={{this.items}} >
        <:item as |item|>
          <div data-test-item={{item.id}}>
            {{item.id}}
          </div>
        </:item>
      </ClientPagination>
    `);
    assert.dom(SELECTORS.emptyStateTitle).doesNotExist('No empty state');
    assert.dom(SELECTORS.pagination).exists('Pagination is rendered');
    assert.dom('[data-test-item]').exists({ count: 10 }, `10 items are rendered`);
    await fillIn(SELECTORS.pageSizeSelector, 30);
    assert.dom('[data-test-item]').exists({ count: 30 }, `30 items are rendered`);
    await fillIn(SELECTORS.pageSizeSelector, 50);
    assert.dom('[data-test-item]').exists({ count: 50 }, `50 items are rendered`);
  });
});
