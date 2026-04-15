/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupEngine } from 'ember-engines/test-support';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { dateFormat } from 'core/helpers/date-format';

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
        start_time: new Date('2021-08-01T00:00:00Z'),
        end_time: undefined,
      },
      {
        id: '1',
        active: false,
        type: 'modal',
        authenticated: true,
        title: 'Message title 2',
        message: 'Some long long long message blah blah blah',
        link: { title: 'here', href: 'www.example2.com' },
        start_time: new Date('2023-07-01T00:00:00Z'),
        end_time: new Date('2023-08-01T00:00:00Z'),
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
      {
        id: '3',
        active: true,
        type: 'banner',
        authenticated: true,
        title: 'Message title 4',
        message: 'Current event message with timezone',
        link: { title: 'here', href: 'www.example.com' },
        start_time: new Date('2023-07-01T00:00:00Z'),
        end_time: new Date('2023-08-01T00:00:00Z'),
      },
      {
        id: '4',
        active: false,
        type: 'banner',
        authenticated: true,
        title: 'Message title 5',
        message: 'A message from the future',
        link: { title: 'here', href: 'www.example.com' },
        start_time: new Date(Date.now() + 1000000000), // start time in the future
        end_time: undefined,
      },
    ];
    // Pass meta information by value so that each test can modify its own copy
    // without affecting other tests
    Object.defineProperty(this.messages, 'meta', { value: { ...META.value } });

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

    this.formatTime = (time) => dateFormat([time, 'MMM d, yyyy hh:mm aaa'], { withTimeZone: true });
  });

  test('it should show the messages empty state', async function (assert) {
    this.messages = [];

    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Add a custom message for all users after they log into Vault. Create message to get started.'
      );
  });

  test('it should show the list of custom messages', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.icon('message-circle')).exists();
    for (const message of this.messages) {
      assert.dom(GENERAL.listItem(message.title)).exists('Message title is displayed');
    }
  });

  test('it should show max message warning modal', async function (assert) {
    // Reset messages to 100 to test max message limit
    this.messages = [];
    for (let i = 0; i < 100; i++) {
      this.messages.push({
        id: `${i}-a`,
        active: true,
        type: 'banner',
        authenticated: false,
        title: `Message title ${i}`,
        message: 'Some long long long message',
        link: { title: 'here', href: 'www.example.com' },
        start_time: new Date('2021-08-01T00:00:00Z'),
      });
    }

    Object.defineProperty(this.messages, 'meta', {
      value: {
        currentPage: 1,
        lastPage: 1,
        nextPage: 1,
        prevPage: 1,
        total: this.messages.length,
        pageSize: 100,
      },
    });

    await this.renderComponent();
    await click(GENERAL.button('Create message'));
    assert
      .dom(CUSTOM_MESSAGES.modalTitle('maximum-message-modal'))
      .hasText('Maximum number of messages reached');
    assert
      .dom(CUSTOM_MESSAGES.modalBody('maximum-message-modal'))
      .hasText(
        'Vault can only store up to 100 messages. To create a message, delete one of your messages to clear up space.'
      );
    await click(GENERAL.button('close-maximum-message'));
  });

  // Badge tests
  test('it should show active custom messages', async function (assert) {
    await this.renderComponent();

    const activeMessage = this.messages[0];
    assert
      .dom(` ${GENERAL.listItem(activeMessage.title)} ${GENERAL.badge()}`)
      .hasText('Active')
      .hasClass('hds-badge--color-success');
  });

  test('it should show active custom messages with end time if present', async function (assert) {
    await this.renderComponent();

    const activeMessageWithTimeZone = this.messages[3];
    assert
      .dom(` ${GENERAL.listItem(activeMessageWithTimeZone.title)} ${GENERAL.badge()}`)
      .hasText(`Active until ${this.formatTime(activeMessageWithTimeZone.end_time)}`)
      .hasClass('hds-badge--color-success');
  });

  test('it should show scheduled messages with future start date and time', async function (assert) {
    await this.renderComponent();

    const scheduledMessage = this.messages[4];
    assert
      .dom(` ${GENERAL.listItem(scheduledMessage.title)} ${GENERAL.badge()}`)
      .hasText(`Scheduled: ${this.formatTime(scheduledMessage.start_time)}`)
      .hasClass('hds-badge--color-highlight');
  });

  test('it should show inactive messages with message expiration time', async function (assert) {
    await this.renderComponent();

    const inactiveMessage = this.messages[1];
    assert
      .dom(` ${GENERAL.listItem(inactiveMessage.title)} ${GENERAL.badge()}`)
      .hasText(`Inactive: ${this.formatTime(inactiveMessage.end_time)}`)
      .hasClass('hds-badge--color-neutral');
  });
});
