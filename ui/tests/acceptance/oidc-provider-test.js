/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import { visit, settled, currentURL, waitFor, currentRouteName, fillIn, click } from '@ember/test-helpers';
import { clearRecord } from 'vault/tests/helpers/oidc-config';
import { runCmd } from 'vault/tests/helpers/commands';
import queryParamString from 'vault/utils/query-param-string';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const OIDC_USER = 'end-user';
const USER_PASSWORD = 'mypassword';
const PROVIDER_NAME = `my-provider-${uuidv4()}`;
const WEB_APP_NAME = `my-webapp-${uuidv4()}`;
const OIDC_POLICY = `path "identity/oidc/provider/+/userinfo" {
  capabilities = ["read", "update"]
}`;
const USER_TOKEN_TEMPLATE = `{
  "username": {{identity.entity.aliases.$USERPASS_ACCESSOR.name}},
  "contact": {
      "email": {{identity.entity.metadata.email}},
      "phone_number": {{identity.entity.metadata.phone_number}}
  }
}`;
const GROUP_TOKEN_TEMPLATE = `{
  "groups": {{identity.entity.groups.names}}
}`;
const oidcEntity = async function (name, policy) {
  return await runCmd([
    `write sys/policies/acl/${name} policy=${window.btoa(policy)}`,
    `write identity/entity name="${OIDC_USER}" policies="${name}" metadata="email=vault@hashicorp.com" metadata="phone_number=123-456-7890"`,
    `read -field=id identity/entity/name/${OIDC_USER}`,
  ]);
};

const oidcGroup = async function (entityId) {
  return await runCmd([
    `write identity/group name="engineering" member_entity_ids="${entityId}"`,
    `read -field=id identity/group/name/engineering`,
  ]);
};

const authAccessor = async function (path = 'userpass') {
  await enablePage.enable('userpass', path);
  return await runCmd([
    `write auth/${path}/users/end-user password="${USER_PASSWORD}"`,
    `read -field=accessor sys/internal/ui/mounts/auth/${path}`,
  ]);
};

const entityAlias = async function (entityId, accessor, groupId) {
  const userTokenTemplate = btoa(USER_TOKEN_TEMPLATE);
  const groupTokenTemplate = btoa(GROUP_TOKEN_TEMPLATE);

  const res = await runCmd([
    `write identity/entity-alias name="end-user" canonical_id="${entityId}" mount_accessor="${accessor}"`,
    `write identity/oidc/key/sigkey allowed_client_ids="*"`,
    `write identity/oidc/assignment/my-assignment group_ids="${groupId}" entity_ids="${entityId}"`,
    `write identity/oidc/scope/user description="scope for user metadata" template="${userTokenTemplate}"`,
    `write identity/oidc/scope/groups description="scope for groups" template="${groupTokenTemplate}"`,
  ]);
  return res.includes('Success');
};

const setupProvider = async function (clientId) {
  await runCmd(
    `write identity/oidc/provider/${PROVIDER_NAME} allowed_client_ids="${clientId}" scopes="user,groups"`
  );
};

const getAuthzUrl = (providerName, redirect, clientId, params) => {
  const queryParams = {
    client_id: clientId,
    nonce: 'abc123',
    redirect_uri: redirect,
    response_type: 'code',
    scope: 'openid',
    state: 'foobar',
    ...params,
  };
  const queryString = queryParamString(queryParams);
  return `/vault/identity/oidc/provider/${providerName}/authorize${queryString}`;
};

const setupOidc = async function (uid) {
  const callback = 'http://127.0.0.1:8251/callback';
  const entityId = await oidcEntity('oidc', OIDC_POLICY);
  const groupId = await oidcGroup(entityId);
  const authMethodPath = `oidc-userpass-${uid}`;
  const accessor = await authAccessor(authMethodPath);
  await entityAlias(entityId, accessor, groupId);
  await runCmd([
    `delete identity/oidc/client/${WEB_APP_NAME}`,
    `write identity/oidc/client/${WEB_APP_NAME} redirect_uris="${callback}" assignments="my-assignment" key="sigkey" id_token_ttl="30m" access_token_ttl="1h"`,
    `clear`,
    `read -field=client_id identity/oidc/client/${WEB_APP_NAME}`,
  ]);
  await settled();

  const clientId = await runCmd([`read -field=client_id identity/oidc/client/${WEB_APP_NAME}`]);
  await setupProvider(clientId);
  return {
    providerName: PROVIDER_NAME,
    callback,
    clientId,
    authMethodPath,
  };
};

