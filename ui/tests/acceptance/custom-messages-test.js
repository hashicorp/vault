/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';

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
});
