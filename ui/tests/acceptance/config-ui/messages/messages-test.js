/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit, fillIn, currentRouteName, findAll } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
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
    await login();
  });

  hooks.afterEach(async function () {
    await visit('/vault/logout');
  });

  test('it should hide the sidebar settings section on community', async function (assert) {
    assert.dom(CUSTOM_MESSAGES.navLink).doesNotExist();
  });
});

module('Acceptance | Enterprise | config-ui/message', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';
    await login();

    this.createMessage = async (
      messageId,
      messageType = 'banner',
      endTime = '2023-12-12',
      authenticated = true
    ) => {
      await click(CUSTOM_MESSAGES.navLink);
      await click(CUSTOM_MESSAGES.tab(authenticated ? 'After user logs in' : 'On login page'));
      await click(GENERAL.submitButton);

      await fillIn(CUSTOM_MESSAGES.input('title'), messageId);
      await click(CUSTOM_MESSAGES.radio(messageType));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
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

      await click(GENERAL.submitButton);
    };
  });

  hooks.afterEach(async function () {
    await visit('/vault/logout');
  });

  test('it should show an empty state when no messages are created', async function (assert) {
    await click(CUSTOM_MESSAGES.navLink);
    assert.dom(GENERAL.emptyStateTitle).exists();
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
    await click(CUSTOM_MESSAGES.tab('On login page'));
    assert.dom(GENERAL.emptyStateTitle).exists();
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
  });

  test('authenticated it should create, edit, view, and delete a message', async function (assert) {
    // create first message
    await this.createMessage('new-message');
    assert.dom(GENERAL.title).hasText('new-message', 'message title shows on the details screen');
    // edit message
    await click(GENERAL.linkTo('edit'));
    await fillIn(CUSTOM_MESSAGES.input('title'), `Edited new-message`);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.title)
      .hasText(`Edited new-message`, 'edited message title shows on the details screen');
    // delete the edited message
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    const linkedBlocks = findAll('[data-test-list-item]');
    assert.false(linkedBlocks.includes(`Edited new-message`), 'edited message was deleted.');
  });

  test('authenticated it should show multiple messages modal', async function (assert) {
    await this.createMessage('message-one', 'modal', null);
    // create second message with same model name
    await this.createMessage('message-one', 'modal', null);
    assert.dom(CUSTOM_MESSAGES.modal('multiple modal messages')).exists();
    assert
      .dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal'))
      .hasText('Warning: more than one modal after the user logs in');

    await click(GENERAL.button('cancel-multiple')); // cancel out of the modal
    await click(GENERAL.cancelButton); // cancel out of the create message form
    // delete the created message
    await click(CUSTOM_MESSAGES.listItem('message-one')); // go back to the message list
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
  });

  test('it should filter by type and status', async function (assert) {
    await this.createMessage('filter-status-1', 'banner', null);
    await this.createMessage('filter-status-2', 'banner');
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
    await click(GENERAL.submitButton);
    assert.dom(MESSAGES_LIST.listItem).exists({ count: 1 }, 'list filters by status');

    // check number of messages with type filters
    await click(MESSAGES_LIST.filterReset);
    await fillIn(MESSAGES_LIST.filterBy('type'), 'modal');
    await click(GENERAL.submitButton);
    // because of test pollution, we cannot guarantee that the list will be empty
    // make sure only modal messages or no messages are shown
    const messages = findAll(MESSAGES_LIST.listItem);
    const allMessages = Array.from(messages || []);
    const modalMessages = allMessages.filter((node) => node.querySelector('[data-test-badge="modal"]'));

    const hasMessages = allMessages.length > 0;

    assert.strictEqual(
      modalMessages.length,
      hasMessages ? allMessages.length : 0,
      'if there are items in the list, they are modal messages'
    );
    // unsetting a filter will reset that item in the query
    await fillIn(MESSAGES_LIST.filterBy('type'), '');
    await fillIn(MESSAGES_LIST.filterBy('status'), 'inactive');
    await click(GENERAL.submitButton);
    assert.dom(MESSAGES_LIST.listItem).exists({ count: 1 }, 'list filters by status again');

    // delete the created messages
    await click(MESSAGES_LIST.filterReset);
    await click(CUSTOM_MESSAGES.listItem('filter-status-1'));
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    await click(CUSTOM_MESSAGES.listItem('filter-status-2'));
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
  });

  test('it should display preview a message when all required fields are filled out', async function (assert) {
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('After user logs in'));
    await click(GENERAL.submitButton);
    await fillIn(CUSTOM_MESSAGES.input('title'), 'authenticated display preview');
    await click(CUSTOM_MESSAGES.radio('banner'));
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview image')).exists('preview image of the message shows');

    await click(GENERAL.button('Close preview'));
    await click(CUSTOM_MESSAGES.radio('modal'));
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview modal')).exists('preview modal of the message shows');
  });

  test('it should not display preview a message when all required fields are not filled out', async function (assert) {
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('After user logs in'));
    await click(GENERAL.submitButton);
    await click(CUSTOM_MESSAGES.radio('banner'));
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(GENERAL.button('preview'));
    assert
      .dom(CUSTOM_MESSAGES.modal('preview image'))
      .doesNotExist('preview image does not show because you have a missing title');
    assert.dom(CUSTOM_MESSAGES.input('title')).hasClass('has-error-border', 'error around title shows');
  });

  // unauthenticated messages
  test('unauthenticated it should create, edit, view, and delete a message', async function (assert) {
    await this.createMessage('unauthenticated create edit view delete', 'banner', null, false);
    assert
      .dom(GENERAL.title)
      .hasText('unauthenticated create edit view delete', 'title shows on the details screen');
    // navigate to edit the title
    await click(GENERAL.linkTo('edit'));
    await fillIn(CUSTOM_MESSAGES.input('title'), `Edited ${'unauthenticated create edit view delete'}`);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.title)
      .hasText(
        `Edited ${'unauthenticated create edit view delete'}`,
        'edited title shows on the details screen'
      );
    // delete the edited message
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.config-ui.messages.index',
      'redirects to messages page after delete'
    );
  });

  test('unauthenticated it should show multiple messages modal', async function (assert) {
    await this.createMessage('unauthenticated message 1', 'modal', null, false);
    // create second message with same model name
    await this.createMessage('unauthenticated message 1', 'modal', null, false);
    assert.dom(CUSTOM_MESSAGES.modal('multiple modal messages')).exists('the multiple modal message shows');
    assert
      .dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal'))
      .hasText('Warning: more than one modal on the login page', 'the warning modal title shows');

    await click(GENERAL.button('cancel-multiple')); // cancel out of the modal
    await click(GENERAL.cancelButton); // cancel out of the create message form
    // delete the created message
    await click(CUSTOM_MESSAGES.listItem('unauthenticated message 1')); // go back to the message list
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
  });

  test('it should show info message about sensitive information on create and edit form', async function (assert) {
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('On login page'));
    await click(GENERAL.submitButton);
    assert
      .dom(CUSTOM_MESSAGES.unauthCreateFormInfo)
      .hasText(
        'Note: Do not include sensitive information in this message since users are unauthenticated at this stage.'
      );
  });

  test('it should allow you to preview a message when all required fields are filled out', async function (assert) {
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('On login page'));
    await click(GENERAL.submitButton);
    await fillIn(CUSTOM_MESSAGES.input('title'), 'unauthenticated display preview');
    await click(CUSTOM_MESSAGES.radio('banner'));
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview image')).exists('preview image of the message shows');

    await click(GENERAL.button('Close preview'));
    await click(CUSTOM_MESSAGES.radio('modal'));
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview modal')).exists('preview modal of the message shows');
  });

  test('it should not display a preview of a message when all required fields are not filled out', async function (assert) {
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('On login page'));
    await click(GENERAL.submitButton);
    await click(CUSTOM_MESSAGES.radio('banner'));
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview image')).doesNotExist('preview image does not show');
    assert.dom(CUSTOM_MESSAGES.input('title')).hasClass('has-error-border', 'error around title shows');
  });
});
