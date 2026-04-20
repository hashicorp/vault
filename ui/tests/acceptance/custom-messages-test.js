/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, fillIn, click, select } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { login, logout } from 'vault/tests/helpers/auth/auth-helpers';

// Helper function to generate mock messages
function generateMessages(authenticated, count = 12) {
  const messages = {
    request_id: `${authenticated ? 'auth' : 'unauth'}-request-id`,
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      key_info: {},
      keys: [],
    },
    wrap_info: null,
    warnings: null,
    auth: null,
    mount_type: '',
  };

  for (let i = 0; i < count; i++) {
    const id = `${authenticated ? 'auth' : 'unauth'}-message-${i}`;
    messages.data.keys.push(id);
    messages.data.key_info[id] = {
      authenticated,
      active: i < 3,
      end_time: null,
      message: btoa(`${authenticated ? 'Authenticated' : 'Unauthenticated'} message content ${i}`),
      options: null,
      start_time: '2024-01-01T08:00:00Z',
      title: `${authenticated ? 'Auth' : 'Unauth'} Message ${i}`,
      type: i % 2 === 0 ? 'banner' : 'modal',
    };
  }

  return messages;
}

const authenticatedMessageResponse = {
  request_id: '664fbad0-fcd8-9023-4c5b-81a7962e9f4b',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    key_info: {
      'some-awesome-id-2': {
        authenticated: true,
        end_time: null,
        link: {
          'some link title': 'www.link.com',
        },
        message: 'aGVsbG8gd29ybGQgaGVsbG8gd29scmQ=',
        options: null,
        start_time: '2024-01-04T08:00:00Z',
        title: 'Banner title',
        type: 'banner',
      },
      'some-awesome-id-1': {
        authenticated: true,
        end_time: null,
        message: 'aGVyZSBpcyBhIGNvb2wgbWVzc2FnZQ==',
        options: null,
        start_time: '2024-01-01T08:00:00Z',
        title: 'Modal title',
        type: 'modal',
      },
    },
    keys: ['some-awesome-id-2', 'some-awesome-id-1'],
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: '',
};

