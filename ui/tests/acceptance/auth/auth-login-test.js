/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, fillIn, typeIn, visit, waitFor } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allSupportedAuthBackends, supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import {
  createNS,
  createPolicyCmd,
  deleteNS,
  mountAuthCmd,
  mountEngineCmd,
  runCmd,
} from 'vault/tests/helpers/commands';
import {
  login,
  loginMethod,
  loginNs,
  logout,
  SYS_INTERNAL_UI_MOUNTS,
} from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { v4 as uuidv4 } from 'uuid';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const ENT_AUTH_METHODS = ['saml'];
const { rootToken } = VAULT_KEYS;

module('Acceptance | auth login form', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('it does not request login settings for community versions', async function (assert) {
    assert.expect(1); // should only be one assertion because the stubbed mirage request should NOT be hit
    this.owner.lookup('service:version').type = 'community';
    this.server.get('/sys/internal/ui/default-auth-methods', () => {
      // cannot throw error here because request errors are swallowed
      assert.false(true, 'request made for login settings and it should not have been');
    });
    await visit('/vault/auth');
    assert.strictEqual(currentURL(), '/vault/auth');
  });

  test('it selects auth method if "with" query param is a supported auth method', async function (assert) {
    const backends = supportedAuthBackends();
    assert.expect(backends.length);
    for (const backend of backends.reverse()) {
      await visit(`/vault/auth?with=${backend.type}`);
      assert.dom(AUTH_FORM.selectMethod).hasValue(backend.type);
    }
  });

  test('it selects auth method if "with" query param ends in an unencoded a slash', async function (assert) {
    await visit('/vault/auth?with=userpass/');
    assert.dom(AUTH_FORM.selectMethod).hasValue('userpass');
  });

  test('it selects auth method if "with" query param ends in an encoded slash and matches an auth type', async function (assert) {
    await visit('/vault/auth?with=userpass%2F');
    assert.dom(AUTH_FORM.selectMethod).hasValue('userpass');
  });

  test('it redirects if "with" query param is not a supported auth method', async function (assert) {
    await visit('/vault/auth?with=fake');
    assert.strictEqual(currentURL(), '/vault/auth', 'invalid query param is cleared');
  });

  test('it does not refire route model if query param does not exist', async function (assert) {
    const route = this.owner.lookup('route:vault/cluster/auth');
    const modelSpy = sinon.spy(route, 'model');
    await visit('/vault/auth');
    assert.strictEqual(modelSpy.callCount, 1, 'model hook is only called once');
    modelSpy.restore();
  });

  test('it clears token when changing selected auth method', async function (assert) {
    await visit('/vault/auth');
    await fillIn(AUTH_FORM.selectMethod, 'token');
    await fillIn(GENERAL.inputByAttr('token'), 'token');
    await fillIn(AUTH_FORM.selectMethod, 'github');
    await fillIn(AUTH_FORM.selectMethod, 'token');
    assert.dom(GENERAL.inputByAttr('token')).hasNoValue('it clears the token value when toggling methods');
  });

  test('it does not render tabs if sys/internal/ui/mounts is empty', async function (assert) {
    await logout(); // clear local storage
    await visit('/vault/auth');
    await waitFor(AUTH_FORM.form);
    assert.dom(GENERAL.selectByAttr('auth type')).exists('dropdown renders');
    // dropdown could still render in "Sign in with other methods" view, so make sure we're not in a weird state
    assert.dom(GENERAL.backButton).doesNotExist('it does not render "Back" button');
    assert.dom(AUTH_FORM.authForm('token')).exists('it renders token form');
    assert.dom(AUTH_FORM.tabs).doesNotExist();
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(async function () {
      this.server.get('/sys/internal/ui/mounts', () => {
        return { data: { auth: SYS_INTERNAL_UI_MOUNTS } };
      });
      await logout(); // clear local storage
    });

    test('it renders tabs if sys/internal/ui/mounts returns data', async function (assert) {
      assert.expect(9);
      const expectedTabs = [
        { type: 'userpass', display: 'Userpass' },
        { type: 'oidc', display: 'OIDC' },
        { type: 'ldap', display: 'LDAP' },
      ];
      await visit('/vault/auth');
      await waitFor(AUTH_FORM.tabs);
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
      // there are 4 mount paths returned in the stubbed sys/internal/ui/mounts response above,
      // but two are of the same type so only expect 3 tabs
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'it groups mount paths by type and renders 3 tabs');
      expectedTabs.forEach((m) => {
        assert.dom(AUTH_FORM.tabBtn(m.type)).exists(`${m.type} renders as a tab`);
        assert.dom(AUTH_FORM.tabBtn(m.type)).hasText(m.display, `${m.type} renders expected display name`);
      });
      assert
        .dom(AUTH_FORM.tabBtn('userpass'))
        .hasAttribute('aria-selected', 'true', 'it selects the first type by default');
    });

    test('it renders preferred mount view if "with" query param is a mount path with listing_visibility="unauth"', async function (assert) {
      await visit('/vault/auth?with=my_oidc%2F');
      await waitFor(AUTH_FORM.tabBtn('oidc'));
      assert.dom(AUTH_FORM.authForm('oidc')).exists();
      assert.dom(AUTH_FORM.tabBtn('oidc')).exists();
      assert.dom(GENERAL.inputByAttr('role')).exists();
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert
        .dom(GENERAL.inputByAttr('path'))
        .hasValue('my_oidc/', 'mount path matches server value and is not camelized');
      assert.dom(GENERAL.button('Sign in with other methods')).exists('"Sign in with other methods" renders');

      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
      assert.dom(GENERAL.backButton).doesNotExist();
    });

    test('it selects tab if "with" query param matches a tab type', async function (assert) {
      await visit('/vault/auth?with=oidc');
      await waitFor(AUTH_FORM.tabBtn('oidc'));
      assert
        .dom(AUTH_FORM.tabBtn('oidc'))
        .hasAttribute('aria-selected', 'true', 'it selects tab matching query param');
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my_oidc/');
      assert.dom(GENERAL.button('Sign in with other methods')).exists('"Sign in with other methods" renders');
      assert.dom(GENERAL.backButton).doesNotExist();
    });

    test('it selects type from dropdown if query param is NOT a visible mount, but is a supported method', async function (assert) {
      await visit('/vault/auth?with=token');
      await waitFor(GENERAL.selectByAttr('auth type'));
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('token');
      assert.dom(GENERAL.backButton).exists('it renders "Back" button because tabs do exist');
      assert
        .dom(GENERAL.button('Sign in with other methods'))
        .doesNotExist(
          'Tabs exist but query param does not match so login is showing "other" methods and this button should not render'
        );
    });
  });

  module('it sends the right payload when authenticating', function (hooks) {
    hooks.beforeEach(function () {
      this.assertReq = () => {};
      this.server.get('/auth/token/lookup-self', (schema, req) => {
        this.assertReq(req);
        return req.passthrough();
      });
      this.server.post('/auth/:mount/login', (schema, req) => {
        // github only
        this.assertReq(req);
        return req.passthrough();
      });
      this.server.post('/auth/:mount/oidc/auth_url', (schema, req) => {
        // For JWT and OIDC
        this.assertReq(req);
        return req.passthrough();
      });
      this.server.post('/auth/:mount/login/:username', (schema, req) => {
        this.assertReq(req);
        return req.passthrough();
      });
      this.server.put('/auth/:mount/sso_service_url', (schema, req) => {
        // SAML only (enterprise)
        this.assertReq(req);
        return req.passthrough();
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
        assert.expect(isOidc ? 3 : 2);

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
        await fillIn(AUTH_FORM.selectMethod, type);

        if (type !== 'token') {
          // set custom mount
          await click(AUTH_FORM.advancedSettings);
          await fillIn(GENERAL.inputByAttr('path'), `custom-${type}`);
        }
        for (const key of backend.formAttributes) {
          // fill in all form items, except JWT which is not rendered
          if (key === 'jwt') return;
          await fillIn(GENERAL.inputByAttr(key), `some-${key}`);
        }
        await click(GENERAL.submitButton);
      });
    }
  });

  test('it does not call renew-self after successful login with non-renewable token', async function (assert) {
    this.server.post(
      '/auth/token/renew-self',
      () => new Error('should not call renew-self directly after logging in')
    );

    await login(rootToken);
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  module('Enterprise', function () {
    // this test is specifically to cover a token renewal bug within namespaces
    // namespace_path isn't returned by the renew-self response and so the auth service was
    // incorrectly setting userRootNamespace to '' (which denotes 'root'). this caused
    // subsequent capability checks fail because they would not be queried with the appropriate namespace header
    // if this test fails because a POST /v1/sys/capabilities-self returns a 403, then we have a problem!
    test('it sets namespace when renewing token', async function (assert) {
      // Sinon spy for clipboard
      const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();
      const uid = uuidv4();
      const ns = `admin-${uid}`;
      // log in to root to create namespace
      await login();
      await runCmd(createNS(ns), false);
      // login to namespace, mount userpass, create policy and user
      await loginNs(ns);
      const db = `database-${uid}`;
      const userpass = `userpass-${uid}`;
      const user = 'bob';
      const policyName = `policy-${userpass}`;
      const policy = `
        path "${db}/" {
          capabilities = ["list"]
        }
        path "${db}/roles" {
          capabilities = ["read","list"]
        }
        `;
      await runCmd([
        mountAuthCmd('userpass', userpass),
        mountEngineCmd('database', db),
        createPolicyCmd(policyName, policy),
        `write auth/${userpass}/users/${user} password=${user} token_policies=${policyName}`,
      ]);

      const inputValues = {
        username: user,
        password: user,
        path: userpass,
        namespace: ns,
      };

      // login as user just to get token (this is the only way to generate a token in the UI right now..)
      await loginMethod(inputValues, { authType: 'userpass', toggleOptions: true });
      await click(GENERAL.button('user-menu-trigger'));
      await click(GENERAL.copyButton);
      assert.true(clipboardSpy.calledOnce, 'Clipboard was called once');
      const token = clipboardSpy.firstCall.args[0];
      clipboardSpy.restore(); // restore original clipboard
      // login with token to reproduce bug
      await loginNs(ns, token);
      await visit(`/vault/secrets/${db}/overview?namespace=${ns}`);
      assert
        .dom('[data-test-overview-card="Roles"]')
        .hasText('Roles Create new', 'database overview renders');
      // renew token
      await click(GENERAL.button('user-menu-trigger'));
      await click('[data-test-user-menu-item="renew token"]');
      // navigate out and back to overview tab to re-request capabilities
      await click(GENERAL.secretTab('Roles'));
      await click(GENERAL.tab('overview'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${db}/overview?namespace=${ns}`,
        'it navigates to database overview'
      );

      // cleanup
      await visit(`/vault/logout?namespace=${ns}`);
      await fillIn(GENERAL.inputByAttr('namespace'), ''); // clear login form namespace input
      await login();
      // clean up namespace pollution
      await runCmd(deleteNS(ns));
    });

    test('it sets namespace header for sys/internal/ui/mounts request when namespace is inputted', async function (assert) {
      assert.expect(1);
      await visit('/vault/auth');

      this.server.get('/sys/internal/ui/mounts', (_, req) => {
        assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'admin', 'header contains namespace');
        return req.passthrough();
      });
      await typeIn(GENERAL.inputByAttr('namespace'), 'admin');
    });
  });
});
