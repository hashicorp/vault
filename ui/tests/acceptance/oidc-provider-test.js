import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { visit, settled, currentURL } from '@ember/test-helpers';

const consoleComponent = create(consoleClass);
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
    `write identity/entity name="end-user" policies="oidc" metadata="email=vault@hashicorp.com" metadata="phone_number=123-456-7890"`,
    `read -field=id identity/entity/name/end-user`,
  ]);

  return consoleComponent.lastLogOutput;
};

const oidcGroup = async function(entityId) {
  await consoleComponent.runCommands([
    `write identity/group name="engineering" member_entity_ids="${entityId}"`,
    `read -field=id identity/group/name/engineering`,
  ]);

  return consoleComponent.lastLogOutput;
};

const authAccessor = async function(path = 'userpass') {
  await enablePage.enable('userpass', path);
  await consoleComponent.runCommands([
    `write auth/${path}/users/end-user password="mypassword"`,
    `read -field=accessor sys/internal/ui/mounts/auth/userpass`,
  ]);
  return consoleComponent.lastLogOutput;
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
  return consoleComponent.lastLogOutput;
};
const setupWebapp = async function(redirect) {
  await consoleComponent.runCommands([
    `write identity/oidc/client/my-webapp redirect_uris="${redirect}" assignments="my-assignment" key="sigkey" id_token_ttl="30m" access_token_ttl="1h"`,
    `read -field=client_id identity/oidc/client/my-webapp`,
  ]);
  return consoleComponent.lastLogOutput;
};
const setupProvider = async function(clientId) {
  let providerName = `my-provider`;
  await consoleComponent.runCommands(
    `write identity/oidc/provider/${providerName} allowed_client_ids="${clientId}" scopes="user,groups"`
  );
  return providerName;
};

const getAuthUrl = (providerName, redirect, clientId) => {
  return `/vault/identity/oidc/provider/${providerName}/authorize?scope=openid&response_type=code&client_id=${clientId}&redirect_uri=${redirect}&state=foobar&nonce=1234&prompt=none`;
};

module('Acceptance | oidc provider', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    await logout.visit();
    return authPage.login();
  });

  test('OIDC Provider logs in and redirects correctly', async function(assert) {
    const callback = 'http://127.0.0.1:8251/callback';
    const entityId = await oidcEntity('oidc', OIDC_POLICY);
    const groupId = await oidcGroup(entityId);
    const accessor = await authAccessor('userpass');
    await entityAlias(entityId, accessor, groupId);
    const clientId = await setupWebapp(callback);
    const providerName = await setupProvider(clientId);

    // OIDC now set up
    await logout.visit();
    console.log('🔵🔵🔵LOGGED OUT🔵🔵🔵');
    await settled();
    let url = getAuthUrl(providerName, callback, clientId);
    visit(url);
    await settled();
    assert.equal(currentURL(), '/vault/auth', 'redirects to auth when no token');
    console.log(currentURL(), 'current url');
    await this.pauseTest();
  });
});
