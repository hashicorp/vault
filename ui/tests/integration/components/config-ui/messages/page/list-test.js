/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

// TODO: test file needs to be updated to use mirage handler and to have the correct META pagination numbers

const META = {
  currentPage: 1,
  lastPage: 2,
  nextPage: 2,
  prevPage: 1,
  total: 16,
  filteredTotal: 16,
  pageSize: 15,
};

module('Integration | Component | messages/page/list', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');

    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '01234567-89ab-cdef-0123-456789abcdef',
      active: true,
      type: 'banner',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { title: 'here', href: 'www.example.com' },
      startTime: '2021-08-01T00:00:00Z',
      endTime: '',
    });
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '01234567-89ab-dddd-0123-456789abcdef',
      active: false,
      type: 'modal',
      authenticated: true,
      title: 'Message title 2',
      message: 'Some long long long message blah blah blah',
      link: { title: 'here', href: 'www.example2.com' },
      startTime: '2023-08-01T00:00:00Z',
      endTime: '2023-08-01T00:00:00Z',
    });
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '01234567-89ab-vvvv-0123-456789abcdef',
      active: true,
      type: 'banner',
      authenticated: false,
      title: 'Message title 3',
      message: 'Some long long long message',
      link: { title: 'here', href: 'www.example.com' },
      startTime: '2021-08-01T00:00:00Z',
      endTime: '2090-08-01T00:00:00Z',
    });
  });

  test('it should show the messages empty state', async function (assert) {
    this.messages = [];

    await render(hbs`<Messages::Page::List @messages={{this.messages}} />`, {
      owner: this.engine,
    });

    assert.dom('[data-test-empty-state-title]').hasText('No messages yet');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'Add a custom message for all users after they log into Vault. Create message to get started.'
      );
    assert.dom('[data-test-empty-state-actions] a').hasText('Create message');
  });

  test('it should show the list of custom messages', async function (assert) {
    this.messages = this.store.peekAll('config-ui/message', {});
    this.messages.meta = META;

    await render(hbs`<Messages::Page::List @messages={{this.messages}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-icon="message-circle"]').exists();
    for (const message of this.messages) {
      assert.dom(`[data-test-list-item="${message.id}"]`).exists();
      assert.dom(`[data-linked-block-title="${message.id}"]`).hasText(message.title);
      // TODO: add tests for active/inactive badges
    }
  });
});
