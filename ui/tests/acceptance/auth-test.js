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

const component = create(authForm);
const jwtComponent = create(jwtForm);

module('Acceptance | auth', function (hooks) {
  setupApplicationTest(hooks);

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
      const body = JSON.parse(lastRequest.requestBody);

      let keys;
      let included;
      if (backend.type === 'token') {
        keys = lastRequest.requestHeaders;
        included = 'x-vault-token';
      } else if (backend.type === 'github') {
        keys = body;
        included = 'token';
      } else if (backend.type === 'jwt' || backend.type === 'oidc') {
        const authReq = this.server.passthroughRequests[this.server.passthroughRequests.length - 2];
        keys = JSON.parse(authReq.requestBody);
        included = 'role';
      } else {
        keys = body;
        included = 'password';
      }
      assert.ok(Object.keys(keys).includes(included), `${backend.type} includes ${included}`);
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
