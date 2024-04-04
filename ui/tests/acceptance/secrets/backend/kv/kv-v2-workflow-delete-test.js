/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { personas } from 'vault/tests/helpers/kv/policy-generator';
import {
  clearRecords,
  deleteLatestCmd,
  setupControlGroup,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';
import { click, currentURL, visit } from '@ember/test-helpers';
import { KV_WORKFLOW } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { KV_METADATA_DETAILS } from 'vault/tests/helpers/components/kv/page/secret/metadata/details-selectors';
import { KV_SECRET } from 'vault/tests/helpers/components/kv/page/secret/details-selectors';

const ALL_DELETE_ACTIONS = ['delete', 'destroy', 'undelete'];
const assertDeleteActions = (assert, expected = ['delete', 'destroy']) => {
  ALL_DELETE_ACTIONS.forEach((toolbar) => {
    if (expected.includes(toolbar)) {
      assert.dom(KV_WORKFLOW.detail[toolbar]).exists(`${toolbar} toolbar action exists`);
    } else {
      assert.dom(KV_WORKFLOW.detail[toolbar]).doesNotExist(`${toolbar} toolbar action not rendered`);
    }
  });
};

const makeToken = (name, mountPath, policyGenerator) => {
  return tokenWithPolicyCmd(`${name}-${mountPath}`, policyGenerator(mountPath));
};

/**
 * This test set is for testing delete, undelete, destroy flows
 * Letter(s) in parenthesis at the end are shorthand for the persona,
 * for ease of tracking down specific tests failures from CI
 */
module('Acceptance | kv-v2 workflow | delete, undelete, destroy', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.backend = `kv-delete-${uuidv4()}`;
    this.secretPath = 'bad-secret';
    this.nestedSecretPath = 'app/nested/bad-secret';
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeVersionedSecret(this.backend, this.secretPath, 'foo', 'bar', 4);
    await writeVersionedSecret(this.backend, this.nestedSecretPath, 'foo', 'bar', 1);
    await writeVersionedSecret(this.backend, 'nuke', 'foo', 'bar', 2);
    // Delete latest version for testing undelete for users that can't delete
    await runCmd(deleteLatestCmd(this.backend, 'nuke'));
    return;
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(makeToken('admin', this.backend, personas.admin));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can delete and undelete the latest secret version (a)', async function (assert) {
      assert.expect(17);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      // delete flow
      await click(KV_SECRET.delete);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(KV_SECRET.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(KV_SECRET.deleteOptionLatest).isNotDisabled('delete latest option is selectable');
      await click(KV_SECRET.deleteOptionLatest);
      await click(KV_SECRET.deleteConfirm);
      // details update accordingly
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 4 of this secret has been deleted', 'Shows deleted message');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 4 deleted');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      // undelete flow
      await click(KV_SECRET.undelete);
      // details update accordingly
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 4 created');
      // correct toolbar options
      assertDeleteActions(assert, ['delete', 'destroy']);
    });
    test('can soft delete and undelete an older secret version (a)', async function (assert) {
      assert.expect(17);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      // delete flow
      await click(KV_SECRET.delete);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(KV_SECRET.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(KV_SECRET.deleteOptionLatest).isNotDisabled('delete latest option is selectable');
      await click(KV_SECRET.deleteOption);
      await click(KV_SECRET.deleteConfirm);
      // details update accordingly
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 2 of this secret has been deleted', 'Shows deleted message');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 2 deleted');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      // undelete flow
      await click(KV_SECRET.undelete);
      // details update accordingly
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 2 created');
      // correct toolbar options
      assertDeleteActions(assert, ['delete', 'destroy']);
    });
    test('can destroy a secret version (a)', async function (assert) {
      assert.expect(9);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert);
      // delete flow
      await click(KV_SECRET.destroy);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Destroy version?', 'modal has correct title');
      await click(KV_SECRET.deleteConfirm);
      // details update accordingly
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 3 of this secret has been permanently destroyed', 'Shows destroyed message');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('does not show version timestamp');
      // updated toolbar options
      assertDeleteActions(assert, []);
    });
    test('can permanently delete all secret versions (a)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(KV_METADATA_DETAILS.deleteMetadata)
        .hasText('Permanently delete', 'shows delete metadata button');
      // delete flow
      await click(KV_METADATA_DETAILS.deleteMetadata);
      assert
        .dom(KV_SECRET.deleteModalTitle)
        .includesText('Delete metadata and secret data?', 'modal has correct title');
      await click(KV_SECRET.deleteConfirm);

      // redirects to list
      assert.strictEqual(currentURL(), `/vault/secrets/${this.backend}/kv/list`, 'redirects to list');
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(makeToken('data-reader', this.backend, personas.dataReader));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cannot delete and undelete the latest secret version (dr)', async function (assert) {
      assert.expect(9);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert, []);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');

      // data-reader can't delete, so check undelete with already-deleted version
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      assertDeleteActions(assert, []);
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 2 of this secret has been deleted', 'Shows deleted message');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 2 deleted');
    });
    test('cannot soft delete and undelete an older secret version (dr)', async function (assert) {
      assert.expect(4);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert, []);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
    });
    test('cannot destroy a secret version (dr)', async function (assert) {
      assert.expect(3);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert, []);
    });
    test('cannot permanently delete all secret versions (dr)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(GENERAL.tab('Metadata'));
      assert.dom(KV_METADATA_DETAILS.deleteMetadata).doesNotExist('does not show delete metadata button');
    });
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(makeToken('data-list-reader', this.backend, personas.dataListReader));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can delete and undelete the latest secret version (dlr)', async function (assert) {
      assert.expect(12);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert, ['delete']);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      // delete flow
      await click(KV_SECRET.delete);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(KV_SECRET.deleteOption).isDisabled('delete option is disabled');
      assert.dom(KV_SECRET.deleteOptionLatest).isNotDisabled('delete latest option is selectable');
      await click(KV_SECRET.deleteOptionLatest);
      await click(KV_SECRET.deleteConfirm);
      // details update accordingly
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 4 of this secret has been deleted', 'Shows deleted message');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 4 deleted');
      // updated toolbar options
      assertDeleteActions(assert, []);
      // user can't undelete
    });
    test('can soft delete and undelete an older secret version (dlr)', async function (assert) {
      assert.expect(6);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert, ['delete']);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      // delete flow
      await click(KV_SECRET.delete);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(KV_SECRET.deleteOption).isDisabled('delete this version is not available');
    });
    test('cannot destroy a secret version (dlr)', async function (assert) {
      assert.expect(3);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert, ['delete']);
    });
    test('cannot permanently delete all secret versions (dr)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(GENERAL.tab('Metadata'));
      assert.dom(KV_METADATA_DETAILS.deleteMetadata).doesNotExist('does not show delete metadata button');
    });
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(makeToken('metadata-maintainer', this.backend, personas.metadataMaintainer));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can delete and undelete the latest secret version (mm)', async function (assert) {
      assert.expect(17);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      // delete flow
      await click(KV_SECRET.delete);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(KV_SECRET.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(KV_SECRET.deleteOptionLatest).isDisabled('delete latest option is disabled');

      // Can't delete latest, try with pre-deleted secret
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      // undelete flow
      await click(KV_SECRET.undelete);
      // details update accordingly
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // correct toolbar options
      assertDeleteActions(assert, ['delete', 'destroy']);
    });
    test('can soft delete and undelete an older secret version (mm)', async function (assert) {
      assert.expect(18);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // delete flow
      await click(KV_SECRET.delete);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(KV_SECRET.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(KV_SECRET.deleteOptionLatest).isDisabled('delete latest option is disabled');
      await click(KV_SECRET.deleteOption);
      await click(KV_SECRET.deleteConfirm);
      // details update accordingly
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      // undelete flow
      await click(KV_SECRET.undelete);
      // details update accordingly
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // correct toolbar options
      assertDeleteActions(assert);
    });
    test('can destroy a secret version (mm)', async function (assert) {
      assert.expect(9);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert);
      // delete flow
      await click(KV_SECRET.destroy);
      assert.dom(KV_SECRET.deleteModalTitle).includesText('Destroy version?', 'modal has correct title');
      await click(KV_SECRET.deleteConfirm);
      // details update accordingly
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'Shows permissions message');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('does not show version timestamp');
      // updated toolbar options
      assertDeleteActions(assert, []);
    });
    test('cannot permanently delete all secret versions (mm)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(GENERAL.tab('Metadata'));
      assert.dom(KV_METADATA_DETAILS.deleteMetadata).doesNotExist('does not show delete metadata button');
    });
  });

  module('secret-nested-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        makeToken('secret-nested-creator', this.backend, personas.secretNestedCreator)
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can delete all secret versions from the nested list view (snc)', async function (assert) {
      assert.expect(1);
      // go to nested secret directory list view
      await visit(`/vault/secrets/${this.backend}/kv/list/app/nested`);
      // correct popup menu items appear on list view
      const popupSelector = `${KV_WORKFLOW.list.item('bad-secret')} ${GENERAL.menuTrigger}`;
      await click(popupSelector);
      assert.dom(KV_WORKFLOW.list.listMenuDelete).exists('shows the option to permanently delete');
    });
    test('can not delete all secret versions from root list view (snc)', async function (assert) {
      assert.expect(1);
      // go to root secret directory list view
      await visit(`/vault/secrets/${this.backend}/kv/list`);
      // shows overview card and not list view
      assert.dom(KV_WORKFLOW.list.overviewCard).exists('renders overview card');
    });
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(makeToken('secret-creator', this.backend, personas.secretCreator));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cannot delete and undelete the latest secret version (sc)', async function (assert) {
      assert.expect(9);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert, []);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');

      // test with already deleted method
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(KV_SECRET.versionTimestamp).doesNotExist('version timestamp not rendered');
      // updated toolbar options
      assertDeleteActions(assert, []);
    });
    test('cannot soft delete and undelete an older secret version (sc)', async function (assert) {
      assert.expect(4);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert, []);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
    });
    test('cannot destroy a secret version (sc)', async function (assert) {
      assert.expect(3);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert, []);
    });
    test('can permanently delete all secret versions (sc)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(KV_METADATA_DETAILS.deleteMetadata)
        .hasText('Permanently delete', 'shows delete metadata button');
      // delete flow
      await click(KV_METADATA_DETAILS.deleteMetadata);
      assert
        .dom(KV_SECRET.deleteModalTitle)
        .includesText('Delete metadata and secret data?', 'modal has correct title');
      await click(KV_SECRET.deleteConfirm);

      // redirects to list
      assert.strictEqual(currentURL(), `/vault/secrets/${this.backend}/kv/list`, 'redirects to list');
    });
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      const userPolicy = `
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

      const { userToken } = await setupControlGroup({ userPolicy, backend: this.backend });
      this.userToken = userToken;
      await authPage.login(userToken);
      clearRecords(this.store);
      return;
    });
    // Copy test outline from admin persona
  });
});
