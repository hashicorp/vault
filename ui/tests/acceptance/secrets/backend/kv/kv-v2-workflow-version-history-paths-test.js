import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { personas } from 'vault/tests/helpers/policy-generator/kv';
import {
  clearRecords,
  deleteVersionCmd,
  destroyVersionCmd,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';

import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { click, currentRouteName, currentURL, visit, waitUntil } from '@ember/test-helpers';
import { grantAccess, setupControlGroup } from 'vault/tests/helpers/control-groups';

/**
 * This test set is for testing version history & path pages for secret.
 * Letter(s) in parenthesis at the end are shorthand for the persona,
 * for ease of tracking down specific tests failures from CI
 */
module('Acceptance | kv-v2 workflow | version history, paths', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.backend = `kv-workflow-${uuidv4()}`;
    this.secretPath = 'app/first-secret';
    this.urlPath = `${this.backend}/kv/${encodeURIComponent(this.secretPath)}`;
    this.navToSecret = async () => {
      return visit(`/vault/secrets/${this.urlPath}/details?version=4`);
    };
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeVersionedSecret(this.backend, this.secretPath, 'hello', 'there', 6);
    await runCmd([
      destroyVersionCmd(this.backend, this.secretPath),
      deleteVersionCmd(this.backend, this.secretPath, 2),
    ]);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(tokenWithPolicyCmd('admin', personas.admin(this.backend)));
      await authPage.login(token);
      clearRecords(this.store);
    });
    test('can navigate to the version history page (a)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Version History'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/metadata/versions`,
        'navigates to version history'
      );
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 6 });

      assert.dom(PAGE.versions.linkedBlock(6)).hasTextContaining('Version 6');
      assert.dom(PAGE.versions.icon(6)).hasTextContaining('Current');

      assert.dom(PAGE.versions.linkedBlock(2)).hasTextContaining('Version 2');
      assert.dom(PAGE.versions.icon(2)).hasTextContaining('Deleted');

      assert.dom(PAGE.versions.linkedBlock(1)).hasTextContaining('Version 1');
      assert.dom(PAGE.versions.icon(1)).hasText('Destroyed');

      await click(PAGE.versions.linkedBlock(5));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/details?version=5`,
        'navigates to detail at specific version'
      );
    });
    test('can navigate to the paths page (a)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/paths`,
        'navigates to secret paths route'
      );
      assert.dom(PAGE.infoRow).exists({ count: 3 }, 'shows 3 rows of information');
      assert.dom(PAGE.infoRowValue('API path')).hasText(`/v1/${this.backend}/data/${this.secretPath}`);
      assert.dom(PAGE.infoRowValue('CLI path')).hasText(`-mount="${this.backend}" "${this.secretPath}"`);
      assert
        .dom(PAGE.infoRowValue('API path for metadata'))
        .hasText(`/v1/${this.backend}/metadata/${this.secretPath}`);
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(tokenWithPolicyCmd('data-reader', personas.dataReader(this.backend)));
      await authPage.login(token);
      clearRecords(this.store);
    });
    test('cannot navigate to the version history page (dr)', async function (assert) {
      await this.navToSecret();
      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Does not render Version History tab');
    });
    test('can navigate to the paths page (dr)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/paths`,
        'navigates to secret paths route'
      );
      assert.dom(PAGE.infoRow).exists({ count: 3 }, 'shows 3 rows of information');
      assert.dom(PAGE.infoRowValue('API path')).hasText(`/v1/${this.backend}/data/${this.secretPath}`);
      assert.dom(PAGE.infoRowValue('CLI path')).hasText(`-mount="${this.backend}" "${this.secretPath}"`);
      assert
        .dom(PAGE.infoRowValue('API path for metadata'))
        .hasText(`/v1/${this.backend}/metadata/${this.secretPath}`);
    });
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd('data-list-reader', personas.dataListReader(this.backend))
      );
      await authPage.login(token);
      clearRecords(this.store);
    });
    test('cannot navigate to the version history page (dlr)', async function (assert) {
      await this.navToSecret();
      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Does not render Version History tab');
    });
    test('can navigate to the paths page (dlr)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/paths`,
        'navigates to secret paths route'
      );
      assert.dom(PAGE.infoRow).exists({ count: 3 }, 'shows 3 rows of information');
      assert.dom(PAGE.infoRowValue('API path')).hasText(`/v1/${this.backend}/data/${this.secretPath}`);
      assert.dom(PAGE.infoRowValue('CLI path')).hasText(`-mount="${this.backend}" "${this.secretPath}"`);
      assert
        .dom(PAGE.infoRowValue('API path for metadata'))
        .hasText(`/v1/${this.backend}/metadata/${this.secretPath}`);
    });
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd('metadata-maintainer', personas.metadataMaintainer(this.backend))
      );
      await authPage.login(token);
      clearRecords(this.store);
    });
    test('can navigate to the version history page (mm)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Version History'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/metadata/versions`,
        'navigates to version history'
      );
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 6 });

      assert.dom(PAGE.versions.linkedBlock(6)).hasTextContaining('Version 6');
      assert.dom(PAGE.versions.icon(6)).hasTextContaining('Current');

      assert.dom(PAGE.versions.linkedBlock(2)).hasTextContaining('Version 2');
      assert.dom(PAGE.versions.icon(2)).hasTextContaining('Deleted');

      assert.dom(PAGE.versions.linkedBlock(1)).hasTextContaining('Version 1');
      assert.dom(PAGE.versions.icon(1)).hasText('Destroyed');

      await click(PAGE.versions.linkedBlock(5));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/details?version=5`,
        'navigates to detail at specific version'
      );
    });
    test('can navigate to the paths page (mm)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/paths`,
        'navigates to secret paths route'
      );
      assert.dom(PAGE.infoRow).exists({ count: 3 }, 'shows 3 rows of information');
      assert.dom(PAGE.infoRowValue('API path')).hasText(`/v1/${this.backend}/data/${this.secretPath}`);
      assert.dom(PAGE.infoRowValue('CLI path')).hasText(`-mount="${this.backend}" "${this.secretPath}"`);
      assert
        .dom(PAGE.infoRowValue('API path for metadata'))
        .hasText(`/v1/${this.backend}/metadata/${this.secretPath}`);
    });
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(tokenWithPolicyCmd('secret-creator', personas.secretCreator(this.backend)));
      await authPage.login(token);
      clearRecords(this.store);
    });
    test('cannot navigate to the version history page (sc)', async function (assert) {
      await this.navToSecret();
      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Does not render Version History tab');
    });
    test('can navigate to the paths page (sc)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/paths`,
        'navigates to secret paths route'
      );
      assert.dom(PAGE.infoRow).exists({ count: 3 }, 'shows 3 rows of information');
      assert.dom(PAGE.infoRowValue('API path')).hasText(`/v1/${this.backend}/data/${this.secretPath}`);
      assert.dom(PAGE.infoRowValue('CLI path')).hasText(`-mount="${this.backend}" "${this.secretPath}"`);
      assert
        .dom(PAGE.infoRowValue('API path for metadata'))
        .hasText(`/v1/${this.backend}/metadata/${this.secretPath}`);
    });
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      const userPolicy = `
path "${this.backend}/metadata/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
  control_group = {
    max_ttl = "24h"
    factor "approver" {
      controlled_capabilities = ["read"]
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
`;
      const { userToken } = await setupControlGroup({ userPolicy });
      this.userToken = userToken;
      await authPage.login(userToken);
      clearRecords(this.store);
    });
    test('can navigate to the version history page (cg)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Version History'));
      assert.ok(
        await waitUntil(() => currentRouteName() === 'vault.cluster.access.control-group-accessor'),
        'redirects to access control group route'
      );
      await grantAccess({
        apiPath: `${this.backend}/metadata/${this.secretPath}`,
        originUrl: `/vault/secrets/${this.urlPath}/details`,
        userToken: this.userToken,
      });
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/details`,
        'navigates back to secret overview after authorized'
      );
      await click(PAGE.secretTab('Version History'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/metadata/versions`,
        'goes to version history page'
      );
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 6 });

      assert.dom(PAGE.versions.linkedBlock(6)).hasTextContaining('Version 6');
      assert.dom(PAGE.versions.icon(6)).hasTextContaining('Current');

      assert.dom(PAGE.versions.linkedBlock(2)).hasTextContaining('Version 2');
      assert.dom(PAGE.versions.icon(2)).hasTextContaining('Deleted');

      assert.dom(PAGE.versions.linkedBlock(1)).hasTextContaining('Version 1');
      assert.dom(PAGE.versions.icon(1)).hasText('Destroyed');

      await click(PAGE.versions.linkedBlock(5));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/details?version=5`,
        'navigates to detail at specific version'
      );
    });
    test('can navigate to the paths page (cg)', async function (assert) {
      await this.navToSecret();
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.urlPath}/paths`,
        'navigates to secret paths route'
      );
      assert.dom(PAGE.infoRow).exists({ count: 3 }, 'shows 3 rows of information');
      assert.dom(PAGE.infoRowValue('API path')).hasText(`/v1/${this.backend}/data/${this.secretPath}`);
      assert.dom(PAGE.infoRowValue('CLI path')).hasText(`-mount="${this.backend}" "${this.secretPath}"`);
      assert
        .dom(PAGE.infoRowValue('API path for metadata'))
        .hasText(`/v1/${this.backend}/metadata/${this.secretPath}`);
    });
  });
});
