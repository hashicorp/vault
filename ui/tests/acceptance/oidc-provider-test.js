import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import authForm from 'vault/tests/pages/components/auth-form';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { visit, settled, currentURL, find } from '@ember/test-helpers';

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
const oidcEntity = async function(name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write identity/entity name="${OIDC_USER}" policies="${name}" metadata="email=vault@hashicorp.com" metadata="phone_number=123-456-7890"`,
    `read -field=id identity/entity/name/${OIDC_USER}`,
  ]);
  return consoleComponent.lastLogOutput;
  // b7762a9b-1a77-e4ba-47e6-d6d2ca87c0cb
};

const oidcGroup = async function(entityId) {
  await consoleComponent.runCommands([
    `write identity/group name="engineering" member_entity_ids="${entityId}"`,
    `read -field=id identity/group/name/engineering`,
  ]);
  return consoleComponent.lastLogOutput;
  // 9b6303d1-3fa3-9a50-327c-e11b057104a3
};

const authAccessor = async function(path = 'userpass') {
  await enablePage.enable('userpass', path);
  await consoleComponent.runCommands([
    `write auth/${path}/users/end-user password="${USER_PASSWORD}"`,
    `read -field=accessor sys/internal/ui/mounts/auth/${path}`,
  ]);
  return consoleComponent.lastLogOutput;
  // auth_userpass_d1b33f35
};

const entityAlias = async function(entityId, accessor, groupId) {
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
const setupWebapp = async function(redirect) {
  let webappName = `my-webapp-${new Date().getTime()}`;
  await consoleComponent.runCommands([
    `write identity/oidc/client/${webappName} redirect_uris="${redirect}" assignments="my-assignment" key="sigkey" id_token_ttl="30m" access_token_ttl="1h"`,
    `read -field=client_id identity/oidc/client/${webappName}`,
  ]);
  let output = consoleComponent.lastLogOutput;
  if (output.includes('error occurred')) {
    throw new Error(`OIDC setup failed: ${output}`);
  }
  return output;
  // RM5mU4p6y1xMzc24LjGd7pir85gNxQUA
};
const setupProvider = async function(clientId) {
  let providerName = `my-provider-${new Date().getTime()}`;
  await consoleComponent.runCommands(
    `write identity/oidc/provider/${providerName} allowed_client_ids="${clientId}" scopes="user,groups"`
  );
  return providerName;
};

const getAuthzUrl = (providerName, redirect, clientId, params) => {
  let queryParams = {
    client_id: clientId,
    nonce: 'abc123',
    redirect_uri: encodeURIComponent(redirect),
    response_type: 'code',
    scope: 'openid',
    state: 'foobar',
    ...params,
  };
  let queryString = Object.keys(queryParams).reduce((prev, key, idx) => {
    if (idx === 0) {
      return `${prev}${key}=${queryParams[key]}`;
    }
    return `${prev}&${key}=${queryParams[key]}`;
  }, '?');
  return `/vault/identity/oidc/provider/${providerName}/authorize${queryString}`;
  // /vault/identity/oidc/provider/my-provider-1635875876837/authorize?=true&client_id=aRMuW4lEZ3KQxFTOv8GK3Nt7Hd3aRlhg&nonce=abc123&redirect_uri=http%3A%2F%2F127.0.0.1%3A8251%2Fcallback&response_type=code&scope=openid&state=foobar
  // return `/vault/identity/oidc/provider/${providerName}/authorize?scope=openid&response_type=code&client_id=${clientId}&redirect_uri=${encodeURIComponent(
  //   redirect
  // )}&state=foobar&nonce=1234&prompt=none`;
};

const setupOidc = async function() {
  const callback = 'http://127.0.0.1:8251/callback';
  const entityId = await oidcEntity('oidc', OIDC_POLICY);
  const groupId = await oidcGroup(entityId);
  let authMethodPath = `userpass-${new Date().getTime()}`;
  const accessor = await authAccessor(authMethodPath);
  await entityAlias(entityId, accessor, groupId);
  const clientId = await setupWebapp(callback);
  let providerName = await setupProvider(clientId);
  return {
    providerName,
    callback,
    clientId,
    authMethodPath,
  };
};

module('Acceptance | oidc provider', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    await logout.visit();
    return authPage.login();
  });

  test('OIDC Provider logs in and redirects correctly', async function(assert) {
    let { providerName, callback, clientId, authMethodPath } = await setupOidc();

    await logout.visit();
    await settled();
    let url = getAuthzUrl(providerName, callback, clientId);
    visit(url);
    await settled();
    assert.ok(currentURL().startsWith('/vault/auth'), 'redirects to auth when no token');
    assert.ok(
      currentURL().includes(`redirect_to=${encodeURIComponent(url)}`),
      'encodes url for the query param'
    );
    // await this.pauseTest();
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
    assert.equal(currentURL(), url, 'URL is as expected after login');
    assert.dom('[data-test-oidc-redirect]').exists('redirect text exists');
    assert.ok(
      find('[data-test-oidc-redirect]').innerText.includes(`${callback}?code=`),
      'Successful redirect'
    );
    // let noPromptUrl = getAuthUrl(providerName, callback, clientId, { prompt: 'none' });
    // await logout.visit();
    // await settled();
    // visit(noPromptUrl);
    // await this.pauseTest();
  });

  test('OIDC Provider returns error if prompt=none and no token', async function(assert) {
    let { providerName, callback, clientId } = await setupOidc();

    await logout.visit();
    await settled();
    let url = getAuthzUrl(providerName, callback, clientId, { prompt: 'none' });
    console.log({ url });
    visit(url);
    await settled();
    await this.pauseTest();
  });
});
