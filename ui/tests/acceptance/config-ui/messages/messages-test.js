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
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | Community | config-ui/messages', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.get('/sys/health', function () {
      return {
        enterprise: false,
        initialized: true,
        sealed: false,
        standby: false,
        license: {
          expiry: '2024-01-12T23:20:50.52Z',
          state: 'stored',
        },
        performance_standby: false,
        replication_performance_mode: 'disabled',
        replication_dr_mode: 'disabled',
        server_time_utc: 1622562585,
        version: '1.16.0',
        cluster_name: 'vault-cluster-e779cd7c',
        cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
        last_wal: 121,
      };
    });
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
    this.server.get('/sys/health', function () {
      return {
        enterprise: true,
        initialized: true,
        sealed: false,
        standby: false,
        license: {
          expiry: '2024-01-12T23:20:50.52Z',
          state: 'stored',
        },
        performance_standby: false,
        replication_performance_mode: 'disabled',
        replication_dr_mode: 'disabled',
        server_time_utc: 1622562585,
        version: '1.16.0+ent',
        cluster_name: 'vault-cluster-e779cd7c',
        cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
        last_wal: 121,
      };
    });
    await authPage.login();
  });

  hooks.afterEach(async function () {
    await logout.visit();
  });
  test('it should show an empty state when no messages are created', async function (assert) {
    assert.expect(4);
    await click(CUSTOM_MESSAGES.navLink);
    assert.dom('[data-test-component="empty-state"]').exists();
    assert.dom(GENERAL.emptyStateTitle).hasText('No messages yet');
    await click(CUSTOM_MESSAGES.tab('On login page'));
    assert.dom('[data-test-component="empty-state"]').exists();
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
      assert.dom('[data-test-component="empty-state"]').exists('Message was deleted');
    });
    test('it should filter by type and status', async function (assert) {
      await this.createMessage('banner', null);
      const msg1 = this.messageDetailId();
      await this.createMessage('banner');
      const msg2 = this.messageDetailId();
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
      assert.dom('[data-test-component="empty-state"]').exists('Message was deleted');
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
