import { test } from 'qunit';
import { create } from 'ember-cli-page-object';

import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import console from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(console);

moduleForAcceptance('Acceptance | Enterprise | control groups', {
  beforeEach() {
    authLogin();
  },
});
const POLICY = `path "kv/foo" {
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
}`;

test('it creates a thing', function(assert) {
  visit('/vault/secrets');
  let token;
  andThen(() => {
    consoleComponent.runCommands([
      'write sys/mounts/kv type=kv',
      'write kv/foo bar=baz',
      'write sys/auth/userpass type=userpass',
      `write sys/policies/acl/kv-control-group policy='${POLICY}'`,
      'write -field=client_token auth/token/create policies=kv-control-group'
    ]);
  });
  andThen(() => {
    token = consoleComponent.lastLogOutput;
    authLogout();
    return authLogin(token);
  });
  visit('/vault/secrets/kv/show/foo');
  andThen(() => {
    debugger;
  });
})
