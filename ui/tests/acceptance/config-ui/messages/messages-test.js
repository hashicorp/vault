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
import { PAGE as MESSAGE_SELECTORS } from 'vault/tests/helpers/config-ui/message-selectors';
import { SELECTORS as GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | config-ui', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.createMessage = async (messageType = 'banner', endTime = '2023-12-12', authenticated = true) => {
      await visit(`vault/config-ui/messages?authenticated=${authenticated}`);
      await click(MESSAGE_SELECTORS.button('create message'));
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Awesome custom message title');
      await click(MESSAGE_SELECTORS.radio(messageType));
      await fillIn(
        MESSAGE_SELECTORS.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn(
        MESSAGE_SELECTORS.input('startTime'),
        format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
      );
      if (endTime) {
        await click('#specificDate');
        await fillIn(
          MESSAGE_SELECTORS.input('endTime'),
          format(addDays(startOfDay(new Date('2023-12-12')), 10), datetimeLocalStringFormat)
        );
      }
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');

      await click(MESSAGE_SELECTORS.button('create-message'));
    };
    await authPage.login();
  });
  hooks.afterEach(async function () {
    await logout.visit();
  });

  test('it should show an empty state when no messages are created', async function (assert) {
    assert.expect(4);
    await visit('vault/config-ui/messages');
    await assert.dom('[data-test-component="empty-state"]').exists();
    await assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
    await click(MESSAGE_SELECTORS.tab('On login page'));
    await assert.dom('[data-test-component="empty-state"]').exists();
    await assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
  });

  module('Authenticated messages', function () {
    test('it should create, edit, view, and delete a message', async function (assert) {
      assert.expect(3);
      await this.createMessage();
      assert.dom(GENERAL.title).hasText('Awesome custom message title', 'on the details screen');
      await click('[data-test-link="edit"]');
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Edited custom message title');
      await click(MESSAGE_SELECTORS.button('create-message'));
      assert.dom(GENERAL.title).hasText('Edited custom message title');
      await click(MESSAGE_SELECTORS.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.messages.index',
        'redirects to messages page after delete'
      );
    });
    test('it should show multiple messages modal', async function (assert) {
      assert.expect(3);
      await this.createMessage('modal', null);
      assert.dom(GENERAL.title).hasText('Awesome custom message title');
      await this.createMessage('modal', null);
      assert.dom(MESSAGE_SELECTORS.modal('multiple modal messages')).exists();
      assert
        .dom(MESSAGE_SELECTORS.modalTitle('Warning: more than one modal'))
        .hasText('Warning: more than one modal after the user logs in');
      await click(MESSAGE_SELECTORS.modalButton('cancel'));
      await visit('vault/config-ui/messages');
      await click(MESSAGE_SELECTORS.listItem('Awesome custom message title'));
      await click(MESSAGE_SELECTORS.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
    });
    test('it should display preview a message when all required fields are filled out', async function (assert) {
      assert.expect(2);
      await visit(`vault/config-ui/messages`);
      await click(MESSAGE_SELECTORS.button('create message'));
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Awesome custom message title');
      await click(MESSAGE_SELECTORS.radio('banner'));
      await fillIn(
        MESSAGE_SELECTORS.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(MESSAGE_SELECTORS.button('preview'));
      assert.dom(MESSAGE_SELECTORS.modal('preview image')).exists();
      await click(MESSAGE_SELECTORS.modalButton('Close'));
      await click(MESSAGE_SELECTORS.radio('modal'));
      await click(MESSAGE_SELECTORS.button('preview'));
      assert.dom(MESSAGE_SELECTORS.modal('preview modal')).exists();
      await click(MESSAGE_SELECTORS.modalButton('Close'));
    });
    test('it should not display preview a message when all required fields are not filled out', async function (assert) {
      assert.expect(2);
      await visit(`vault/config-ui/messages`);
      await click(MESSAGE_SELECTORS.button('create message'));
      await click(MESSAGE_SELECTORS.radio('banner'));
      await fillIn(
        MESSAGE_SELECTORS.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(MESSAGE_SELECTORS.button('preview'));
      assert.dom(MESSAGE_SELECTORS.modal('preview image')).doesNotExist();
      assert.dom(MESSAGE_SELECTORS.input('title')).hasClass('has-error-border');
    });
  });

  module('Unauthenticated messages', function () {
    test('it should create, edit, view, and delete a message', async function (assert) {
      assert.expect(3);
      await this.createMessage('banner', null, false);
      assert.dom(GENERAL.title).hasText('Awesome custom message title', 'on the details screen');
      await click('[data-test-link="edit"]');
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Edited custom message title');
      await click(MESSAGE_SELECTORS.button('create-message'));
      assert.dom(GENERAL.title).hasText('Edited custom message title');
      await click(MESSAGE_SELECTORS.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.messages.index',
        'redirects to messages page after delete'
      );
    });
    test('it should show multiple messages modal', async function (assert) {
      assert.expect(3);
      await this.createMessage('modal', null, false);
      assert.dom(GENERAL.title).hasText('Awesome custom message title');
      await this.createMessage('modal', null, false);
      assert.dom(MESSAGE_SELECTORS.modal('multiple modal messages')).exists();
      assert
        .dom(MESSAGE_SELECTORS.modalTitle('Warning: more than one modal'))
        .hasText('Warning: more than one modal on the login page');
      await click(MESSAGE_SELECTORS.modalButton('cancel'));
      await visit('vault/config-ui/messages?authenticated=false');
      await click(MESSAGE_SELECTORS.listItem('Awesome custom message title'));
      await click(MESSAGE_SELECTORS.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
    });
    test('it should show info message on create and edit form', async function (assert) {
      assert.expect(1);
      await visit('vault/config-ui/messages?authenticated=false');
      await click(MESSAGE_SELECTORS.button('create message'));
      assert
        .dom(MESSAGE_SELECTORS.unauthCreateFormInfo)
        .hasText(
          'Note: Do not include sensitive info in this message since users are unauthenticated at this stage.'
        );
    });
    test('it should display preview a message when all required fields are filled out', async function (assert) {
      assert.expect(2);
      await visit('vault/config-ui/messages?authenticated=false');
      await click(MESSAGE_SELECTORS.button('create message'));
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Awesome custom message title');
      await click(MESSAGE_SELECTORS.radio('banner'));
      await fillIn(
        MESSAGE_SELECTORS.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(MESSAGE_SELECTORS.button('preview'));
      assert.dom(MESSAGE_SELECTORS.modal('preview image')).exists();
      await click(MESSAGE_SELECTORS.modalButton('Close'));
      await click(MESSAGE_SELECTORS.radio('modal'));
      await click(MESSAGE_SELECTORS.button('preview'));
      assert.dom(MESSAGE_SELECTORS.modal('preview modal')).exists();
      await click(MESSAGE_SELECTORS.modalButton('Close'));
    });
    test('it should not display preview a message when all required fields are not filled out', async function (assert) {
      assert.expect(2);
      await visit('vault/config-ui/messages?authenticated=false');
      await click(MESSAGE_SELECTORS.button('create message'));
      await click(MESSAGE_SELECTORS.radio('banner'));
      await fillIn(
        MESSAGE_SELECTORS.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(MESSAGE_SELECTORS.button('preview'));
      assert.dom(MESSAGE_SELECTORS.modal('preview image')).doesNotExist();
      assert.dom(MESSAGE_SELECTORS.input('title')).hasClass('has-error-border');
    });
  });
});
