/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint qunit/no-conditional-assertions: "warn" */
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import sinon from 'sinon';
import { click, currentURL, visit, waitUntil, find } from '@ember/test-helpers';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import authForm from '../pages/components/auth-form';
import jwtForm from '../pages/components/auth-jwt';
import { create } from 'ember-cli-page-object';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import { PAGE } from 'vault/tests/helpers/config-ui/message-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';

const component = create(authForm);
const jwtComponent = create(jwtForm);

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

module('Acceptance | auth route template', function (hooks) {
  setupApplicationTest(hooks);

  module('auth template', function (hooks) {
    hooks.beforeEach(function () {
      this.clock = sinon.useFakeTimers({
        now: Date.now(),
        shouldAdvanceTime: true,
      });
      this.server = apiStub({ usePassthrough: true });
    });

    hooks.afterEach(function () {
      this.clock.restore();
      this.server.shutdown();
    });

    test('auth query params', async function (assert) {
      const backends = supportedAuthBackends();
      assert.expect(backends.length + 1);
      await visit('/vault/auth');
      assert.strictEqual(currentURL(), '/vault/auth?with=token');
      for (const backend of backends.reverse()) {
        await component.selectMethod(backend.type);
        assert.strictEqual(
          currentURL(),
          `/vault/auth?with=${backend.type}`,
          `has the correct URL for ${backend.type}`
        );
      }
    });

    test('it clears token when changing selected auth method', async function (assert) {
      await visit('/vault/auth');
      await component.token('token').selectMethod('github');
      await component.selectMethod('token');
      assert.strictEqual(component.tokenValue, '', 'it clears the token value when toggling methods');
    });

    test('it sends the right attributes when authenticating', async function (assert) {
      assert.expect(8);
      const backends = supportedAuthBackends();
      await visit('/vault/auth');
      for (const backend of backends.reverse()) {
        await component.selectMethod(backend.type);
        if (backend.type === 'github') {
          await component.token('token');
        }
        if (backend.type === 'jwt' || backend.type === 'oidc') {
          await jwtComponent.role('test');
        }
        await component.login();
        const lastRequest = this.server.passthroughRequests[this.server.passthroughRequests.length - 1];
        let body = JSON.parse(lastRequest.requestBody);
        // Note: x-vault-token used to be lowercase prior to upgrade
        if (backend.type === 'token') {
          assert.ok(
            Object.keys(lastRequest.requestHeaders).includes('X-Vault-Token'),
            'token uses vault token header'
          );
        } else if (backend.type === 'github') {
          assert.ok(Object.keys(body).includes('token'), 'GitHub includes token');
        } else if (backend.type === 'jwt' || backend.type === 'oidc') {
          const authReq = this.server.passthroughRequests[this.server.passthroughRequests.length - 2];
          body = JSON.parse(authReq.requestBody);
          assert.ok(Object.keys(body).includes('role'), `${backend.type} includes role`);
        } else {
          assert.ok(Object.keys(body).includes('password'), `${backend.type} includes password`);
        }
      }
    });

    test('it shows the push notification warning after submit', async function (assert) {
      assert.expect(1);

      this.server.get('/v1/auth/token/lookup-self', async () => {
        assert.ok(
          await waitUntil(() => find('[data-test-auth-message="push"]')),
          'shows push notification message'
        );
        return [204, { 'Content-Type': 'application/json' }, JSON.stringify({})];
      });

      await visit('/vault/auth');
      await component.selectMethod('token');
      await click('[data-test-auth-submit]');
    });
  });

  module('custom messages auth tests', function (hooks) {
    setupMirage(hooks);
    test('it shows the alert banner and modal message', async function (assert) {
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
});
