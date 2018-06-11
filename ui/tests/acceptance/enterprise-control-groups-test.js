import { test } from 'qunit';
import { create } from 'ember-cli-page-object';

import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import console from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(console);

moduleForAcceptance('Acceptance | Enterprise | control groups', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  }
});
const POLICY = `'path "kv/foo" {
    capabilities = ["create", "read", "update", "delete", "list"]
    control_group = {
        max_ttl = "24h"
        factor "ops_manager" {
            identity {
                group_names = ["managers"]
                approvals = 2
            }
        }
        factor "superman" {
            identity {
                group_names = ["superman"]
                approvals = 1
            }
        }
    }
}'`;

const setupControlGroup = () => {
  let token;
  visit('/vault/secrets');
  andThen(() => {
    consoleComponent.runCommands([
      'write sys/mounts/kv type=kv',
      'write kv/foo bar=baz',
      'write sys/auth/userpass type=userpass',
      `write sys/policies/acl/kv-control-group policy=${POLICY}`,
      'write -field=client_token auth/token/create policies=kv-control-group'
    ]);
  });
  andThen(() => {
    token = consoleComponent.lastLogOutput;
    authLogout();
    return authLogin(token);
  });
};

test('it redirects you if you try to navigate to a Control Group restricted path', function(assert) {
  setupControlGroup();
  visit('/vault/secrets/kv/show/foo');
  andThen(() => {
    assert.equal(
      currentPath(),
      'vault.cluster.access.control-group-accessor',
      'redirects to access control group route'
    );
  });
});

test('it displays the warning in the console when making a request to a Control Group path', function(assert) {
  setupControlGroup();
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
