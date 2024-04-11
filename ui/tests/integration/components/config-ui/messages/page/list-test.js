/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

const META = {
  currentPage: 1,
  lastPage: 1,
  nextPage: 1,
  prevPage: 1,
  total: 3,
  pageSize: 15,
};

module('Integration | Component | messages/page/list', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');

    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '0',
      active: true,
      type: 'banner',
      authenticated: true,
      title: 'Message title 1',
      message: 'Some long long long message',
      link: { title: 'here', href: 'www.example.com' },
      start_time: '2021-08-01T00:00:00Z',
      end_time: '',
    });
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '1',
      active: false,
      type: 'modal',
      authenticated: true,
      title: 'Message title 2',
      message: 'Some long long long message blah blah blah',
      link: { title: 'here', href: 'www.example2.com' },
      start_time: '2023-07-01T00:00:00Z',
      end_time: '2023-08-01T00:00:00Z',
    });
    this.store.pushPayload('config-ui/message', {
      modelName: 'config-ui/message',
      id: '2',
      active: false,
      type: 'banner',
      authenticated: false,
      title: 'Message title 3',
      message: 'Some long long long message',
      link: { title: 'here', href: 'www.example.com' },
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
  });

  test('it should show the list of custom messages', async function (assert) {
    this.messages = this.store.peekAll('config-ui/message', {});
    this.messages.meta = META;
    await render(hbs`<Messages::Page::List @messages={{this.messages}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-icon="message-circle"]').exists();
    for (const message of this.messages) {
      assert.dom(CUSTOM_MESSAGES.listItem('Message title 1')).exists();
      assert.dom(`[data-linked-block-title="${message.id}"]`).hasText(message.title);
    }
  });

  test('it should show max message warning modal', async function (assert) {
    for (let i = 0; i < 97; i++) {
      this.store.pushPayload('config-ui/message', {
        modelName: 'config-ui/message',
        id: `${i}-a`,
        active: true,
        type: 'banner',
        authenticated: false,
        title: `Message title ${i}`,
        message: 'Some long long long message',
        link: { title: 'here', href: 'www.example.com' },
        start_time: '2021-08-01T00:00:00Z',
      });
    }

    this.messages = this.store.peekAll('config-ui/message', {});
    this.messages.meta = {
      currentPage: 1,
      lastPage: 1,
      nextPage: 1,
      prevPage: 1,
      total: this.messages.length,
      pageSize: 100,
    };
    await render(hbs`<Messages::Page::List @messages={{this.messages}} />`, {
      owner: this.engine,
    });
    await click(CUSTOM_MESSAGES.button('create message'));
    assert
      .dom(CUSTOM_MESSAGES.modalTitle('maximum-message-modal'))
      .hasText('Maximum number of messages reached');
    assert
      .dom(CUSTOM_MESSAGES.modalBody('maximum-message-modal'))
      .hasText(
        'Vault can only store up to 100 messages. To create a message, delete one of your messages to clear up space.'
      );
    await click(CUSTOM_MESSAGES.modalButton('maximum-message-modal'));
  });

  test('it should show the correct badge colors based on badge status', async function (assert) {
    this.messages = this.store.peekAll('config-ui/message', {});
    this.messages.meta = META;
    await render(hbs`<Messages::Page::List @messages={{this.messages}} />`, {
      owner: this.engine,
    });
    assert.dom(CUSTOM_MESSAGES.badge('0')).hasClass('hds-badge--color-success');
    assert.dom(CUSTOM_MESSAGES.badge('1')).hasClass('hds-badge--color-neutral');
    assert.dom(CUSTOM_MESSAGES.badge('2')).hasClass('hds-badge--color-highlight');
  });
});
