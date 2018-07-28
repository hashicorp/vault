import { test } from 'qunit';
import { create } from 'ember-cli-page-object';

import { storageKey } from 'vault/services/control-group';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import console from 'vault/tests/pages/components/console/ui-panel';
import authForm from 'vault/tests/pages/components/auth-form';
import controlGroup from 'vault/tests/pages/components/control-group';
import controlGroupSuccess from 'vault/tests/pages/components/control-group-success';

const consoleComponent = create(console);
const authFormComponent = create(authForm);
const controlGroupComponent = create(controlGroup);
const controlGroupSuccessComponent = create(controlGroupSuccess);
moduleForAcceptance('Acceptance | Enterprise | control groups', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
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
const setupControlGroup = context => {
  let userpassAccessor;
  visit('/vault/secrets');
  consoleComponent.toggle();
  andThen(() => {
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
  });
  andThen(() => {
    userpassAccessor = consoleComponent.lastTextOutput;
    consoleComponent.runCommands([
      // lookup entity id for our authorizer
      `write -field=id identity/lookup/entity name=${ADMIN_USER}`,
    ]);
  });

  andThen(() => {
    let authorizerEntityId = consoleComponent.lastTextOutput;
    consoleComponent.runCommands([
      // create alias for authorizor and add them to the managers group
      `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${ADMIN_USER}`,
      `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
      // create a token to request access to kv/foo
      'write -field=client_token auth/token/create policies=kv-control-group',
    ]);
  });

  andThen(() => {
    context.userToken = consoleComponent.lastLogOutput;
  });
  authLogout();
  andThen(() => {
    authLogin(context.userToken);
  });
};

test('it redirects you if you try to navigate to a Control Group restricted path', function(assert) {
  setupControlGroup(this);
  visit('/vault/secrets/kv/show/foo');
  andThen(() => {
    assert.equal(
      currentPath(),
      'vault.cluster.access.control-group-accessor',
      'redirects to access control group route'
    );
  });
});

const workflow = (assert, context, shouldStoreToken) => {
  let controlGroupToken;
  let accessor;
  let url = '/vault/secrets/kv/show/foo';
  setupControlGroup(context);

  // as the requestor, go to the URL that's blocked by the control group
  // and store the values
  visit(url);
  andThen(() => {
    accessor = controlGroupComponent.accessor;
    controlGroupToken = controlGroupComponent.token;
  });
  authLogout();

  // log in as the admin, navigate to the accessor page,
  // and authorize the control group request
  visit('/vault/auth?with=userpass');
  andThen(() => {
    authFormComponent.username(ADMIN_USER);
    authFormComponent.password(ADMIN_PASSWORD);
    authFormComponent.login();
  });
  andThen(() => {
    visit(`/vault/access/control-groups/${accessor}`);
  });
  andThen(() => {
    controlGroupComponent.authorize();
  });
  andThen(() => {
    assert.equal(controlGroupComponent.bannerPrefix, 'Thanks!', 'text display changes');
  });
  authLogout();

  // log _back_ in as the requestor
  andThen(() => {
    authLogin(context.userToken);
  });

  if (shouldStoreToken) {
    // stuff localStorage full of the necessary details
    // so that they can nav back to the control group'd path in the UI
    andThen(() => {
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
    });
    andThen(() => {
      visit(`/vault/access/control-groups/${accessor}`);
    });
    andThen(() => {
      assert.ok(controlGroupSuccessComponent.showsNavigateMessage, 'shows user the navigate message');
    });
    controlGroupSuccessComponent.navigate();
    andThen(() => {
      assert.equal(currentURL(), url, 'successfully loads the target url');
    });
  } else {
    andThen(() => {
      visit(`/vault/access/control-groups/${accessor}`);
    });
    andThen(() => {
      controlGroupSuccessComponent.token(controlGroupToken);
      controlGroupSuccessComponent.unwrap();
    });
    andThen(() => {
      assert.ok(controlGroupSuccessComponent.showsJsonViewer, 'shows the json viewer');
    });
  }
};

test('it allows the full flow to work with a saved token', function(assert) {
  workflow(assert, this, true);
});

test('it allows the full flow to work without a saved token', function(assert) {
  workflow(assert, this);
});

test('it displays the warning in the console when making a request to a Control Group path', function(
  assert
) {
  setupControlGroup(this);
  andThen(() => {
    consoleComponent.toggle();
  });
  andThen(() => {
    consoleComponent.runCommands('read kv/foo');
  });
  andThen(() => {
    let output = consoleComponent.lastLogOutput;
    assert.ok(output.includes('A Control Group was encountered at kv/foo'));
    assert.ok(output.includes('The Control Group Token is'));
    assert.ok(output.includes('The Accessor is'));
    assert.ok(output.includes('Visit /ui/vault/access/control-groups/'));
  });
});
