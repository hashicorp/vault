/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit, fillIn, findAll, waitFor } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { format, addDays, startOfDay } from 'date-fns';
import { datetimeLocalStringFormat } from 'core/utils/date-formatters';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { encodeString } from 'core/utils/b64';

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
    await click(GENERAL.navLink('Operational tools'));
    assert.dom(CUSTOM_MESSAGES.navLink).doesNotExist();
  });
});

module('Acceptance | Enterprise | config-ui/message', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  const messageIdObject = {}; // Object that holds the message title as the key and the message ID as the value. using an object to help debug if a specific createMessage is causing issues.

  hooks.beforeEach(async function () {
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';
    await login();

    // Use the CLI command to create a message and retrieve its ID (assigned by the API on POST).
    // The message ID is required for cleanup to prevent test pollution.
    this.createMessageRepl = async ({
      title = 'Test Message',
      type = 'banner',
      message = encodeString('Lorem ipsum dolor sit amet, consectetur adipiscing elit.'),
      end_time = null,
      start_time = '2023-12-12T00:00:00.000Z',
      authenticated = true,
    } = {}) => {
      const payloadParts = [
        `title="${title}"`,
        `message="${message}"`,
        `type="${type}"`,
        `start_time="${start_time}"`,
        `authenticated=${authenticated}`,
      ];
      if (end_time) {
        payloadParts.push(`end_time="${end_time}"`);
      }
      const payload = payloadParts.join(' ');
      const result = await runCmd(`vault write sys/config/ui/custom-messages/ ${payload}`);
      // The result will contain the message ID in the response, but the response is a giant string not an object.
      const match = result.match(/id\s+([a-f0-9-]+)/i);
      const messageId = match ? match[1] : null;
      messageIdObject[title] = messageId;
      // visit the details page to ensure the message is created
      await visit(`/vault/config-ui/messages/${messageId}/details`);
    };

    this.deleteMessages = async () => {
      // Store message IDs in an object to ensure all are deleted, even if a test is interrupted.
      for (const id of Object.values(messageIdObject)) {
        await runCmd(`vault delete sys/config/ui/custom-messages/${id}`);
      }
      await visit(`/vault/config-ui/messages`); // redirect to messages index after delete to ensure the state is refreshed
    };

    this.createMessageBrowser = async ({
      title,
      type = 'banner',
      end_time = '2023-12-12',
      authenticated = true,
    }) => {
      await visit('/vault/config-ui/messages');
      await click(CUSTOM_MESSAGES.tab(authenticated ? 'After user logs in' : 'On login page'));
      await click(GENERAL.button('Create message'));
      await fillIn(CUSTOM_MESSAGES.input('title'), title);

      await click(CUSTOM_MESSAGES.radio(type));
      await fillIn(
        CUSTOM_MESSAGES.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
      );
      await fillIn(
        CUSTOM_MESSAGES.input('start_time'),
        format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
      );
      if (end_time) {
        await click('#specificDate');
        await fillIn(
          CUSTOM_MESSAGES.input('end_time'),
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

  test('authenticated it should create, edit, view, and delete a message', async function (assert) {
    // create first message
    await this.createMessageRepl({ title: 'new-message' });
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText('new-message', 'message title shows on the details screen');
    // edit message
    await click(GENERAL.linkTo('edit'));
    await fillIn(CUSTOM_MESSAGES.input('title'), `Edited new-message`);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText(`Edited new-message`, 'edited message title shows on the details screen');
    await this.deleteMessages();
    const linkedBlocks = findAll('[data-test-list-item]');
    assert.false(linkedBlocks.includes(`Edited new-message`), 'edited message was deleted.');
  });

  test('authenticated it should show multiple messages modal', async function (assert) {
    await this.createMessageRepl({ title: 'message-one', type: 'modal' });
    // create second message with same model name through the UI (not the webrepl)
    await this.createMessageBrowser({ title: 'message-one', type: 'modal', end_time: null });

    assert.dom(CUSTOM_MESSAGES.modal('multiple modal messages')).exists();
    assert
      .dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal'))
      .hasText('Warning: more than one modal after the user logs in');

    await click(GENERAL.button('cancel-multiple')); // cancel out of the modal
    await click(GENERAL.cancelButton); // cancel out of the create message form
    // delete the created message to avoid test pollution
    await this.deleteMessages();
  });

  test('it should filter by type and status', async function (assert) {
    await this.createMessageRepl({ title: 'filter-status-1', type: 'banner' });
    await this.createMessageRepl({
      title: 'filter-status-2',
      type: 'banner',
      end_time: '2023-12-22T00:00:00.000Z',
    });
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
    await this.deleteMessages();
  });

  test('it should display preview a message when all required fields are filled out', async function (assert) {
    await click(GENERAL.navLink('Operational tools'));
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('After user logs in'));
    await click(GENERAL.button('Create message'));
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
    await click(GENERAL.navLink('Operational tools'));
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('After user logs in'));
    await click(GENERAL.button('Create message'));
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
    assert.dom(GENERAL.validationErrorByAttr('title')).exists();
  });

  // unauthenticated messages
  test('unauthenticated it should create, edit, view, and delete a message', async function (assert) {
    await this.createMessageRepl({
      title: 'unauthenticated create edit view delete',
      type: 'banner',
      authenticated: false,
    });
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText('unauthenticated create edit view delete', 'title shows on the details screen');
    // navigate to edit the title
    await click(GENERAL.linkTo('edit'));
    await fillIn(CUSTOM_MESSAGES.input('title'), `Edited ${'unauthenticated create edit view delete'}`);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText(
        `Edited ${'unauthenticated create edit view delete'}`,
        'edited title shows on the details screen'
      );
    // delete the edited message
    await this.deleteMessages();
  });

  test('unauthenticated it should show multiple messages modal', async function (assert) {
    await this.createMessageRepl({
      title: 'unauthenticated message 1',
      type: 'modal',
      authenticated: false,
    });
    // create second message with same model name
    await this.createMessageBrowser({
      title: 'unauthenticated message 1',
      type: 'modal',
      authenticated: false,
    });
    assert.dom(CUSTOM_MESSAGES.modal('multiple modal messages')).exists('the multiple modal message shows');
    assert
      .dom(CUSTOM_MESSAGES.modalTitle('Warning: more than one modal'))
      .hasText('Warning: more than one modal on the login page', 'the warning modal title shows');

    await click(GENERAL.button('cancel-multiple')); // cancel out of the modal
    await click(GENERAL.cancelButton); // cancel out of the create message form
    // delete the created message
    await this.deleteMessages();
  });

  test('it should show info message about sensitive information on create and edit form', async function (assert) {
    await click(GENERAL.navLink('Operational tools'));
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('On login page'));
    await click(GENERAL.button('Create message'));
    assert
      .dom(CUSTOM_MESSAGES.unauthCreateFormInfo)
      .hasText(
        'Note: Do not include sensitive information in this message since users are unauthenticated at this stage.'
      );
  });

  test('it should allow you to preview a message when all required fields are filled out', async function (assert) {
    await click(GENERAL.navLink('Operational tools'));
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('On login page'));
    await click(GENERAL.button('Create message'));
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
    await click(GENERAL.navLink('Operational tools'));
    await click(CUSTOM_MESSAGES.navLink);
    await click(CUSTOM_MESSAGES.tab('On login page'));
    await click(GENERAL.button('Create message'));
    await click(CUSTOM_MESSAGES.radio('banner'));
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');
    await click(GENERAL.button('preview'));
    assert.dom(CUSTOM_MESSAGES.modal('preview image')).doesNotExist('preview image does not show');
    assert.dom(GENERAL.validationErrorByAttr('title')).exists();
  });

  test('cleanup message pollution', async function (assert) {
    // Visit the messages page and delete any remaining messages.
    await visit('/vault/config-ui/messages');
    const rows = findAll('.list-item-row');
    for (const row of rows) {
      const trigger = row.querySelector('[data-test-popup-menu-trigger]');
      if (trigger) {
        await click(GENERAL.menuTrigger);
        await click(GENERAL.menuItem('delete'));
        await click(GENERAL.confirmButton);
      }
    }

    // Redirect to the dashboard and revisit the messages page to refresh the state.
    await visit('/vault/dashboard');
    await visit('/vault/config-ui/messages');

    // Wait for the empty state to render and assert that no messages exist.
    await waitFor(GENERAL.emptyStateTitle, {
      timeout: 2000,
      timeoutMessage: 'Timed out waiting for empty state title to render',
    });
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet', 'No messages exist after cleanup');
  });
});
