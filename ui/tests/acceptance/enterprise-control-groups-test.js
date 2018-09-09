import { currentURL, currentPath, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';

import { storageKey } from 'vault/services/control-group';
import console from 'vault/tests/pages/components/console/ui-panel';
import authForm from 'vault/tests/pages/components/auth-form';
import controlGroup from 'vault/tests/pages/components/control-group';
import controlGroupSuccess from 'vault/tests/pages/components/control-group-success';

const consoleComponent = create(console);
const authFormComponent = create(authForm);
const controlGroupComponent = create(controlGroup);
const controlGroupSuccessComponent = create(controlGroupSuccess);

module('Acceptance | Enterprise | control groups', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  hooks.afterEach(function() {
    return authLogout();
  });

  const POLICY = `'
    path "kv/foo" {
      capabilities = ["create", "read", "update", "delete", "list"]
      control_group = {
        max_ttl = "24h"
        factor "ops_manager" {
            identity {
                group_names = ["managers"]
                approvals = 1
            }
         }
      }
    }
  '`;

  const AUTHORIZER_POLICY = `'
    path "sys/control-group/authorize" {
      capabilities = ["update"]
    }

    path "sys/control-group/request" {
      capabilities = ["update"]
    }
  '`;

  const ADMIN_USER = 'authorizer';
  const ADMIN_PASSWORD = 'test';
  const setupControlGroup = async context => {
    let userpassAccessor;
    await visit('/vault/secrets');
    consoleComponent.toggle();
    consoleComponent.runCommands([
      //enable kv mount and write some data
      'write sys/mounts/kv type=kv',
      'write kv/foo bar=baz',
      //enable userpass, create user and associated entity
      'write sys/auth/userpass type=userpass',
      `write auth/userpass/users/${ADMIN_USER} password=${ADMIN_PASSWORD} policies=default`,
      `write identity/entity name=${ADMIN_USER} policies=test`,
      // write policies for control group + authorization
      `write sys/policies/acl/kv-control-group policy=${POLICY}`,
      `write sys/policies/acl/authorizer policy=${AUTHORIZER_POLICY}`,
      // read out mount to get the accessor
      'read -field=accessor sys/internal/ui/mounts/auth/userpass',
    ]);
    userpassAccessor = consoleComponent.lastTextOutput;
    consoleComponent.runCommands([
      // lookup entity id for our authorizer
      `write -field=id identity/lookup/entity name=${ADMIN_USER}`,
    ]);
    let authorizerEntityId = consoleComponent.lastTextOutput;
    consoleComponent.runCommands([
      // create alias for authorizor and add them to the managers group
      `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${ADMIN_USER}`,
      `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
      // create a token to request access to kv/foo
      'write -field=client_token auth/token/create policies=kv-control-group',
    ]);
    context.userToken = consoleComponent.lastLogOutput;
    authLogout();
    authLogin(context.userToken);
  };

  test('it redirects you if you try to navigate to a Control Group restricted path', async function(assert) {
    setupControlGroup(this);
    await visit('/vault/secrets/kv/show/foo');
    assert.equal(
      currentPath(),
      'vault.cluster.access.control-group-accessor',
      'redirects to access control group route'
    );
  });

  const workflow = async (assert, context, shouldStoreToken) => {
    let controlGroupToken;
    let accessor;
    let url = '/vault/secrets/kv/show/foo';
    setupControlGroup(context);

    // as the requestor, go to the URL that's blocked by the control group
    // and store the values
    await visit(url);
    accessor = controlGroupComponent.accessor;
    controlGroupToken = controlGroupComponent.token;
    authLogout();

    // log in as the admin, navigate to the accessor page,
    // and authorize the control group request
    await visit('/vault/auth?with=userpass');
    authFormComponent.username(ADMIN_USER);
    authFormComponent.password(ADMIN_PASSWORD);
    authFormComponent.login();
    await visit(`/vault/access/control-groups/${accessor}`);
    controlGroupComponent.authorize();
    assert.equal(controlGroupComponent.bannerPrefix, 'Thanks!', 'text display changes');
    authLogout();

    authLogin(context.userToken);

    if (shouldStoreToken) {
      localStorage.setItem(
        storageKey(accessor, 'kv/foo'),
        JSON.stringify({
          accessor,
          token: controlGroupToken,
          creation_path: 'kv/foo',
          uiParams: {
            url,
          },
        })
      );
      await visit(`/vault/access/control-groups/${accessor}`);
      assert.ok(controlGroupSuccessComponent.showsNavigateMessage, 'shows user the navigate message');
      controlGroupSuccessComponent.navigate();
      assert.equal(currentURL(), url, 'successfully loads the target url');
    } else {
      await visit(`/vault/access/control-groups/${accessor}`);
      controlGroupSuccessComponent.token(controlGroupToken);
      controlGroupSuccessComponent.unwrap();
      assert.ok(controlGroupSuccessComponent.showsJsonViewer, 'shows the json viewer');
    }
  };

  test('it allows the full flow to work with a saved token', function(assert) {
    workflow(assert, this, true);
  });

  test('it allows the full flow to work without a saved token', function(assert) {
    workflow(assert, this);
  });

  test('it displays the warning in the console when making a request to a Control Group path', function(assert) {
    setupControlGroup(this);
    consoleComponent.toggle();
    consoleComponent.runCommands('read kv/foo');
    let output = consoleComponent.lastLogOutput;
    assert.ok(output.includes('A Control Group was encountered at kv/foo'));
    assert.ok(output.includes('The Control Group Token is'));
    assert.ok(output.includes('The Accessor is'));
    assert.ok(output.includes('Visit /ui/vault/access/control-groups/'));
  });
});
