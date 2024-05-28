/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { STANDARD_META } from 'vault/tests/helpers/pagination';

module('Integration | Component | pki-paginated-list', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.store.pushPayload('pki/key', {
      modelName: 'pki/key',
      data: {
        key_id: '724862ff-6438-bad0-b598-77a6c7f4e934',
        key_type: 'ec',
        key_name: 'test-key',
      },
    });
    this.store.pushPayload('pki/key', {
      modelName: 'pki/key',
      data: {
        key_id: '9fdddf12-9ce3-0268-6b34-dc1553b00175',
        key_type: 'rsa',
        key_name: 'another-key',
      },
    });
    // mimic what happens in lazyPaginatedQuery
    const keyModels = this.store.peekAll('pki/key');
    keyModels.meta = STANDARD_META;
    this.list = keyModels;
    const emptyList = this.store.peekAll('pki/foo');
    emptyList.meta = {
      meta: {
        total: 0,
        currentPage: 1,
        pageSize: 100,
      },
    };
    this.emptyList = emptyList;
  });

  test('it renders correctly with a list', async function (assert) {
    this.set('hasConfig', null);
    await render(
      hbs`
      <PkiPaginatedList @backend="pki-mount" @list={{this.list}} @hasConfig={{this.hasConfig}}>
        <:list as |items|>
          {{#each items as |item|}}
            <div data-test-item={{item.keyId}}>{{item.keyName}}</div>
          {{/each}}
        </:list>
        <:empty>
          No items found
        </:empty>
        <:configure>
          Not configured
        </:configure>
      </PkiPaginatedList>
    `,
      { owner: this.engine }
    );

    assert.dom(this.element).doesNotContainText('Not configured', 'defaults to has config if not boolean');
    assert.dom(this.element).doesNotContainText('No items found', 'does not render empty state');
    assert.dom('[data-test-item]').exists({ count: 2 }, 'lists the items');
    assert.dom('[data-test-item="724862ff-6438-bad0-b598-77a6c7f4e934"]').hasText('test-key');
    assert.dom('[data-test-item="9fdddf12-9ce3-0268-6b34-dc1553b00175"]').hasText('another-key');
    assert.dom('[data-test-pagination]').exists('shows pagination');
    await this.set('hasConfig', false);
    assert.dom(this.element).doesNotContainText('No items found', 'does not render empty state');
    assert.dom(this.element).containsText('Not configured', 'shows configuration prompt');
    assert.dom('[data-test-item]').doesNotExist('Does not show list items when not configured');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination');
  });

  test('it renders correctly with an empty list', async function (assert) {
    this.set('hasConfig', true);
    await render(
      hbs`
      <PkiPaginatedList @backend="pki-mount" @list={{this.emptyList}} @hasConfig={{this.hasConfig}}>
        <:list>
          List item
        </:list>
        <:empty>
          No items found
        </:empty>
        <:configure>
          Not configured
        </:configure>
      </PkiPaginatedList>
    `,
      { owner: this.engine }
    );

    assert.dom(this.element).doesNotContainText('list item', 'does not render list items if empty');
    assert.dom(this.element).hasText('No items found', 'shows empty block');
    assert.dom(this.element).doesNotContainText('Not configured', 'does not show configuration prompt');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination');
    await this.set('hasConfig', false);
    assert.dom(this.element).doesNotContainText('list item', 'does not render list items if empty');
    assert.dom(this.element).doesNotContainText('No items found', 'does not show empty state');
    assert.dom(this.element).hasText('Not configured', 'shows configuration prompt');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination');
  });

  test('it renders actions, description, pagination', async function (assert) {
    this.set('hasConfig', true);
    this.set('model', this.list);
    await render(
      hbs`
      <PkiPaginatedList @backend="pki-mount" @list={{this.model}} @hasConfig={{this.hasConfig}}>
        <:actions>
          <div data-test-button>Action</div>
        </:actions>
        <:description>
          Description goes here
        </:description>
        <:list>
          List items
        </:list>
        <:empty>
          No items found
        </:empty>
        <:configure>
          Not configured
        </:configure>
      </PkiPaginatedList>
    `,
      { owner: this.engine }
    );
    assert
      .dom('[data-test-button]')
      .includesText('Action', 'Renders actions in toolbar when list and config');
    assert
      .dom(this.element)
      .includesText('Description goes here', 'renders description when list and config');
    assert.dom('[data-test-pagination]').exists('shows pagination when list and config');

    this.set('model', this.emptyList);
    assert
      .dom('[data-test-button]')
      .hasText('Action', 'Renders actions in toolbar when empty list and config');
    assert
      .dom(this.element)
      .doesNotIncludeText('Description goes here', 'hides description when empty list and config');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination when empty list and config');

    this.set('hasConfig', false);
    assert.dom('[data-test-button]').hasText('Action', 'Renders actions in toolbar when no config');
    assert.dom(this.element).doesNotIncludeText('Description goes here', 'hides description when no config');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination when no config');
  });
});
