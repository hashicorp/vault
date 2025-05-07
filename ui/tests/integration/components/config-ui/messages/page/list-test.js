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
import { addDays, startOfDay } from 'date-fns';
import timestamp from 'core/utils/timestamp';

const META = {
  value: {
    currentPage: 1,
    lastPage: 1,
    nextPage: 1,
    prevPage: 1,
    total: 3,
    pageSize: 15,
  },
};

module('Integration | Component | messages/page/list', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'config-ui');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.messages = [
      {
        id: '0',
        active: true,
        type: 'banner',
        authenticated: true,
        title: 'Message title 1',
        message: 'Some long long long message',
        link: { title: 'here', href: 'www.example.com' },
        startTime: new Date('2021-08-01T00:00:00Z'),
        endTime: undefined,
      },
      {
        id: '1',
        active: false,
        type: 'modal',
        authenticated: true,
        title: 'Message title 2',
        message: 'Some long long long message blah blah blah',
        link: { title: 'here', href: 'www.example2.com' },
        startTime: new Date('2023-07-01T00:00:00Z'),
        endTime: new Date('2023-08-01T00:00:00Z'),
      },
      {
        id: '2',
        active: false,
        type: 'banner',
        authenticated: false,
        title: 'Message title 3',
        message: 'Some long long long message',
        link: { title: 'here', href: 'www.example.com' },
      },
    ];
    Object.defineProperty(this.messages, 'meta', META);

    this.renderComponent = () => {
      const capabilitiesService = this.owner.lookup('service:capabilities');
      this.capabilities = this.messages.reduce((obj, { id }) => {
        const path = capabilitiesService.pathFor('customMessages', { id });
        obj[path] = { canUpdate: true, canDelete: true };
        return obj;
      }, {});

      return render(
        hbs`<Messages::Page::List @messages={{this.messages}} @capabilities={{this.capabilities}} />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should show the messages empty state', async function (assert) {
    this.messages = [];

    await this.renderComponent();
    assert.dom('[data-test-empty-state-title]').hasText('No messages yet');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'Add a custom message for all users after they log into Vault. Create message to get started.'
      );
  });

  test('it should show the list of custom messages', async function (assert) {
    await this.renderComponent();
    assert.dom('[data-test-icon="message-circle"]').exists();
    for (const message of this.messages) {
      assert.dom(CUSTOM_MESSAGES.listItem('Message title 1')).exists();
      assert.dom(`[data-linked-block-title="${message.id}"]`).hasText(message.title);
    }
  });

  test('it should show max message warning modal', async function (assert) {
    for (let i = 0; i < 97; i++) {
      this.messages.push({
        id: `${i}-a`,
        active: true,
        type: 'banner',
        authenticated: false,
        title: `Message title ${i}`,
        message: 'Some long long long message',
        link: { title: 'here', href: 'www.example.com' },
        startTime: new Date('2021-08-01T00:00:00Z'),
      });
    }
    this.messages.meta.total = this.messages.length;
    this.messages.meta.pageSize = 100;

    await this.renderComponent();
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
    this.messages[2].startTime = addDays(startOfDay(timestamp.now()), 1);

    await this.renderComponent();
    assert.dom(CUSTOM_MESSAGES.badge('0')).hasClass('hds-badge--color-success');
    assert.dom(CUSTOM_MESSAGES.badge('1')).hasClass('hds-badge--color-neutral');
    assert.dom(CUSTOM_MESSAGES.badge('2')).hasClass('hds-badge--color-highlight');
  });
});
