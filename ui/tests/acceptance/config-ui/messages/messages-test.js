/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit, fillIn, currentRouteName } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import { format, addDays, startOfDay } from 'date-fns';
import { datetimeLocalStringFormat } from 'core/utils/date-formatters';
import { PAGE } from 'vault/tests/helpers/config-ui/message-selectors';
import { clickTrigger } from 'ember-power-select/test-support/helpers';

module('Acceptance | config-ui', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.createMessage = async (messageType = 'banner', endTime = '2023-12-12', authenticated = true) => {
      await visit(`vault/config-ui/messages?authenticated=${authenticated}`);
      await click(PAGE.button('create message'));
      await fillIn(PAGE.input('title'), 'Awesome custom message title');
      await click(PAGE.radio(messageType));
      await fillIn(
        PAGE.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn(
        PAGE.input('startTime'),
        format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
      );
      if (endTime) {
        await click('#specificDate');
        await fillIn(
          PAGE.input('endTime'),
          format(addDays(startOfDay(new Date('2023-12-12')), 10), datetimeLocalStringFormat)
        );
      }
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');

      await click(PAGE.button('create-message'));
    };
    await authPage.login();
  });
  hooks.afterEach(async function () {
    await logout.visit();
  });

  test('it should show an empty state when no messages are created', async function (assert) {
    assert.expect(4);
    await visit('vault/config-ui/messages');
    assert.dom('[data-test-component="empty-state"]').exists();
    assert.dom(PAGE.emptyStateTitle).hasText('No messages yet');
    await click(PAGE.tab('On login page'));
    assert.dom('[data-test-component="empty-state"]').exists();
    assert.dom(PAGE.emptyStateTitle).hasText('No messages yet');
  });

  module('Authenticated messages', function () {
    test('it should create, edit, view, and delete a message', async function (assert) {
      assert.expect(3);
      await this.createMessage();
      assert.dom(PAGE.title).hasText('Awesome custom message title', 'on the details screen');
      await click('[data-test-link="edit"]');
      await fillIn(PAGE.input('title'), 'Edited custom message title');
      await click(PAGE.button('create-message'));
      assert.dom(PAGE.title).hasText('Edited custom message title');
      await click(PAGE.confirmActionButton('Delete message'));
      await click(PAGE.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.messages.index',
        'redirects to messages page after delete'
      );
    });

    test('it should show multiple messages modal', async function (assert) {
      assert.expect(4);
      await this.createMessage('modal', null);
      assert.dom(PAGE.title).hasText('Awesome custom message title');
      await this.createMessage('modal', null);
      assert.dom(PAGE.modal('multiple modal messages')).exists();
      assert
        .dom(PAGE.modalTitle('Warning: more than one modal'))
        .hasText('Warning: more than one modal after the user logs in');
      await click(PAGE.modalButton('cancel'));
      await visit('vault/config-ui/messages');
      await click(PAGE.listItem('Awesome custom message title'));
      await click(PAGE.confirmActionButton('Delete message'));
      await click(PAGE.confirmButton);
      assert.dom('[data-test-component="empty-state"]').exists('Message was deleted');
    });
    test('it should filter by type and status', async function (assert) {
      assert.expect(6);
      await this.createMessage('banner', null);
      await this.createMessage('banner');
      await visit('vault/config-ui/messages');

      // check number of messages with status filters
      await clickTrigger('#filter-by-message-status');
      await click('.ember-power-select-options [data-option-index="0"]');
      assert.dom('.linked-block').exists({ count: 1 }, 'filtered by active');
      await click('[data-test-selected-list-button="delete"]');
      await clickTrigger('#filter-by-message-status');
      await click('.ember-power-select-options [data-option-index="1"]');
      assert.dom('.linked-block').exists({ count: 1 }, 'filtered by inactive');
      await click('[data-test-selected-list-button="delete"]');

      // check number of messages with type filters
      await clickTrigger('#filter-by-message-type');
      await click('.ember-power-select-options [data-option-index="0"]');
      assert.dom('.linked-block').exists({ count: 0 }, 'filtered by modal');
      await click('[data-test-selected-list-button="delete"]');
      await clickTrigger('#filter-by-message-type');
      await click('.ember-power-select-options [data-option-index="1"]');
      assert.dom('.linked-block').exists({ count: 2 }, 'filtered by banner');
      await click('[data-test-selected-list-button="delete"]');

      // check number of messages with no filters
      assert.dom('.linked-block').exists({ count: 2 }, 'no filters selected');

      // clean up custom messages
      await click(PAGE.listItem('Awesome custom message title'));
      await click(PAGE.confirmActionButton('Delete message'));
      await click(PAGE.confirmButton);
      await click(PAGE.listItem('Awesome custom message title'));
      await click(PAGE.confirmActionButton('Delete message'));
      await click(PAGE.confirmButton);
      assert.dom('[data-test-component="empty-state"]').exists('Message was deleted');
    });
    test('it should display preview a message when all required fields are filled out', async function (assert) {
      assert.expect(2);
      await visit(`vault/config-ui/messages`);
      await click(PAGE.button('create message'));
      await fillIn(PAGE.input('title'), 'Awesome custom message title');
      await click(PAGE.radio('banner'));
      await fillIn(
        PAGE.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(PAGE.button('preview'));
      assert.dom(PAGE.modal('preview image')).exists();
      await click(PAGE.modalButton('Close'));
      await click(PAGE.radio('modal'));
      await click(PAGE.button('preview'));
      assert.dom(PAGE.modal('preview modal')).exists();
    });
    test('it should not display preview a message when all required fields are not filled out', async function (assert) {
      assert.expect(2);
      await visit(`vault/config-ui/messages`);
      await click(PAGE.button('create message'));
      await click(PAGE.radio('banner'));
      await fillIn(
        PAGE.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(PAGE.button('preview'));
      assert.dom(PAGE.modal('preview image')).doesNotExist();
      assert.dom(PAGE.input('title')).hasClass('has-error-border');
    });
  });

  module('Unauthenticated messages', function () {
    test('it should create, edit, view, and delete a message', async function (assert) {
      assert.expect(3);
      await this.createMessage('banner', null, false);
      assert.dom(PAGE.title).hasText('Awesome custom message title', 'on the details screen');
      await click('[data-test-link="edit"]');
      await fillIn(PAGE.input('title'), 'Edited custom message title');
      await click(PAGE.button('create-message'));
      assert.dom(PAGE.title).hasText('Edited custom message title');
      await click(PAGE.confirmActionButton('Delete message'));
      await click(PAGE.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.messages.index',
        'redirects to messages page after delete'
      );
    });
    test('it should show multiple messages modal', async function (assert) {
      assert.expect(4);
      await this.createMessage('modal', null, false);
      assert.dom(PAGE.title).hasText('Awesome custom message title');
      await this.createMessage('modal', null, false);
      assert.dom(PAGE.modal('multiple modal messages')).exists();
      assert
        .dom(PAGE.modalTitle('Warning: more than one modal'))
        .hasText('Warning: more than one modal on the login page');
      await click(PAGE.modalButton('cancel'));
      await visit('vault/config-ui/messages?authenticated=false');
      await click(PAGE.listItem('Awesome custom message title'));
      await click(PAGE.confirmActionButton('Delete message'));
      await click(PAGE.confirmButton);
      assert.dom('[data-test-component="empty-state"]').exists('Message was deleted');
    });
    test('it should show info message on create and edit form', async function (assert) {
      assert.expect(1);
      await visit('vault/config-ui/messages?authenticated=false');
      await click(PAGE.button('create message'));
      assert
        .dom(PAGE.unauthCreateFormInfo)
        .hasText(
          'Note: Do not include sensitive information in this message since users are unauthenticated at this stage.'
        );
    });
    test('it should display preview a message when all required fields are filled out', async function (assert) {
      assert.expect(2);
      await visit('vault/config-ui/messages?authenticated=false');
      await click(PAGE.button('create message'));
      await fillIn(PAGE.input('title'), 'Awesome custom message title');
      await click(PAGE.radio('banner'));
      await fillIn(
        PAGE.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(PAGE.button('preview'));
      assert.dom(PAGE.modal('preview image')).exists();
      await click(PAGE.modalButton('Close'));
      await click(PAGE.radio('modal'));
      await click(PAGE.button('preview'));
      assert.dom(PAGE.modal('preview modal')).exists();
    });
    test('it should not display preview a message when all required fields are not filled out', async function (assert) {
      assert.expect(2);
      await visit('vault/config-ui/messages?authenticated=false');
      await click(PAGE.button('create message'));
      await click(PAGE.radio('banner'));
      await fillIn(
        PAGE.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(PAGE.button('preview'));
      assert.dom(PAGE.modal('preview image')).doesNotExist();
      assert.dom(PAGE.input('title')).hasClass('has-error-border');
    });
  });
});