module('Acceptance | oidc provider', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    this.store = this.owner.lookup('service:store');
    await login();
    await settled();
    this.oidcSetupInformation = await setupOidc(this.uid);
    return;
  });

  hooks.afterEach(async function () {
    await login();
  });

  test('OIDC Provider logs in and redirects correctly', async function (assert) {
    const { providerName, callback, clientId, authMethodPath } = this.oidcSetupInformation;
    await visit('/vault/access/oidc');
    assert
      .dom(`[data-test-oidc-client-linked-block='${WEB_APP_NAME}']`)
      .exists({ count: 1 }, 'shows webapp in oidc provider list');
    await visit('/vault/logout');
    const url = getAuthzUrl(providerName, callback, clientId);
    await visit(url);

    assert.ok(currentURL().startsWith('/vault/auth'), 'redirects to auth when no token');

    await waitFor('[data-test-auth-form]', { timeout: 5000 });
    assert.ok(
      currentURL().includes(`redirect_to=${encodeURIComponent(url)}`),
      `encodes url for the query param in: ${currentURL()}`
    );
    assert.dom('[data-test-auth-logo]').exists('Vault logo exists on auth page');
    assert
      .dom('[data-test-auth-helptext]')
      .hasText(
        'Once you log in, you will be redirected back to your application. If you require login credentials, contact your administrator.',
        'Has updated text for client authorization flow'
      );

    await fillIn(AUTH_FORM.selectMethod, 'userpass');
    await fillIn(GENERAL.inputByAttr('username'), OIDC_USER);
    await fillIn(GENERAL.inputByAttr('password'), USER_PASSWORD);
    await click(AUTH_FORM.advancedSettings);
    await fillIn(GENERAL.inputByAttr('path'), authMethodPath);
    await click(GENERAL.submitButton);
    assert.strictEqual(currentURL(), url, 'URL is as expected after login');
    assert
      .dom('[data-test-oidc-redirect]')
      .hasTextContaining(`click here to go back to app`, 'Shows link back to app');
    const link = document.querySelector('[data-test-oidc-redirect]').getAttribute('href');
    assert.ok(link.includes('/callback?code='), 'Redirects to correct url');

    //* clean up test state
    await clearRecord(this.store, 'oidc/client', WEB_APP_NAME);
    await clearRecord(this.store, 'oidc/provider', PROVIDER_NAME);
  });

  test('OIDC Provider redirects to auth if current token and prompt = login', async function (assert) {
    const { providerName, callback, clientId, authMethodPath } = this.oidcSetupInformation;
    await settled();
    await visit('/vault/dashboard');
    assert.strictEqual(currentURL(), '/vault/dashboard', 'User is logged in before oidc login attempt');
    const url = getAuthzUrl(providerName, callback, clientId, { prompt: 'login' });
    await visit(url);

    assert.ok(currentURL().startsWith('/vault/auth'), 'redirects to auth when no token');
    assert.notOk(
      currentURL().includes('prompt=login'),
      'Url params no longer include prompt=login after redirect'
    );

    await fillIn(AUTH_FORM.selectMethod, 'userpass');
    await click(AUTH_FORM.advancedSettings);
    await fillIn(GENERAL.inputByAttr('path'), authMethodPath);
    await fillIn(GENERAL.inputByAttr('username'), OIDC_USER);
    await fillIn(GENERAL.inputByAttr('password'), USER_PASSWORD);
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-oidc-redirect]')
      .hasTextContaining(`click here to go back to app`, 'Shows link back to app');
    const link = document.querySelector('[data-test-oidc-redirect]').getAttribute('href');
    assert.true(link.includes('/callback?code='), 'Redirects to correct url');
  });

  test('OIDC Provider shows consent form when prompt = consent', async function (assert) {
    const { providerName, callback, clientId, authMethodPath } = this.oidcSetupInformation;
    const url = getAuthzUrl(providerName, callback, clientId, { prompt: 'consent' });
    await visit('/vault/logout');
    await fillIn(AUTH_FORM.selectMethod, 'userpass');
    await click(AUTH_FORM.advancedSettings);
    await fillIn(GENERAL.inputByAttr('path'), authMethodPath);
    await fillIn(GENERAL.inputByAttr('username'), OIDC_USER);
    await fillIn(GENERAL.inputByAttr('password'), USER_PASSWORD);
    await click(GENERAL.submitButton);
    await visit(url);

    assert.notOk(
      currentURL().startsWith('/vault/auth'),
      'Does not redirect to auth because user is already logged in'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.oidc-provider');
    assert.dom('[data-test-consent-form]').exists('Consent form exists');

    //* clean up test state
    await clearRecord(this.store, 'oidc/client', WEB_APP_NAME);
    await clearRecord(this.store, 'oidc/provider', PROVIDER_NAME);
  });

  // Error handling test coverage, see issue for more context https://github.com/hashicorp/vault/issues/27772
  test('OIDC Provider redirects if authorization request throws a permission denied error', async function (assert) {
    this.auth = this.owner.lookup('service:auth');
    const { providerName, callback, clientId, authMethodPath } = this.oidcSetupInformation;
    // oidc provider authorization url, see https://developer.hashicorp.com/vault/docs/concepts/oidc-provider#authorization-endpoint
    const url = getAuthzUrl(providerName, callback, clientId);

    // stub ajax request made by the model hook in routes/vault/cluster/oidc-provider.js
    const authStub = sinon.stub(this.auth, 'ajax');
    authStub.rejects({
      json: () =>
        Promise.resolve({
          errors: ['2 errors occurred:\n\t* permission denied\n\t* invalid token\n\n'],
        }),
    });

    await visit('/vault/logout');

    // set spy here so they only spy on the relevant logic
    const deleteTokenSpy = sinon.spy(this.auth, 'deleteToken');

    // visit the OIDC authorization url to trigger the stubbed (and rejected) auth service ajax request
    await visit(url);

    await waitFor('[data-test-auth-form]', { timeout: 5000 });
    await fillIn(AUTH_FORM.selectMethod, 'userpass');
    await fillIn(GENERAL.inputByAttr('username'), OIDC_USER);
    await fillIn(GENERAL.inputByAttr('password'), USER_PASSWORD);
    await click(AUTH_FORM.advancedSettings);
    await fillIn(GENERAL.inputByAttr('path'), authMethodPath);
    await click(GENERAL.submitButton);

    // permission denied error redirect user to log in
    // if the route remains "vault.cluster.oidc-provider" - it did not redirect
    assert.strictEqual(currentRouteName(), 'vault.cluster.auth', 'it redirects to auth route');

    // assert permission denied error deletes OIDC user's token
    assert.true(
      deleteTokenSpy.calledOnce,
      'deleteToken is called because _redirectToAuth was called with logout:true'
    );

    //* clean up test state
    authStub.restore();
    await login();
    await clearRecord(this.store, 'oidc/client', WEB_APP_NAME);
    await clearRecord(this.store, 'oidc/provider', PROVIDER_NAME);
  });
});
