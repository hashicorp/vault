/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, visit, waitUntil, find } from '@ember/test-helpers';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import authForm from '../pages/components/auth-form';
import jwtForm from '../pages/components/auth-jwt';
import { create } from 'ember-cli-page-object';
import { setupMirage } from 'ember-cli-mirage/test-support';

const component = create(authForm);
const jwtComponent = create(jwtForm);

module('Acceptance | auth', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

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

  module('it sends the right attributes when authenticating', function (hooks) {
    hooks.beforeEach(function () {
      this.assertReq = () => {};
      this.server.get('/auth/token/lookup-self', (schema, req) => {
        this.assertReq(req);
        req.passthrough();
      });
      this.server.post('/auth/github/login', (schema, req) => {
        // This one is for github only
        this.assertReq(req);
        req.passthrough();
      });
      this.server.post('/auth/:mount/oidc/auth_url', (schema, req) => {
        // For JWT and OIDC
        this.assertReq(req);
        req.passthrough();
      });
      this.server.post('/auth/:mount/login/:username', (schema, req) => {
        this.assertReq(req);
        req.passthrough();
      });
      this.expected = {
        token: {
          included: 'X-Vault-Token',
          url: '/v1/auth/token/lookup-self',
        },
        userpass: {
          included: 'password',
          url: '/v1/auth/userpass/login/null',
        },
        ldap: {
          included: 'password',
          url: '/v1/auth/ldap/login/null',
        },
        okta: {
          included: 'password',
          url: '/v1/auth/okta/login/null',
        },
        jwt: {
          included: 'role',
          url: '/v1/auth/jwt/oidc/auth_url',
        },
        oidc: {
          included: 'role',
          url: '/v1/auth/oidc/oidc/auth_url',
        },
        radius: {
          included: 'password',
          url: '/v1/auth/radius/login/null',
        },
        github: {
          included: 'token',
          url: '/v1/auth/github/login',
        },
      };
    });

    for (const backend of supportedAuthBackends().reverse()) {
      test(`for ${backend.type}`, async function (assert) {
        const { type } = backend;
        const isOidc = ['jwt', 'oidc'].includes(type);
        // OIDC types make 3 requests, each time the role changes
        assert.expect(isOidc ? 6 : 2);
        this.assertReq = (req) => {
          const body = type === 'token' ? req.requestHeaders : JSON.parse(req.requestBody);
          const { included, url } = this.expected[type];
          assert.true(Object.keys(body).includes(included), `${type} includes ${included}`);
          assert.strictEqual(req.url, url, `${type} calls the correct URL`);
        };
        await visit('/vault/auth');
        await component.selectMethod(type);
        if (type === 'github') {
          await component.token('token');
        }
        if (isOidc) {
          await jwtComponent.role('test');
        }
        await component.login();
      });
    }
  });

  test('it shows the push notification warning after submit', async function (assert) {
    assert.expect(1);

    this.server.get(
      '/auth/token/lookup-self',
      async () => {
        assert.ok(
          await waitUntil(() => find('[data-test-auth-message="push"]')),
          'shows push notification message'
        );
        return {};
      },
      { timing: 1000 }
    );
    await visit('/vault/auth');
    await component.selectMethod('token');
    await click('[data-test-auth-submit]');
  });
});