module('Acceptance | custom messages', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('it shows the alert banner and modal message', async function (assert) {
    this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
      return authenticatedMessageResponse;
    });
    await visit('/vault/dashboard');

    const modalId = 'some-awesome-id-1';
    const alertId = 'some-awesome-id-2';
    assert.dom(CUSTOM_MESSAGES.modal(modalId)).exists();
    assert.dom(CUSTOM_MESSAGES.modalTitle(modalId)).hasText('Modal title');
    assert.dom(CUSTOM_MESSAGES.modalBody(modalId)).exists();
    assert.dom(CUSTOM_MESSAGES.modalBody(modalId)).hasText('here is a cool message');
    assert.dom(CUSTOM_MESSAGES.alertTitle(alertId)).hasText('Banner title');
    assert.dom(CUSTOM_MESSAGES.alertDescription(alertId)).hasText('hello world hello wolrd');
    assert.dom(CUSTOM_MESSAGES.alertAction('link')).hasText('some link title');
  });

  test('it shows the multiple modal messages', async function (assert) {
    const modalIdOne = 'some-awesome-id-2';
    const modalIdTwo = 'some-awesome-id-1';

    this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
      authenticatedMessageResponse.data.key_info[modalIdOne].type = 'modal';
      authenticatedMessageResponse.data.key_info[modalIdOne].title = 'Modal title 1';
      authenticatedMessageResponse.data.key_info[modalIdTwo].type = 'modal';
      authenticatedMessageResponse.data.key_info[modalIdTwo].title = 'Modal title 2';
      return authenticatedMessageResponse;
    });
    await visit('/vault/dashboard');

    assert.dom(CUSTOM_MESSAGES.modal(modalIdOne)).exists();
    assert.dom(CUSTOM_MESSAGES.modalTitle(modalIdOne)).hasText('Modal title 1');
    assert.dom(CUSTOM_MESSAGES.modalBody(modalIdOne)).exists();
    assert.dom(CUSTOM_MESSAGES.modalBody(modalIdOne)).hasText('hello world hello wolrd some link title');
    assert.dom(CUSTOM_MESSAGES.modal(modalIdTwo)).exists();
    assert.dom(CUSTOM_MESSAGES.modalTitle(modalIdTwo)).hasText('Modal title 2');
    assert.dom(CUSTOM_MESSAGES.modalBody(modalIdTwo)).exists();
    assert.dom(CUSTOM_MESSAGES.modalBody(modalIdTwo)).hasText('here is a cool message');
  });

  test('it shows the multiple banner messages', async function (assert) {
    const bannerIdOne = 'some-awesome-id-2';
    const bannerIdTwo = 'some-awesome-id-1';

    this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
      authenticatedMessageResponse.data.key_info[bannerIdOne].type = 'banner';
      authenticatedMessageResponse.data.key_info[bannerIdOne].title = 'Banner title 1';
      authenticatedMessageResponse.data.key_info[bannerIdTwo].type = 'banner';
      authenticatedMessageResponse.data.key_info[bannerIdTwo].title = 'Banner title 2';
      return authenticatedMessageResponse;
    });
    await visit('/vault/dashboard');

    assert.dom(CUSTOM_MESSAGES.alertTitle(bannerIdOne)).hasText('Banner title 1');
    assert.dom(CUSTOM_MESSAGES.alertDescription(bannerIdOne)).hasText('hello world hello wolrd');
    assert.dom(CUSTOM_MESSAGES.alertAction('link')).hasText('some link title');
    assert.dom(CUSTOM_MESSAGES.alertTitle(bannerIdTwo)).hasText('Banner title 2');
    assert.dom(CUSTOM_MESSAGES.alertDescription(bannerIdTwo)).hasText('here is a cool message');
  });

  test('it should filter by message title and paginate', async function (assert) {
    assert.expect(7);

    await login();

    // Mock authenticated messages endpoint with 12+ messages
    this.server.get('/sys/config/ui/custom-messages', function () {
      return generateMessages(true, 12);
    });

    // Visit the custom messages page
    await visit('/vault/config-ui/messages');

    // Verify pagination is present (more than 10 messages)
    assert.dom(GENERAL.pagination).exists('Pagination exists for messages list');
    // Test filtering by title
    await fillIn(GENERAL.filter('pageFilter'), 'Auth Message 1');
    await click(GENERAL.submitButton);

    // After filtering, we should see only messages matching "Auth Message 1"
    // This would include "Auth Message 1", "Auth Message 10", "Auth Message 11"
    assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'Filter shows matching messages');

    // Clear filter
    await click(GENERAL.button('reset'));

    // Verify all messages are shown again (first page of 10)
    assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'All messages shown after clearing filter');

    // Test filtering with no results
    await fillIn(GENERAL.filter('pageFilter'), 'NonExistentMessage');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.emptyStateTitle).exists('No messages yet');

    // Clear filter again
    await click(GENERAL.button('reset'));

    // Test pagination after filtering
    await fillIn(GENERAL.filter('pageFilter'), 'Message 0');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.listItem()).exists('Filtered results are displayed');

    // Clear and navigate to next page
    await click(GENERAL.button('reset'));
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.listItem()).exists('Second page of messages loads');

    // Test filter message status
    await select(GENERAL.filter('status'), 'active');

    await click(GENERAL.submitButton);
    assert.dom(GENERAL.listItem()).exists('Filtered results are displayed');

    await logout();
  });

  test('it should filter by message status and type', async function (assert) {
    // Mock the list endpoint to respond based on query params
    this.server.get('/sys/config/ui/custom-messages', function (schema, request) {
      const { active, authenticated, type } = request.queryParams;
      const isAuthenticated = authenticated === 'true';
      const allMessages = generateMessages(isAuthenticated, 12);

      // Filter by active status if specified
      let filteredKeys = allMessages.data.keys;
      if (active !== undefined) {
        const isActive = active === 'true';
        filteredKeys = filteredKeys.filter((key) => {
          return allMessages.data.key_info[key].active === isActive;
        });
      }

      // Filter by type if specified
      if (type) {
        filteredKeys = filteredKeys.filter((key) => {
          return allMessages.data.key_info[key].type === type;
        });
      }

      // Build filtered response
      allMessages.data.keys = filteredKeys;
      const filteredKeyInfo = {};
      filteredKeys.forEach((key) => {
        filteredKeyInfo[key] = allMessages.data.key_info[key];
      });
      allMessages.data.key_info = filteredKeyInfo;

      return allMessages;
    });

    await login();

    // Visit the custom messages page
    await visit('/vault/config-ui/messages');

    // Verify initial state - should show all messages (first page of 10)
    assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'Initial page shows 10 messages');

    // Test filtering by status (active)
    await fillIn(GENERAL.filter('status'), 'active');
    await click(GENERAL.submitButton);

    // Verify URL contains status parameter
    assert.ok(
      this.owner.lookup('service:router').currentURL.includes('status=active'),
      'URL includes status=active parameter'
    );
    // With hardcoded active messages (0-7), we should have 8 active messages
    assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'Shows 3 active messages when filtered by status');

    // Test filtering by type (banner)
    await fillIn(GENERAL.filter('type'), 'banner');
    await click(GENERAL.submitButton);

    // Verify URL contains type parameter
    assert.ok(
      this.owner.lookup('service:router').currentURL.includes('type=banner'),
      'URL includes type=banner parameter'
    );

    // Verify filtered results show only active banner messages
    // Active messages: 0-7, Banners (even indices): 0, 2, 4, 6 = 4 active banners
    assert
      .dom(GENERAL.listItem())
      .exists({ count: 2 }, 'Shows 2 active banner messages when filtered by status and type');
    assert.dom(CUSTOM_MESSAGES.badge('banner')).exists('Banner badge is shown');

    // Clear status filter and test type filter alone
    await fillIn(GENERAL.filter('status'), '');
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.listItem()).exists({ count: 6 }, 'Shows 6 banner messages when filtered by type only');

    // Test filtering by both status (active) and type (modal)
    await fillIn(GENERAL.filter('status'), 'active');
    await fillIn(GENERAL.filter('type'), 'modal');
    await click(GENERAL.submitButton);

    // Verify URL contains both parameters
    assert.ok(
      this.owner.lookup('service:router').currentURL.includes('status=active'),
      'URL includes status parameter'
    );
    assert.ok(
      this.owner.lookup('service:router').currentURL.includes('type=modal'),
      'URL includes type parameter'
    );

    assert
      .dom(GENERAL.listItem())
      .exists({ count: 1 }, 'Shows 1 active modal messages when filtered by status and type');
    assert.dom(CUSTOM_MESSAGES.badge('modal')).exists('Modal badge is shown');

    // Test reset filters button
    await click(GENERAL.button('reset'));

    // Verify filters are cleared
    assert.notOk(
      this.owner.lookup('service:router').currentURL.includes('status='),
      'Status filter is cleared from URL'
    );
    assert.notOk(
      this.owner.lookup('service:router').currentURL.includes('type='),
      'Type filter is cleared from URL'
    );

    // Verify all messages are shown again
    assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'All messages shown after clearing filters');

    await logout();
  });
});
