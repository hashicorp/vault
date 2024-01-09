/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import authForm from 'vault/tests/pages/components/auth-form';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { visit, settled, currentURL } from '@ember/test-helpers';
import { clearRecord } from 'vault/tests/helpers/oidc-config';
const consoleComponent = create(consoleClass);
const authFormComponent = create(authForm);

const OIDC_USER = 'end-user';
const USER_PASSWORD = 'mypassword';
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
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${window.btoa(policy)}`,
    `write identity/entity name="${OIDC_USER}" policies="${name}" metadata="email=vault@hashicorp.com" metadata="phone_number=123-456-7890"`,
    `read -field=id identity/entity/name/${OIDC_USER}`,
  ]);
  return consoleComponent.lastLogOutput;
};

const oidcGroup = async function (entityId) {
  await consoleComponent.runCommands([
    `write identity/group name="engineering" member_entity_ids="${entityId}"`,
    `read -field=id identity/group/name/engineering`,
  ]);
  return consoleComponent.lastLogOutput;
};

const authAccessor = async function (path = 'userpass') {
  await enablePage.enable('userpass', path);
  await consoleComponent.runCommands([
    `write auth/${path}/users/end-user password="${USER_PASSWORD}"`,
    `read -field=accessor sys/internal/ui/mounts/auth/${path}`,
  ]);
  return consoleComponent.lastLogOutput;
};

const entityAlias = async function (entityId, accessor, groupId) {
  const userTokenTemplate = btoa(USER_TOKEN_TEMPLATE);
  const groupTokenTemplate = btoa(GROUP_TOKEN_TEMPLATE);

  await consoleComponent.runCommands([
    `write identity/entity-alias name="end-user" canonical_id="${entityId}" mount_accessor="${accessor}"`,
    `write identity/oidc/key/sigkey allowed_client_ids="*"`,
    `write identity/oidc/assignment/my-assignment group_ids="${groupId}" entity_ids="${entityId}"`,
    `write identity/oidc/scope/user description="scope for user metadata" template="${userTokenTemplate}"`,
    `write identity/oidc/scope/groups description="scope for groups" template="${groupTokenTemplate}"`,
  ]);
  return consoleComponent.lastLogOutput.includes('Success');
};
const setupWebapp = async function (redirect) {
  const webappName = 'my-webapp';
  await consoleComponent.runCommands([
    `write identity/oidc/client/${webappName} redirect_uris="${redirect}" assignments="my-assignment" key="sigkey" id_token_ttl="30m" access_token_ttl="1h"`,
    `read -field=client_id identity/oidc/client/${webappName}`,
  ]);
  const output = consoleComponent.lastLogOutput;
  if (output.includes('error occurred')) {
    throw new Error(`OIDC setup failed: ${output}`);
  }
  return output;
};
const setupProvider = async function (clientId) {
  const providerName = `my-provider`;
  await consoleComponent.runCommands(
    `write identity/oidc/provider/${providerName} allowed_client_ids="${clientId}" scopes="user,groups"`
  );
  return providerName;
};

const getAuthzUrl = (providerName, redirect, clientId, params) => {
  const queryParams = {
    client_id: clientId,
    nonce: 'abc123',
    redirect_uri: encodeURIComponent(redirect),
    response_type: 'code',
    scope: 'openid',
    state: 'foobar',
    ...params,
  };
  const queryString = Object.keys(queryParams).reduce((prev, key, idx) => {
    if (idx === 0) {
      return `${prev}${key}=${queryParams[key]}`;
    }
    return `${prev}&${key}=${queryParams[key]}`;
  }, '?');
  return `/vault/identity/oidc/provider/${providerName}/authorize${queryString}`;
};

const setupOidc = async function (uid) {
  const callback = 'http://127.0.0.1:8251/callback';
  const entityId = await oidcEntity('oidc', OIDC_POLICY);
  const groupId = await oidcGroup(entityId);
  const authMethodPath = `oidc-userpass-${uid}`;
  const accessor = await authAccessor(authMethodPath);
  await entityAlias(entityId, accessor, groupId);
  const clientId = await setupWebapp(callback);
  const providerName = await setupProvider(clientId);
  return {
    providerName,
    callback,
    clientId,
    authMethodPath,
  };
};

module('Acceptance | oidc provider', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  test('OIDC Provider logs in and redirects correctly', async function (assert) {
    const { providerName, callback, clientId, authMethodPath } = await setupOidc(this.uid);

    await logout.visit();
    await settled();
    const url = getAuthzUrl(providerName, callback, clientId);
    await visit(url);

    assert.ok(currentURL().startsWith('/vault/auth'), 'redirects to auth when no token');
    assert.ok(
      currentURL().includes(`redirect_to=${encodeURIComponent(url)}`),
      'encodes url for the query param'
    );
    assert.dom('[data-test-auth-logo]').exists('Vault logo exists on auth page');
    assert
      .dom('[data-test-auth-helptext]')
      .hasText(
        'Once you log in, you will be redirected back to your application. If you require login credentials, contact your administrator.',
        'Has updated text for client authorization flow'
      );
    await authFormComponent.selectMethod(authMethodPath);
    await authFormComponent.username(OIDC_USER);
    await authFormComponent.password(USER_PASSWORD);
    await authFormComponent.login();
    await settled();
    assert.strictEqual(currentURL(), url, 'URL is as expected after login');
    assert
      .dom('[data-test-oidc-redirect]')
      .hasTextContaining(`click here to go back to app`, 'Shows link back to app');
    const link = document.querySelector('[data-test-oidc-redirect]').getAttribute('href');
    assert.ok(link.includes('/callback?code='), 'Redirects to correct url');

    //* clean up test state
    await clearRecord(this.store, 'oidc/client', 'my-webapp');
    await clearRecord(this.store, 'oidc/provider', 'my-provider');
  });

  test('OIDC Provider redirects to auth if current token and prompt = login', async function (assert) {
    const { providerName, callback, clientId, authMethodPath } = await setupOidc(this.uid);
    await settled();
    const url = getAuthzUrl(providerName, callback, clientId, { prompt: 'login' });
    await visit(url);

    assert.ok(currentURL().startsWith('/vault/auth'), 'redirects to auth when no token');
    assert.notOk(
      currentURL().includes('prompt=login'),
      'Url params no longer include prompt=login after redirect'
    );
    await authFormComponent.selectMethod(authMethodPath);
    await authFormComponent.username(OIDC_USER);
    await authFormComponent.password(USER_PASSWORD);
    await authFormComponent.login();
    await settled();
    assert
      .dom('[data-test-oidc-redirect]')
      .hasTextContaining(`click here to go back to app`, 'Shows link back to app');
    const link = document.querySelector('[data-test-oidc-redirect]').getAttribute('href');
    assert.ok(link.includes('/callback?code='), 'Redirects to correct url');

    //* clean up test state
    await clearRecord(this.store, 'oidc/client', 'my-webapp');
    await clearRecord(this.store, 'oidc/provider', 'my-provider');
  });

  test('OIDC Provider shows consent form when prompt = consent', async function (assert) {
    const { providerName, callback, clientId, authMethodPath } = await setupOidc(this.uid);
    const url = getAuthzUrl(providerName, callback, clientId, { prompt: 'consent' });
    await logout.visit();
    await authFormComponent.selectMethod(authMethodPath);
    await authFormComponent.username(OIDC_USER);
    await authFormComponent.password(USER_PASSWORD);
    await authFormComponent.login();
    await visit(url);

    assert.notOk(
      currentURL().startsWith('/vault/auth'),
      'Does not redirect to auth because user is already logged in'
    );
    assert.dom('[data-test-consent-form]').exists('Consent form exists');

    //* clean up test state
    await clearRecord(this.store, 'oidc/client', 'my-webapp');
    await clearRecord(this.store, 'oidc/provider', 'my-provider');
  });
});
