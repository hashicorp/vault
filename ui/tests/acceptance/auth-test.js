/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, visit, waitUntil, find, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allSupportedAuthBackends, supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import {
  createNS,
  createPolicyCmd,
  mountAuthCmd,
  mountEngineCmd,
  runCmd,
} from 'vault/tests/helpers/commands';
import { login, loginMethod, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { v4 as uuidv4 } from 'uuid';
import { GENERAL } from '../helpers/general-selectors';

const ENT_AUTH_METHODS = ['saml'];
const { rootToken } = VAULT_KEYS;

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
    await fillIn(AUTH_FORM.input('token'), 'token');
    await fillIn(AUTH_FORM.method, 'github');
    await fillIn(AUTH_FORM.method, 'token');
    assert.dom(AUTH_FORM.input('token')).hasNoValue('it clears the token value when toggling methods');
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

  test('it does not call renew-self after successful login with non-renewable token', async function (assert) {
    this.server.post(
      '/auth/token/renew-self',
      () => new Error('should not call renew-self directly after logging in')
    );

    await login(rootToken);
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  module('Enterprise', function (hooks) {
    hooks.beforeEach(async function () {
      const uid = uuidv4();
      this.ns = `admin-${uid}`;
      // log in to root to create namespace
      await login();
      await runCmd(createNS(this.ns), false);
      // login to namespace, mount userpass, create policy and user
      await loginNs(this.ns);
      this.db = `database-${uid}`;
      this.userpass = `userpass-${uid}`;
      this.user = 'bob';
      this.policyName = `policy-${this.userpass}`;
      this.policy = `
        path "${this.db}/" {
          capabilities = ["list"]
        }
        path "${this.db}/roles" {
          capabilities = ["read","list"]
        }
        `;
      await runCmd([
        mountAuthCmd('userpass', this.userpass),
        mountEngineCmd('database', this.db),
        createPolicyCmd(this.policyName, this.policy),
        `write auth/${this.userpass}/users/${this.user} password=${this.user} token_policies=${this.policyName}`,
      ]);
      return await logout();
    });

    hooks.afterEach(async function () {
      await visit(`/vault/logout?namespace=${this.ns}`);
      await fillIn(AUTH_FORM.namespaceInput, ''); // clear login form namespace input
      await login();
      await runCmd([`delete sys/namespaces/${this.ns}`], false);
    });

    // this test is specifically to cover a token renewal bug within namespaces
    // namespace_path isn't returned by the renew-self response and so the auth service was
    // incorrectly setting userRootNamespace to '' (which denotes 'root')
    // making subsequent capability checks fail because they would not be queried with the appropriate namespace header
    // if this test fails because a POST /v1/sys/capabilities-self returns a 403, then we have a problem!
    test('it sets namespace when renewing token', async function (assert) {
      await login();
      await runCmd([
        mountAuthCmd('userpass', this.userpass),
        mountEngineCmd('database', this.db),
        createPolicyCmd(this.policyName, this.policy),
        `write auth/${this.userpass}/users/${this.user} password=${this.user} token_policies=${this.policyName}`,
      ]);

      const options = { username: this.user, password: this.user, 'auth-form-mount-path': this.userpass };

      // login as user just to get token (this is the only way to generate a token in the UI right now..)
      await loginMethod('userpass', options, { toggleOptions: true, ns: this.ns });
      await click('[data-test-user-menu-trigger=""]');
      const token = find('[data-test-copy-button]').getAttribute('data-test-copy-button');

      // login with token to reproduce bug
      await loginNs(this.ns, token);
      await visit(`/vault/secrets/${this.db}/overview?namespace=${this.ns}`);
      assert
        .dom('[data-test-overview-card="Roles"]')
        .hasText('Roles Create new', 'database overview renders');
      // renew token
      await click('[data-test-user-menu-trigger=""]');
      await click('[data-test-user-menu-item="renew token"]');
      // navigate out and back to overview tab to re-request capabilities
      await click(GENERAL.secretTab('Roles'));
      await click(GENERAL.tab('overview'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.db}/overview?namespace=${this.ns}`,
        'it navigates to database overview'
      );
    });
  });
});
