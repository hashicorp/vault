import { currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';

import { storageKey } from 'vault/services/control-group';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authForm from 'vault/tests/pages/components/auth-form';
import controlGroup from 'vault/tests/pages/components/control-group';
import controlGroupSuccess from 'vault/tests/pages/components/control-group-success';
import authPage from 'vault/tests/pages/auth';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import listPage from 'vault/tests/pages/secrets/backend/list';

const consoleComponent = create(consoleClass);
const authFormComponent = create(authForm);
const controlGroupComponent = create(controlGroup);
const controlGroupSuccessComponent = create(controlGroupSuccess);

module('Acceptance | Enterprise | control groups', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  const POLICY = `
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

    path "kv-v2-mount/data/foo" {
      capabilities = ["create", "read", "update", "list"]
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

    path "kv-v2-mount/*" {
      capabilities = ["list"]
    }
  `;

  const AUTHORIZER_POLICY = `
    path "sys/control-group/authorize" {
      capabilities = ["update"]
    }

    path "sys/control-group/request" {
      capabilities = ["update"]
    }
  `;

  const ADMIN_USER = 'authorizer';
  const ADMIN_PASSWORD = 'test';
  const setupControlGroup = async context => {
    let userpassAccessor;
    await visit('/vault/secrets');
    await consoleComponent.toggle();
    await consoleComponent.runCommands([
      //enable kv-v1 mount and write a secret
      'write sys/mounts/kv type=kv',
      'write kv/foo bar=baz',

      //enable userpass, create user and associated entity
      'write sys/auth/userpass type=userpass',
      `write auth/userpass/users/${ADMIN_USER} password=${ADMIN_PASSWORD} policies=default`,
      `write identity/entity name=${ADMIN_USER} policies=test`,
      // write policies for control group + authorization
      `write sys/policies/acl/kv-control-group policy=${btoa(POLICY)}`,
      `write sys/policies/acl/authorizer policy=${btoa(AUTHORIZER_POLICY)}`,
      // read out mount to get the accessor
      'read -field=accessor sys/internal/ui/mounts/auth/userpass',
    ]);

    userpassAccessor = consoleComponent.lastTextOutput;

    await consoleComponent.runCommands([
      // lookup entity id for our authorizer
      `write -field=id identity/lookup/entity name=${ADMIN_USER}`,
    ]);
    let authorizerEntityId = consoleComponent.lastTextOutput;
    await consoleComponent.runCommands([
      // create alias for authorizor and add them to the managers group
      `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${ADMIN_USER}`,
      `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
      // create a token to request access to kv/foo
      'write -field=client_token auth/token/create policies=kv-control-group',
    ]);
    context.userToken = consoleComponent.lastLogOutput;

    await authPage.login(context.userToken);
    return this;
  };

  const writeSecret = async function(backend, path, key, val) {
    await listPage.visitRoot({ backend });
    await listPage.create();
    await editPage.createSecret(path, key, val);
  };

  test('for v2 secrets it redirects you if you try to navigate to a Control Group restricted path', async function(assert) {
    await consoleComponent.runCommands([
      'write sys/mounts/kv-v2-mount type=kv-v2',
      'delete kv-v2-mount/metadata/foo',
    ]);
    await writeSecret('kv-v2-mount', 'foo', 'bar', 'baz');
    await setupControlGroup(this);
    await visit('/vault/secrets/kv-v2-mount/show/foo');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.control-group-accessor',
      'redirects to access control group route'
    );
  });

  const workflow = async (assert, context, shouldStoreToken) => {
    let controlGroupToken;
    let accessor;
    let url = '/vault/secrets/kv/show/foo';
    await setupControlGroup(context);

    // as the requestor, go to the URL that's blocked by the control group
    // and store the values
    await visit(url);
    accessor = controlGroupComponent.accessor;
    controlGroupToken = controlGroupComponent.token;
    await authPage.logout();

    // log in as the admin, navigate to the accessor page,
    // and authorize the control group request
    await visit('/vault/auth?with=userpass');
    await authFormComponent.username(ADMIN_USER);
    await authFormComponent.password(ADMIN_PASSWORD);
    await authFormComponent.login();
    await visit(`/vault/access/control-groups/${accessor}`);
    await controlGroupComponent.authorize();
    assert.equal(controlGroupComponent.bannerPrefix, 'Thanks!', 'text display changes');
    await authPage.logout();

    await authPage.login(context.userToken);

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
      await controlGroupSuccessComponent.navigate();
      assert.equal(currentURL(), url, 'successfully loads the target url');
    } else {
      await visit(`/vault/access/control-groups/${accessor}`);
      await controlGroupSuccessComponent.token(controlGroupToken);
      await controlGroupSuccessComponent.unwrap();
      assert.ok(controlGroupSuccessComponent.showsJsonViewer, 'shows the json viewer');
    }
  };

  test('it allows the full flow to work with a saved token', async function(assert) {
    await workflow(assert, this, true);
  });

  test('it allows the full flow to work without a saved token', async function(assert) {
    await workflow(assert, this);
  });

  test('it displays the warning in the console when making a request to a Control Group path', async function(assert) {
    await setupControlGroup(this);
    await consoleComponent.toggle();
    await consoleComponent.runCommands('read kv/foo');
    let output = consoleComponent.lastLogOutput;
    assert.ok(output.includes('A Control Group was encountered at kv/foo'));
    assert.ok(output.includes('The Control Group Token is'));
    assert.ok(output.includes('The Accessor is'));
    assert.ok(output.includes('Visit /ui/vault/access/control-groups/'));
  });
});
