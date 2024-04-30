/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit, fillIn, currentRouteName, currentURL } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import { format, addDays, startOfDay } from 'date-fns';
import { datetimeLocalStringFormat } from 'core/utils/date-formatters';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const MESSAGES_LIST = {
  listItem: '.linked-block',
  filterBy: (name) => `[data-test-filter-by="${name}"]`,
  filterSubmit: '[data-test-filter-submit]',
  filterReset: '[data-test-filter-reset]',
};

module('Acceptance | Community | config-ui/messages', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const version = this.owner.lookup('service:version');
    version.type = 'community';
    await authPage.login();
  });

  hooks.afterEach(async function () {
    await logout.visit();
  });

  test('it should hide the sidebar settings section on community', async function (assert) {
    assert.expect(1);
    assert.dom(CUSTOM_MESSAGES.navLink).doesNotExist();
  });
});

module('Acceptance | Enterprise | config-ui/message', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';
    this.messageDetailId = () => {
      return currentURL().match(/messages\/(.*)\/details/)[1];
    };
    this.createMessage = async (messageType = 'banner', endTime = '2023-12-12', authenticated = true) => {
      await click(CUSTOM_MESSAGES.navLink);
      if (authenticated) {
        await click(CUSTOM_MESSAGES.tab('After user logs in'));
      } else {
        await click(CUSTOM_MESSAGES.tab('On login page'));
      }
      await click(CUSTOM_MESSAGES.button('create message'));

      await fillIn(CUSTOM_MESSAGES.input('title'), 'Awesome custom message title');
      await click(CUSTOM_MESSAGES.radio(messageType));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn(
        CUSTOM_MESSAGES.input('startTime'),
        format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
      );
      if (endTime) {
        await click('#specificDate');
        await fillIn(
          CUSTOM_MESSAGES.input('endTime'),
          format(addDays(startOfDay(new Date('2023-12-12')), 10), datetimeLocalStringFormat)
        );
      }
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');

      await click(CUSTOM_MESSAGES.button('create-message'));
    };
    this.deleteMessage = async (id) => {
      await visit(`vault/config-ui/messages/${id}/details`);
      await click(CUSTOM_MESSAGES.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
    };
    await authPage.login();
  });

  hooks.afterEach(async function () {
    await logout.visit();
  });
  test('it should show an empty state when no messages are created', async function (assert) {
    assert.expect(4);
    await click(CUSTOM_MESSAGES.navLink);
    assert.dom(GENERAL.emptyStateTitle).exists();
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
    await click(CUSTOM_MESSAGES.tab('On login page'));
    assert.dom(GENERAL.emptyStateTitle).exists();
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
  });

  module('Authenticated messages', function () {
    test('it should create, edit, view, and delete a message', async function (assert) {
      assert.expect(3);
      await this.createMessage();
      assert.dom(GENERAL.title).hasText('Awesome custom message title', 'on the details screen');
      await click('[data-test-link="edit"]');
      await fillIn(CUSTOM_MESSAGES.input('title'), 'Edited custom message title');
      await click(CUSTOM_MESSAGES.button('create-message'));
      assert.dom(GENERAL.title).hasText('Edited custom message title');
      await click(CUSTOM_MESSAGES.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.messages.index',
        'redirects to messages page after delete'
      );
    });

    test('it should show multiple messages modal', async function (assert) {
      assert.expect(4);
      await this.createMessage('modal', null);
      assert.dom(GENERAL.title).hasText('Awesome custom message title');
      await this.createMessage('modal', null);
      assert.dom(CUSTOM_MESSAGES.modal('multiple modal messages')).exists();
      assert
        .dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal'))
        .hasText('Warning: more than one modal after the user logs in');
      await click(CUSTOM_MESSAGES.modalButton('cancel'));
      await visit('vault/config-ui/messages');
      await click(CUSTOM_MESSAGES.listItem('Awesome custom message title'));
      await click(CUSTOM_MESSAGES.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.emptyStateTitle).exists('Message was deleted');
    });
    test('it should filter by type and status', async function (assert) {
      await this.createMessage('banner', null);
      const msg1 = this.messageDetailId();
      await this.createMessage('banner');
      const msg2 = this.messageDetailId();
      await visit('vault/config-ui/messages?pageFilter=foobar&status=inactive&type=banner');
      // check that filters inherit param values
      assert.dom(MESSAGES_LIST.filterBy('pageFilter')).hasValue('foobar');
      assert.dom(MESSAGES_LIST.filterBy('status')).hasValue('inactive');
      assert.dom(MESSAGES_LIST.filterBy('type')).hasValue('banner');
      assert.dom(GENERAL.emptyStateTitle).exists();

      // clear filters works
      await click(MESSAGES_LIST.filterReset);
      assert.dom(MESSAGES_LIST.listItem).exists({ count: 2 });

      // check number of messages with status filters
      await fillIn(MESSAGES_LIST.filterBy('status'), 'active');
      assert.dom(MESSAGES_LIST.listItem).exists({ count: 2 }, 'list does not filter before clicking submit');
      await click(MESSAGES_LIST.filterSubmit);
      assert.dom(MESSAGES_LIST.listItem).exists({ count: 1 });

      // check number of messages with type filters
      await click(MESSAGES_LIST.filterReset);
      await fillIn(MESSAGES_LIST.filterBy('type'), 'modal');
      await click(MESSAGES_LIST.filterSubmit);
      assert.dom(GENERAL.emptyStateTitle).exists();

      // unsetting a filter will reset that item in the query
      await fillIn(MESSAGES_LIST.filterBy('type'), '');
      await fillIn(MESSAGES_LIST.filterBy('status'), 'inactive');
      await click(MESSAGES_LIST.filterSubmit);
      assert.dom(MESSAGES_LIST.listItem).exists({ count: 1 });

      // clean up custom messages
      await this.deleteMessage(msg1);
      await this.deleteMessage(msg2);
    });
    test('it should display preview a message when all required fields are filled out', async function (assert) {
      assert.expect(2);
      await click(CUSTOM_MESSAGES.navLink);
      await click(CUSTOM_MESSAGES.tab('After user logs in'));
      await click(CUSTOM_MESSAGES.button('create message'));
      await fillIn(CUSTOM_MESSAGES.input('title'), 'Awesome custom message title');
      await click(CUSTOM_MESSAGES.radio('banner'));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(CUSTOM_MESSAGES.button('preview'));
      assert.dom(CUSTOM_MESSAGES.modal('preview image')).exists();
      await click(CUSTOM_MESSAGES.modalButton('Close'));
      await click(CUSTOM_MESSAGES.radio('modal'));
      await click(CUSTOM_MESSAGES.button('preview'));
      assert.dom(CUSTOM_MESSAGES.modal('preview modal')).exists();
    });
    test('it should not display preview a message when all required fields are not filled out', async function (assert) {
      assert.expect(2);
      await click(CUSTOM_MESSAGES.navLink);
      await click(CUSTOM_MESSAGES.tab('After user logs in'));
      await click(CUSTOM_MESSAGES.button('create message'));
      await click(CUSTOM_MESSAGES.radio('banner'));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(CUSTOM_MESSAGES.button('preview'));
      assert.dom(CUSTOM_MESSAGES.modal('preview image')).doesNotExist();
      assert.dom(CUSTOM_MESSAGES.input('title')).hasClass('has-error-border');
    });
  });

  module('Unauthenticated messages', function () {
    test('it should create, edit, view, and delete a message', async function (assert) {
      assert.expect(3);
      await this.createMessage('banner', null, false);
      assert.dom(GENERAL.title).hasText('Awesome custom message title', 'on the details screen');
      await click('[data-test-link="edit"]');
      await fillIn(CUSTOM_MESSAGES.input('title'), 'Edited custom message title');
      await click(CUSTOM_MESSAGES.button('create-message'));
      assert.dom(GENERAL.title).hasText('Edited custom message title');
      await click(CUSTOM_MESSAGES.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.messages.index',
        'redirects to messages page after delete'
      );
    });
    test('it should show multiple messages modal', async function (assert) {
      assert.expect(4);
      await this.createMessage('modal', null, false);
      assert.dom(GENERAL.title).hasText('Awesome custom message title');
      await this.createMessage('modal', null, false);
      assert.dom(CUSTOM_MESSAGES.modal('multiple modal messages')).exists();
      assert
        .dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal'))
        .hasText('Warning: more than one modal on the login page');
      await click(CUSTOM_MESSAGES.modalButton('cancel'));
      await visit('vault/config-ui/messages?authenticated=false');
      await click(CUSTOM_MESSAGES.listItem('Awesome custom message title'));
      await click(CUSTOM_MESSAGES.confirmActionButton('Delete message'));
      await click(GENERAL.confirmButton);
      assert.dom(GENERAL.emptyStateTitle).exists('Message was deleted');
    });
    test('it should show info message on create and edit form', async function (assert) {
      assert.expect(1);
      await click(CUSTOM_MESSAGES.navLink);
      await click(CUSTOM_MESSAGES.tab('On login page'));
      await click(CUSTOM_MESSAGES.button('create message'));
      assert
        .dom(CUSTOM_MESSAGES.unauthCreateFormInfo)
        .hasText(
          'Note: Do not include sensitive information in this message since users are unauthenticated at this stage.'
        );
    });
    test('it should display preview a message when all required fields are filled out', async function (assert) {
      assert.expect(2);
      await click(CUSTOM_MESSAGES.navLink);
      await click(CUSTOM_MESSAGES.tab('On login page'));
      await click(CUSTOM_MESSAGES.button('create message'));
      await fillIn(CUSTOM_MESSAGES.input('title'), 'Awesome custom message title');
      await click(CUSTOM_MESSAGES.radio('banner'));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(CUSTOM_MESSAGES.button('preview'));
      assert.dom(CUSTOM_MESSAGES.modal('preview image')).exists();
      await click(CUSTOM_MESSAGES.modalButton('Close'));
      await click(CUSTOM_MESSAGES.radio('modal'));
      await click(CUSTOM_MESSAGES.button('preview'));
      assert.dom(CUSTOM_MESSAGES.modal('preview modal')).exists();
    });
    test('it should not display preview a message when all required fields are not filled out', async function (assert) {
      assert.expect(2);
      await click(CUSTOM_MESSAGES.navLink);
      await click(CUSTOM_MESSAGES.tab('On login page'));
      await click(CUSTOM_MESSAGES.button('create message'));
      await click(CUSTOM_MESSAGES.radio('banner'));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn('[data-test-kv-key="0"]', 'Learn more');
      await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
      await click(CUSTOM_MESSAGES.button('preview'));
      assert.dom(CUSTOM_MESSAGES.modal('preview image')).doesNotExist();
      assert.dom(CUSTOM_MESSAGES.input('title')).hasClass('has-error-border');
    });
  });
});
