/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, visit, waitUntil, find, fillIn } from '@ember/test-helpers';
import { allSupportedAuthBackends, supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { setupMirage } from 'ember-cli-mirage/test-support';

const AUTH_FORM = {
  method: '[data-test-select=auth-method]',
  token: '[data-test-token]',
  login: '[data-test-auth-submit]',
};
const ENT_AUTH_METHODS = ['saml'];

module('Acceptance | auth', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('auth query params', async function (assert) {
    const backends = supportedAuthBackends();
    assert.expect(backends.length + 1);
    await visit('/vault/auth');
    assert.strictEqual(currentURL(), '/vault/auth?with=token');
    for (const backend of backends.reverse()) {
      await fillIn(AUTH_FORM.method, backend.type);
      assert.strictEqual(
        currentURL(),
        `/vault/auth?with=${backend.type}`,
        `has the correct URL for ${backend.type}`
      );
    }
  });

  test('it clears token when changing selected auth method', async function (assert) {
    await visit('/vault/auth');
    await fillIn(AUTH_FORM.token, 'token');
    await fillIn(AUTH_FORM.method, 'github');
    await fillIn(AUTH_FORM.method, 'token');
    assert.dom(AUTH_FORM.token).hasNoValue('it clears the token value when toggling methods');
  });

  module('it sends the right payload when authenticating', function (hooks) {
    hooks.beforeEach(function () {
      this.assertReq = () => {};
      this.server.get('/auth/token/lookup-self', (schema, req) => {
        this.assertReq(req);
        req.passthrough();
      });
      this.server.post('/auth/:mount/login', (schema, req) => {
        // github only
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
      this.server.put('/auth/:mount/sso_service_url', (schema, req) => {
        // SAML only (enterprise)
        this.assertReq(req);
        req.passthrough();
      });
      this.expected = {
        token: {
          url: '/v1/auth/token/lookup-self',
          payload: {
            'X-Vault-Token': 'some-token',
          },
        },
        userpass: {
          url: '/v1/auth/custom-userpass/login/some-username',
          payload: {
            password: 'some-password',
          },
        },
        ldap: {
          url: '/v1/auth/custom-ldap/login/some-username',
          payload: {
            password: 'some-password',
          },
        },
        okta: {
          url: '/v1/auth/custom-okta/login/some-username',
          payload: {
            password: 'some-password',
          },
        },
        jwt: {
          url: '/v1/auth/custom-jwt/oidc/auth_url',
          payload: {
            redirect_uri: 'http://localhost:7357/ui/vault/auth/custom-jwt/oidc/callback',
            role: 'some-role',
          },
        },
        oidc: {
          url: '/v1/auth/custom-oidc/oidc/auth_url',
          payload: {
            redirect_uri: 'http://localhost:7357/ui/vault/auth/custom-oidc/oidc/callback',
            role: 'some-role',
          },
        },
        radius: {
          url: '/v1/auth/custom-radius/login/some-username',
          payload: {
            password: 'some-password',
          },
        },
        github: {
          url: '/v1/auth/custom-github/login',
          payload: {
            token: 'some-token',
          },
        },
        saml: {
          url: '/v1/auth/custom-saml/sso_service_url',
          payload: {
            role: 'some-role',
          },
        },
      };
    });

    for (const backend of allSupportedAuthBackends().reverse()) {
      test(`for ${backend.type} ${
        ENT_AUTH_METHODS.includes(backend.type) ? '(enterprise)' : ''
      }`, async function (assert) {
        const { type } = backend;
        const expected = this.expected[type];
        const isOidc = ['oidc', 'jwt'].includes(type);
        assert.expect(isOidc ? 6 : 2);

        this.assertReq = (req) => {
          const body = type === 'token' ? req.requestHeaders : JSON.parse(req.requestBody);
          if (isOidc && !body.role) {
            // OIDC and JWT auth form calls the endpoint every time the role or mount is updated.
            // if role is not provided, it means we haven't filled out the full info yet so don't
            // validate the payload until all data is provided
            // eslint-disable-next-line qunit/no-early-return
            return {};
          }
          assert.strictEqual(req.url, expected.url, `${type} calls the correct URL`);
          Object.keys(expected.payload).forEach((expKey) => {
            assert.strictEqual(
              body[expKey],
              expected.payload[expKey],
              `${type} payload includes ${expKey} with expected value`
            );
          });
        };
        await visit('/vault/auth');
        await fillIn(AUTH_FORM.method, type);

        if (type !== 'token') {
          // set custom mount
          await click('[data-test-auth-form-options-toggle]');
          await fillIn('[data-test-auth-form-mount-path]', `custom-${type}`);
        }
        backend.formAttributes.forEach(async (key) => {
          // fill in all form items, except JWT which is not rendered
          if (key === 'jwt') return;
          await fillIn(`[data-test-${key}]`, `some-${key}`);
        });

        await click(AUTH_FORM.login);
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
    await fillIn(AUTH_FORM.method, 'token');
    await click('[data-test-auth-submit]');
  });
});
