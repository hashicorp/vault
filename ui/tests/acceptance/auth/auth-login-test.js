/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, fillIn, typeIn, visit, waitFor } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import { createNS, createPolicyCmd, deleteNS, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import {
  AUTH_METHOD_LOGIN_DATA,
  fillInLoginFields,
  login,
  loginNs,
  logout,
  SYS_INTERNAL_UI_MOUNTS,
} from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { v4 as uuidv4 } from 'uuid';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { supportedTypes } from 'vault/utils/auth-form-helpers';

const { rootToken } = VAULT_KEYS;

module('Acceptance | auth login', function (hooks) {
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
    const authTypes = supportedTypes(false);
    assert.expect(authTypes.length);
    for (const backend of authTypes.reverse()) {
      await visit(`/vault/auth?with=${backend}`);
      assert.dom(AUTH_FORM.selectMethod).hasValue(backend);
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

  // PAYLOAD TESTS FOR EACH AUTH METHOD
  // Assertion count is one for the URL and one for each payload key
  module('it sends the right payload when authenticating', function (hooks) {
    hooks.beforeEach(function () {
      this.assertAuthRequest = (assert, req, expectedPayload) => {
        const body = JSON.parse(req.requestBody);
        assert.true(true, `it calls the correct URL: ${req.url}`);

        for (const [expKey, expValue] of Object.entries(expectedPayload)) {
          assert.strictEqual(body[expKey], expValue, `payload includes ${expKey}: ${expValue}`);
        }
      };

      this.fillAndLogIn = async () => {
        await visit('/vault/auth');
        await fillIn(AUTH_FORM.selectMethod, this.authType);

        const loginData = { ...AUTH_METHOD_LOGIN_DATA[this.authType], path: `custom-${this.authType}` };
        await fillInLoginFields(loginData, { toggleOptions: true });
        await click(GENERAL.submitButton);
      };
    });

    test('token', async function (assert) {
      assert.expect(2);
      this.authType = 'token';
      this.expectedPayload = { 'x-vault-token': 'mysupersecuretoken' };
      const headerKey = 'x-vault-token';
      const expectedToken = this.expectedPayload[headerKey];

      this.server.get('/auth/token/lookup-self', (schema, req) => {
        const actualToken = req.requestHeaders[headerKey];
        assert.true(true, `it calls the correct URL: ${req.url}`);
        assert.strictEqual(actualToken, expectedToken, 'headers include token');
        req.passthrough();
      });

      await visit('/vault/auth');
      await fillIn(AUTH_FORM.selectMethod, this.authType);
      const loginData = AUTH_METHOD_LOGIN_DATA[this.authType];
      await fillInLoginFields(loginData);
      await click(GENERAL.submitButton);
    });

    test('github', async function (assert) {
      assert.expect(2);
      this.authType = 'github';
      this.expectedPayload = { token: 'mysupersecuretoken' };
      this.server.post('/auth/custom-github/login', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('ldap', async function (assert) {
      assert.expect(2);
      this.authType = 'ldap';
      this.expectedPayload = { password: 'some-password' };
      this.server.post('/auth/custom-ldap/login/matilda', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('jwt', async function (assert) {
      // auth_url is hit twice (once when inputs are filled and again on submit)
      // so the assertion count is doubled
      assert.expect(6);
      this.authType = 'jwt';
      this.expectedPayload = {
        redirect_uri: 'http://localhost:7357/ui/vault/auth/custom-jwt/oidc/callback',
        role: 'some-dev',
      };
      this.server.post('/auth/custom-jwt/oidc/auth_url', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('oidc', async function (assert) {
      // auth_url is hit twice (once when inputs are filled and again on submit)
      // so the assertion count is doubled
      assert.expect(6);
      this.authType = 'oidc';
      this.expectedPayload = {
        redirect_uri: 'http://localhost:7357/ui/vault/auth/custom-oidc/oidc/callback',
        role: 'some-dev',
      };
      this.server.post('/auth/custom-oidc/oidc/auth_url', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('okta', async function (assert) {
      assert.expect(2);
      this.authType = 'okta';
      this.expectedPayload = { password: 'some-password' };
      this.server.post('/auth/custom-okta/login/matilda', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('radius', async function (assert) {
      assert.expect(2);
      this.authType = 'radius';
      this.expectedPayload = { password: 'some-password' };
      this.server.post('/auth/custom-radius/login/matilda', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('userpass', async function (assert) {
      assert.expect(2);
      this.authType = 'userpass';
      this.expectedPayload = { password: 'some-password' };
      this.server.post('/auth/custom-userpass/login/matilda', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });

    test('enterprise: saml', async function (assert) {
      assert.expect(2);
      this.authType = 'saml';
      this.expectedPayload = { role: 'some-dev' };
      this.server.post('/auth/custom-saml/sso_service_url', (schema, req) => {
        this.assertAuthRequest(assert, req, this.expectedPayload);
        req.passthrough();
      });

      await this.fillAndLogIn();
    });
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
    // this test is to cover a token renewal bug within namespaces.
    // namespace_path isn't returned by the renew-self response and so the auth service was
    // incorrectly setting userRootNamespace to '' (which denotes 'root').
    // this caused subsequent capability checks to fail because they would not be queried with the appropriate namespace header.
    // if this test fails because a POST /v1/sys/capabilities-self returns a 403, then we have a problem!
    test('it sets namespace when renewing token', async function (assert) {
      const setTokenDataSpy = sinon.spy(this.owner.lookup('service:auth'), 'setTokenData');
      const uid = uuidv4();
      const ns = `admin-${uid}`;

      // log in to root to create namespace
      await login();
      await runCmd(createNS(ns), false);
      // log in to namespace, create policy and generate token
      await loginNs(ns);
      const db = `database-${uid}`;
      const policyName = `policy-${uid}`;
      const policy = `
        path "${db}/" {
          capabilities = ["list"]
        }
        path "${db}/roles" {
          capabilities = ["read","list"]
        }
        `;
      const token = await runCmd([
        mountEngineCmd('database', db),
        createPolicyCmd(policyName, policy),
        `write auth/token/create policies=${policyName} -field=client_token`,
      ]);

      // login with token to reproduce bug
      await loginNs(ns, token);
      await visit(`/vault/secrets/${db}/overview?namespace=${ns}`);
      assert
        .dom('[data-test-overview-card="Roles"]')
        .hasText('Roles Create new', 'database overview renders');

      // renew token
      await click(GENERAL.button('user-menu-trigger'));
      await click('[data-test-user-menu-item="renew token"]');
      // confirm setTokenData is called with correct args
      const [tokenName, { authMethodType, displayName, userRootNamespace, namespacePath }] =
        setTokenDataSpy.lastCall.args;
      assert.strictEqual(tokenName, 'vault-token☃1', 'setTokenData is called with tokenName');
      assert.strictEqual(authMethodType, 'token', 'setTokenData is called with authMethodType');
      assert.strictEqual(displayName, 'token', 'setTokenData is called with displayName');
      assert.strictEqual(userRootNamespace, ns, 'setTokenData is called with userRootNamespace');
      assert.strictEqual(namespacePath, `${ns}/`, 'setTokenData is called with namespacePath');

      // navigate out and back to overview tab to re-request capabilities
      // (before the bug fix, the view would not render and instead would show a 403)
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
        assert.strictEqual(req.requestHeaders['x-vault-namespace'], 'admin', 'header contains namespace');
        return req.passthrough();
      });
      await typeIn(GENERAL.inputByAttr('namespace'), 'admin');
    });
  });
});
