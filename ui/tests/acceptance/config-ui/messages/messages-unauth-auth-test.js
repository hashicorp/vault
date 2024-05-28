/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, fillIn, currentRouteName } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';

import authPage from 'vault/tests/pages/auth';
import { datetimeLocalStringFormat } from 'core/utils/date-formatters';
import { format, addDays, startOfDay } from 'date-fns';
import { createNS, runCmd } from '../../../helpers/commands';

const unauthenticatedMessageResponse = {
  request_id: '664fbad0-fcd8-9023-4c5b-81a7962e9f4b',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    key_info: {
      'some-awesome-id-2': {
        authenticated: false,
        end_time: null,
        link: {
          'some alert link': 'world',
        },
        message: 'aGVsbG8gd29ybGQgaGVsbG8gd29scmQ=',
        options: null,
        start_time: '2024-01-04T08:00:00Z',
        title: 'Banner title',
        type: 'banner',
      },
      'some-awesome-id-1': {
        authenticated: false,
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

// custom messages is an enterprise feature, but since we're using mirage enterprise is omitted from the module name.
module('Acceptance | auth custom messages auth tests', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  module('auth and unauth messages', function (hooks) {
    hooks.beforeEach(function () {
      return this.server.get('/sys/internal/ui/mounts', () => ({}));
    });

    test('it shows the alert banner and modal message', async function (assert) {
      this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
        return unauthenticatedMessageResponse;
      });
      await visit('/vault/auth');
      const modalId = 'some-awesome-id-1';
      const alertId = 'some-awesome-id-2';
      assert.dom(CUSTOM_MESSAGES.modal(modalId)).exists();
      assert.dom(CUSTOM_MESSAGES.modalTitle(modalId)).hasText('Modal title');
      assert.dom(CUSTOM_MESSAGES.modalBody(modalId)).exists();
      assert.dom(CUSTOM_MESSAGES.modalBody(modalId)).hasText('here is a cool message');
      await click(CUSTOM_MESSAGES.modalButton(modalId));
      assert.dom(CUSTOM_MESSAGES.alertTitle(alertId)).hasText('Banner title');
      assert.dom(CUSTOM_MESSAGES.alertDescription(alertId)).hasText('hello world hello wolrd');
      assert.dom(CUSTOM_MESSAGES.alertAction('link')).hasText('some alert link');
    });
    test('it shows the multiple modal messages', async function (assert) {
      const modalIdOne = 'some-awesome-id-2';
      const modalIdTwo = 'some-awesome-id-1';

      this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
        unauthenticatedMessageResponse.data.key_info[modalIdOne].type = 'modal';
        unauthenticatedMessageResponse.data.key_info[modalIdOne].title = 'Modal title 1';
        unauthenticatedMessageResponse.data.key_info[modalIdTwo].type = 'modal';
        unauthenticatedMessageResponse.data.key_info[modalIdTwo].title = 'Modal title 2';
        return unauthenticatedMessageResponse;
      });
      await visit('/vault/auth');
      assert.dom(CUSTOM_MESSAGES.modal(modalIdOne)).exists();
      assert.dom(CUSTOM_MESSAGES.modalTitle(modalIdOne)).hasText('Modal title 1');
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdOne)).exists();
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdOne)).hasText('hello world hello wolrd some alert link');
      await click(CUSTOM_MESSAGES.modalButton(modalIdOne));
      assert.dom(CUSTOM_MESSAGES.modal(modalIdTwo)).exists();
      assert.dom(CUSTOM_MESSAGES.modalTitle(modalIdTwo)).hasText('Modal title 2');
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdTwo)).exists();
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdTwo)).hasText('here is a cool message');
      await click(CUSTOM_MESSAGES.modalButton(modalIdTwo));
    });
    test('it shows the multiple banner messages', async function (assert) {
      const bannerIdOne = 'some-awesome-id-2';
      const bannerIdTwo = 'some-awesome-id-1';

      this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
        unauthenticatedMessageResponse.data.key_info[bannerIdOne].type = 'banner';
        unauthenticatedMessageResponse.data.key_info[bannerIdOne].title = 'Banner title 1';
        unauthenticatedMessageResponse.data.key_info[bannerIdTwo].type = 'banner';
        unauthenticatedMessageResponse.data.key_info[bannerIdTwo].title = 'Banner title 2';
        return unauthenticatedMessageResponse;
      });
      await visit('/vault/auth');
      assert.dom(CUSTOM_MESSAGES.alertTitle(bannerIdOne)).hasText('Banner title 1');
      assert.dom(CUSTOM_MESSAGES.alertDescription(bannerIdOne)).hasText('hello world hello wolrd');
      assert.dom(CUSTOM_MESSAGES.alertTitle(bannerIdTwo)).hasText('Banner title 2');
      assert.dom(CUSTOM_MESSAGES.alertDescription(bannerIdTwo)).hasText('here is a cool message');
      assert.dom(CUSTOM_MESSAGES.alertAction('link')).hasText('some alert link');
    });
  });

  test('it should display an active authenticated message after creation on enterprise', async function (assert) {
    assert.expect(4);
    await authPage.login();
    await visit('vault/config-ui/messages');
    await click(CUSTOM_MESSAGES.button('create message'));
    await fillIn(CUSTOM_MESSAGES.input('title'), 'Awesome custom message title');
    await click(CUSTOM_MESSAGES.radio('banner'));
    await fillIn(
      CUSTOM_MESSAGES.input('message'),
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Pulvinar mattis nunc sed blandit libero volutpat sed cras ornare.'
    );
    await fillIn(
      CUSTOM_MESSAGES.input('startTime'),
      format(addDays(startOfDay(new Date('2023-12-12')), 1), datetimeLocalStringFormat)
    );
    await fillIn('[data-test-kv-key="0"]', 'Learn more');
    await fillIn('[data-test-kv-value="0"]', 'www.learn.com');

    await click(CUSTOM_MESSAGES.button('create-message'));
    assert.dom(GENERAL.title).hasText('Awesome custom message title', 'on the details screen');
    assert.dom('.hds-alert').exists('active custom message displays on authenticated.');
    await runCmd(createNS('world'), false);
    await visit('vault/config-ui/messages');
    assert.dom('.hds-alert').exists('active custom message displays on namespace authenticated.');
    await click(CUSTOM_MESSAGES.listItem('Awesome custom message title'));
    await click(CUSTOM_MESSAGES.confirmActionButton('Delete message'));
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.config-ui.messages.index',
      'redirects to messages page after delete'
    );
  });
});
