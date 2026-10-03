/**
 * Copyright IBM Corp. 2016, 2025
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
    this.secretMountPath = this.owner.lookup('service:secret-mount-path').update('pki-test');
    this.list = [
      {
        key_id: '724862ff-6438-bad0-b598-77a6c7f4e934',
        key_type: 'ec',
        key_name: 'test-key',
      },
      {
        key_id: '9fdddf12-9ce3-0268-6b34-dc1553b00175',
        key_type: 'rsa',
        key_name: 'another-key',
      },
    ];
    // mimic what happens in paginate util
    this.list.meta = STANDARD_META;
    this.emptyList = [];
    this.emptyList.meta = {
      meta: {
        total: 0,
        currentPage: 1,
        pageSize: 100,
      },
    };
    this.hasConfig = true;

    this.renderComponent = () =>
      render(
        hbs`
          <PkiPaginatedList @backend="pki-mount" @list={{this.list}} @hasConfig={{this.hasConfig}}>
            <:actions>
              <div data-test-button>Action</div>
            </:actions>
            <:description>
              <span data-test-description>Description goes here</span>
            </:description>
            <:list as |items|>
              {{#each items as |item|}}
                <div data-test-item={{item.key_id}}>{{item.key_name}}</div>
              {{/each}}
            </:list>
            <:empty>
              <span data-test-empty>No items found</span>
            </:empty>
            <:configure>
              <span data-test-no-config>Not configured</span>
            </:configure>
        </PkiPaginatedList>
      `,
        { owner: this.engine }
      );
  });

  test('it renders correctly with a list', async function (assert) {
    this.hasConfig = null;
    await this.renderComponent();

    assert.dom('[data-test-no-config]').doesNotExist('defaults to has config if not boolean');
    assert.dom('[data-test-empty]').doesNotExist('No items found', 'does not render empty state');
    assert.dom('[data-test-item]').exists({ count: 2 }, 'lists the items');
    assert.dom('[data-test-item="724862ff-6438-bad0-b598-77a6c7f4e934"]').hasText('test-key');
    assert.dom('[data-test-item="9fdddf12-9ce3-0268-6b34-dc1553b00175"]').hasText('another-key');
    assert.dom('[data-test-pagination]').exists('shows pagination');

    this.hasConfig = false;
    await this.renderComponent();
    assert.dom('[data-test-empty]').doesNotExist('No items found', 'does not render empty state');
    assert.dom('[data-test-no-config]').hasText('Not configured', 'shows configuration prompt');
    assert.dom('[data-test-item]').doesNotExist('Does not show list items when not configured');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination');
  });

  test('it renders correctly with an empty list', async function (assert) {
    this.list = this.emptyList;
    await this.renderComponent();
    assert.dom('[data-test-item]').doesNotExist('does not render list items if empty');
    assert.dom('[data-test-empty]').hasText('No items found', 'shows empty block');
    assert.dom('[data-test-no-config]').doesNotExist('does not show configuration prompt');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination');

    this.hasConfig = false;
    await this.renderComponent();
    assert.dom('[data-test-item]').doesNotExist('does not render list items if empty');
    assert.dom('[data-test-empty]').doesNotExist('does not show empty state');
    assert.dom('[data-test-no-config]').hasText('Not configured', 'shows configuration prompt');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination');
  });

  test('it renders actions, description, pagination', async function (assert) {
    await this.renderComponent();
    assert
      .dom('[data-test-button]')
      .includesText('Action', 'Renders actions in toolbar when list and config');
    assert
      .dom('[data-test-description]')
      .hasText('Description goes here', 'renders description when list and config');
    assert.dom('[data-test-pagination]').exists('shows pagination when list and config');

    this.list = this.emptyList;
    await this.renderComponent();
    assert
      .dom('[data-test-button]')
      .hasText('Action', 'Renders actions in toolbar when empty list and config');
    assert.dom('[data-test-description]').doesNotExist('hides description when empty list and config');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination when empty list and config');

    this.hasConfig = false;
    await this.renderComponent();
    assert.dom('[data-test-button]').hasText('Action', 'Renders actions in toolbar when no config');
    assert.dom('[data-test-description]').doesNotExist('hides description when no config');
    assert.dom('[data-test-pagination]').doesNotExist('hides pagination when no config');
  });
});
