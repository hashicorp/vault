import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import {
  createPolicyCmd,
  deleteEngineCmd,
  mountEngineCmd,
  runCmd,
  createTokenCmd,
  mountAuthCmd,
} from 'vault/tests/helpers/commands';
import {
  adminPolicy,
  dataPolicy,
  metadataListPolicy,
  metadataPolicy,
} from 'vault/tests/helpers/policy-generator/kv';
import { writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';

/**
 * This test set is for testing delete, undelete, destroy flows
 */
module('Acceptance | kv-v2 workflow | delete, undelete, destroy', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `kv-workflow-${uuidv4()}`;
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeSecret(this.backend, 'app/first-secret', 'foo', 'bar');
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd('admin', adminPolicy(this.backend)),
        createTokenCmd('admin'),
      ]);
      await authPage.login(token);
    });
    test.todo('can soft delete and undelete a secret version', async function (assert) {
      assert.expect(0);
    });
    test.todo('can destroy a secret version', async function (assert) {
      assert.expect(0);
    });
    test.todo('can destroy a secret', async function (assert) {
      assert.expect(0);
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd('data-reader', dataPolicy({ backend: this.backend, capabilities: ['read'] })),
        createTokenCmd('data-reader'),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'data-reader-list',
          dataPolicy({ backend: this.backend, capabilities: ['read'] }) + metadataListPolicy(this.backend)
        ),
        createTokenCmd('data-reader-list'),
      ]);

      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'metadata-maintainer',
          metadataPolicy({ backend: this.backend }) + metadataListPolicy(this.backend)
        ),
        createTokenCmd('metadata-maintainer'),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'secret-creator',
          dataPolicy({ backend: this.backend, capabilities: ['create', 'update'] })
        ),
        createTokenCmd('secret-creator'),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      // Set up control group scenario
      const adminUser = 'admin';
      const adminPassword = 'password';
      const userpassMount = 'userpass';
      const POLICY = `
path "${this.backend}/data/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
  control_group = {
    max_ttl = "24h"
    factor "approver" {
      controlled_capabilities = ["write"]
      identity {
          group_names = ["managers"]
          approvals = 1
      }
    }
  }
}

path "${this.backend}/*" {
  capabilities = ["list"]
}

// Can we allow this so user can self-authorize?
path "sys/control-group/authorize" {
  capabilities = ["update"]
}

path "sys/control-group/request" {
  capabilities = ["update"]
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
      const userpassAccessor = await runCmd([
        //enable userpass, create user and associated entity
        mountAuthCmd('userpass', userpassMount),
        `write auth/${userpassMount}/users/${adminUser} password=${adminPassword} policies=default`,
        `write identity/entity name=${adminUser} policies=default`,
        // write policies for control group + authorization
        createPolicyCmd('kv-control-group', POLICY),
        createPolicyCmd('authorizer', AUTHORIZER_POLICY),
        // read out mount to get the accessor
        `read -field=accessor sys/internal/ui/mounts/auth/${userpassMount}`,
      ]);
      const authorizerEntityId = await runCmd(`write -field=id identity/lookup/entity name=${adminUser}`);
      const userToken = await runCmd([
        // create alias for authorizor and add them to the managers group
        `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${adminUser}`,
        `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
        // create a token to request access to kv/foo
        'write -field=client_token auth/token/create policies=kv-control-group',
      ]);
      this.userToken = userToken;
      await authPage.login(userToken);
    });
    // Copy test outline from admin persona
  });
});
