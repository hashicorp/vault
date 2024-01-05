/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/config-ui/message-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';

const unauthenticatedMessageResponse = {
  request_id: '664fbad0-fcd8-9023-4c5b-81a7962e9f4b',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    key_info: {
      '02180e3f-bd5b-a851-bcc9-6f7983806df0': {
        authenticated: false,
        end_time: null,
        link: {
          title: '',
        },
        message: 'aGVsbG8gd29ybGQgaGVsbG8gd29scmQ=',
        options: null,
        start_time: '2024-01-04T08:00:00Z',
        title: 'Banner title',
        type: 'banner',
      },
      'a7d7d9b1-a1ca-800c-17c5-0783be88e29c': {
        authenticated: false,
        end_time: null,
        link: {
          title: '',
        },
        message: 'aGVyZSBpcyBhIGNvb2wgbWVzc2FnZQ==',
        options: null,
        start_time: '2024-01-01T08:00:00Z',
        title: 'Modal title',
        type: 'modal',
      },
    },
    keys: ['02180e3f-bd5b-a851-bcc9-6f7983806df0', 'a7d7d9b1-a1ca-800c-17c5-0783be88e29c'],
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: '',
};

module('Acceptance | auth custom messages auth tests', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('it shows the alert banner and modal message', async function (assert) {
    this.server.get('/sys/internal/ui/mounts', () => ({}));
    this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
      return unauthenticatedMessageResponse;
    });
    await visit('/vault/auth');
    const modalId = 'a7d7d9b1-a1ca-800c-17c5-0783be88e29c';
    const alertId = '02180e3f-bd5b-a851-bcc9-6f7983806df0';
    assert.dom(PAGE.modal(modalId)).exists();
    assert.dom(PAGE.modalTitle(modalId)).hasText('Modal title');
    assert.dom(PAGE.modalBody(modalId)).exists();
    assert.dom(PAGE.modalBody(modalId)).hasText('here is a cool message');
    await click(PAGE.modalButton(modalId));
    assert.dom(PAGE.alertTitle(alertId)).hasText('Banner title');
    assert.dom(PAGE.alertDescription(alertId)).hasText('hello world hello wolrd');
  });
  test('it shows the multiple modal messages', async function (assert) {
    const modalIdOne = '02180e3f-bd5b-a851-bcc9-6f7983806df0';
    const modalIdTwo = 'a7d7d9b1-a1ca-800c-17c5-0783be88e29c';
    this.server.get('/sys/internal/ui/mounts', () => ({}));

    this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
      unauthenticatedMessageResponse.data.key_info[modalIdOne].type = 'modal';
      unauthenticatedMessageResponse.data.key_info[modalIdOne].title = 'Modal title 1';
      unauthenticatedMessageResponse.data.key_info[modalIdTwo].type = 'modal';
      unauthenticatedMessageResponse.data.key_info[modalIdTwo].title = 'Modal title 2';
      return unauthenticatedMessageResponse;
    });
    await visit('/vault/auth');
    assert.dom(PAGE.modal(modalIdOne)).exists();
    assert.dom(PAGE.modalTitle(modalIdOne)).hasText('Modal title 1');
    assert.dom(PAGE.modalBody(modalIdOne)).exists();
    assert.dom(PAGE.modalBody(modalIdOne)).hasText('hello world hello wolrd');
    await click(PAGE.modalButton(modalIdOne));
    assert.dom(PAGE.modal(modalIdTwo)).exists();
    assert.dom(PAGE.modalTitle(modalIdTwo)).hasText('Modal title 2');
    assert.dom(PAGE.modalBody(modalIdTwo)).exists();
    assert.dom(PAGE.modalBody(modalIdTwo)).hasText('here is a cool message');
    await click(PAGE.modalButton(modalIdTwo));
  });
  test('it shows the multiple banner messages', async function (assert) {
    const bannerIdOne = '02180e3f-bd5b-a851-bcc9-6f7983806df0';
    const bannerIdTwo = 'a7d7d9b1-a1ca-800c-17c5-0783be88e29c';
    this.server.get('/sys/internal/ui/mounts', () => ({}));

    this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
      unauthenticatedMessageResponse.data.key_info[bannerIdOne].type = 'banner';
      unauthenticatedMessageResponse.data.key_info[bannerIdOne].title = 'Banner title 1';
      unauthenticatedMessageResponse.data.key_info[bannerIdTwo].type = 'banner';
      unauthenticatedMessageResponse.data.key_info[bannerIdTwo].title = 'Banner title 2';
      return unauthenticatedMessageResponse;
    });
    await visit('/vault/auth');
    assert.dom(PAGE.alertTitle(bannerIdOne)).hasText('Banner title 1');
    assert.dom(PAGE.alertDescription(bannerIdOne)).hasText('hello world hello wolrd');
    assert.dom(PAGE.alertTitle(bannerIdTwo)).hasText('Banner title 2');
    assert.dom(PAGE.alertDescription(bannerIdTwo)).hasText('here is a cool message');
  });
});
