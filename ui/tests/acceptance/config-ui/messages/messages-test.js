/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit, fillIn } from '@ember/test-helpers';
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
    this.createMessage = async () => {
      await visit('vault/config-ui/messages');
      await click(MESSAGE_SELECTORS.button('create message'));
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Awesome custom message title');
      await fillIn(
        MESSAGE_SELECTORS.input('message'),
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
      );
      await fillIn(
        MESSAGE_SELECTORS.input('startTime'),
        format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
      );
      await click('#specificDate');
      await fillIn(
        MESSAGE_SELECTORS.input('endTime'),
        format(addDays(startOfDay(new Date('2023-12-12')), 10), datetimeLocalStringFormat)
      );
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
    await visit('vault/config-ui/messages');
    await assert.dom('[data-test-component="empty-state"]').exists();
    await assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
    await click(MESSAGE_SELECTORS.tab('On login page'));
    await assert.dom('[data-test-component="empty-state"]').exists();
    await assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
  });

  module('Authenticated messages', function () {
    test('it should create and edit a message', async function (assert) {
      await this.createMessage();
      assert.dom(GENERAL.title).hasText('Awesome custom message title');
      await click('[data-test-link="edit"]');
      await fillIn(MESSAGE_SELECTORS.input('title'), 'Edited custom message title');
      await click(MESSAGE_SELECTORS.button('create-message'));
      assert.dom(GENERAL.title).hasText('Edited custom message title');
      await click('[data-test-confirm-action="Delete message"]');
      await click(GENERAL.confirmButton);
    });
    test('it should delete a message', async function (assert) {
      this.createMessage();
      await visit('vault/config-ui/messages');
      await assert.dom('[data-test-component="empty-state"]').doesNotExist();
      await click(GENERAL.menuTrigger);
      await click(GENERAL.confirmTrigger);
      await click(GENERAL.confirmButton);
      await assert.dom('[data-test-component="empty-state"]').exists();
    });
  });
});
