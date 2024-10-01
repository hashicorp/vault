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
import { clearRecords, deleteLatestCmd, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { setupControlGroup } from 'vault/tests/helpers/control-groups';
import { click, currentRouteName, currentURL, waitUntil, visit } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const ALL_DELETE_ACTIONS = ['delete', 'destroy', 'undelete'];
const assertDeleteActions = (assert, expected = ['delete', 'destroy']) => {
  ALL_DELETE_ACTIONS.forEach((toolbar) => {
    if (expected.includes(toolbar)) {
      assert.dom(PAGE.detail[toolbar]).exists(`${toolbar} toolbar action exists`);
    } else {
      assert.dom(PAGE.detail[toolbar]).doesNotExist(`${toolbar} toolbar action not rendered`);
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
    // Versioned secret for testing delete is created (and deleted) by each module to avoid race condition failures
    return;
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      // patch is an enterprise feature but stubbing the version so assertions that check
      // patch actions exist before/after deletion can run on both CE and ent repos
      this.version = this.owner.lookup('service:version').type = 'enterprise';
      const token = await runCmd(makeToken('admin', this.backend, personas.admin));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can delete and undelete the latest secret version (a)', async function (assert) {
      assert.expect(21);
      const flashSuccess = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(PAGE.detail.patchLatest).exists();
      assert.dom(PAGE.infoRow).exists('shows secret data on load');
      // delete flow
      await click(PAGE.detail.delete);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(PAGE.detail.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(PAGE.detail.deleteOptionLatest).isNotDisabled('delete latest option is selectable');
      await click(PAGE.detail.deleteOptionLatest);
      await click(PAGE.detail.deleteConfirm);
      const expected = `Successfully deleted Version 4 of ${this.secretPath}.`;
      const [actual] = flashSuccess.lastCall.args;
      assert.strictEqual(actual, expected, 'renders correct flash message');

      // details update accordingly
      await click(PAGE.secretTab('Secret'));
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 4 of this secret has been deleted', 'Shows deleted message');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 4 deleted');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      assert.dom(PAGE.detail.patchLatest).doesNotExist('patching a deleted secret is not allowed');
      // undelete flow
      await click(PAGE.detail.undelete);
      // details update accordingly
      await click(PAGE.secretTab('Secret'));
      assert.dom(PAGE.infoRow).exists('shows secret data after undeleting');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 4 created');
      // correct toolbar options
      assertDeleteActions(assert, ['delete', 'destroy']);
      assert.dom(PAGE.detail.patchLatest).exists('patch is allowed after undeleting');
    });
    test('can soft delete and undelete an older secret version (a)', async function (assert) {
      assert.expect(19);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(PAGE.detail.patchLatest).exists();
      assert.dom(PAGE.infoRow).exists('shows secret data on load');
      // delete flow
      await click(PAGE.detail.delete);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(PAGE.detail.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(PAGE.detail.deleteOptionLatest).isNotDisabled('delete latest option is selectable');
      await click(PAGE.detail.deleteOption);
      await click(PAGE.detail.deleteConfirm);
      // we get navigated back to the overview page, so manually go back to deleted version
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 2 of this secret has been deleted', 'Shows deleted message');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 2 deleted');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      assert
        .dom(PAGE.detail.patchLatest)
        .exists('patching the latest version is allowed after deleting an older version');

      // undelete flow
      await click(PAGE.detail.undelete);
      // details update accordingly
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      assert.dom(PAGE.infoRow).exists('shows secret data after undeleting');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 2 created');
      // correct toolbar options
      assertDeleteActions(assert, ['delete', 'destroy']);
    });
    test('can destroy a secret version (a)', async function (assert) {
      assert.expect(12);
      const flashSuccess = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert);
      assert.dom(PAGE.detail.patchLatest).exists();
      // delete flow
      await click(PAGE.detail.destroy);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Destroy version?', 'modal has correct title');
      await click(PAGE.detail.deleteConfirm);
      const expected = `Successfully destroyed Version 3 of ${this.secretPath}.`;
      const [actual] = flashSuccess.lastCall.args;
      assert.strictEqual(actual, expected, 'renders correct flash message');
      // details update accordingly
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 3 of this secret has been permanently destroyed', 'Shows destroyed message');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('does not show version timestamp');
      // updated toolbar options
      assertDeleteActions(assert, []);
      assert
        .dom(PAGE.detail.patchLatest)
        .exists('patching the latest version is allowed after destroying an older version');
    });

    test('can permanently delete all secret versions (a)', async function (assert) {
      await writeVersionedSecret(this.backend, 'nuke', 'foo', 'bar', 2);
      await visit(`/vault/secrets/${this.backend}/kv/nuke/metadata`);
      assert.dom(PAGE.metadata.deleteMetadata).hasText('Permanently delete', 'shows delete metadata button');
      // delete flow
      await click(PAGE.metadata.deleteMetadata);
      assert
        .dom(PAGE.detail.deleteModalTitle)
        .includesText('Delete metadata and secret data?', 'modal has correct title');
      await click(PAGE.detail.deleteConfirm);
      await waitUntil(() => currentRouteName() === 'vault.cluster.secrets.backend.kv.list');
      // redirects to list
      assert.strictEqual(currentURL(), `/vault/secrets/${this.backend}/kv/list`, 'redirects to list');
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      // create and delete a secret as root user
      await authPage.login();
      await writeVersionedSecret(this.backend, 'nuke', 'foo', 'bar', 2);
      await runCmd(deleteLatestCmd(this.backend, 'nuke'));
      // login as data-reader persona
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
      assert.dom(PAGE.infoRow).exists('shows secret data');

      // data-reader can't delete, so check undelete with already-deleted version
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      assertDeleteActions(assert, []);
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 2 of this secret has been deleted', 'Shows deleted message');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 2 deleted');
    });
    test('cannot soft delete and undelete an older secret version (dr)', async function (assert) {
      assert.expect(4);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert, []);
      assert.dom(PAGE.infoRow).exists('shows secret data');
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
      await click(PAGE.secretTab('Metadata'));
      assert.dom(PAGE.metadata.deleteMetadata).doesNotExist('does not show delete metadata button');
    });
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      // create and delete a secret as root user
      await authPage.login();
      await writeVersionedSecret(this.backend, 'nuke', 'foo', 'bar', 2);
      await runCmd(deleteLatestCmd(this.backend, 'nuke'));
      // login as data-list-reader persona
      const token = await runCmd(makeToken('data-list-reader', this.backend, personas.dataListReader));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can delete and cannot undelete the latest secret version (dlr)', async function (assert) {
      assert.expect(12);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert, ['delete']);
      assert.dom(PAGE.infoRow).exists('shows secret data');
      // delete flow
      await click(PAGE.detail.delete);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(PAGE.detail.deleteOption).isDisabled('delete option is disabled');
      assert.dom(PAGE.detail.deleteOptionLatest).isNotDisabled('delete latest option is selectable');
      await click(PAGE.detail.deleteOptionLatest);
      await click(PAGE.detail.deleteConfirm);
      // details update accordingly
      await click(PAGE.secretTab('Secret'));
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 4 of this secret has been deleted', 'Shows deleted message');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 4 deleted');
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
      assert.dom(PAGE.infoRow).exists('shows secret data');
      // delete flow
      await click(PAGE.detail.delete);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(PAGE.detail.deleteOption).isDisabled('delete this version is not available');
    });
    test('cannot destroy a secret version (dlr)', async function (assert) {
      assert.expect(3);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert, ['delete']);
    });
    test('cannot permanently delete all secret versions (dlr)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(PAGE.secretTab('Metadata'));
      assert.dom(PAGE.metadata.deleteMetadata).doesNotExist('does not show delete metadata button');
    });
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      // create and delete a secret as root user
      await authPage.login();
      await writeVersionedSecret(this.backend, 'nuke', 'foo', 'bar', 2);
      await runCmd(deleteLatestCmd(this.backend, 'nuke'));
      // login as metadata-maintainer persona
      const token = await runCmd(makeToken('metadata-maintainer', this.backend, personas.metadataMaintainer));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cannot delete but can undelete the latest secret version (mm)', async function (assert) {
      assert.expect(18);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      // delete flow
      await click(PAGE.detail.delete);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(PAGE.detail.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(PAGE.detail.deleteOptionLatest).isDisabled('delete latest option is disabled');

      // Can't delete latest, try with pre-deleted secret
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      // undelete flow
      await click(PAGE.detail.undelete);
      await waitUntil(() => currentRouteName() === 'vault.cluster.secrets.backend.kv.secret.index');
      assert
        .dom(GENERAL.overviewCard.container('Current version'))
        .hasText(`Current version The current version of this secret. 2`);
      // details update accordingly
      await click(PAGE.secretTab('Secret'));
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // correct toolbar options
      assertDeleteActions(assert, ['delete', 'destroy']);
    });
    test('can soft delete and undelete an older secret version (mm)', async function (assert) {
      assert.expect(18);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // delete flow
      await click(PAGE.detail.delete);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Delete version?', 'shows correct modal title');
      assert.dom(PAGE.detail.deleteOption).isNotDisabled('delete option is selectable');
      assert.dom(PAGE.detail.deleteOptionLatest).isDisabled('delete latest option is disabled');
      await click(PAGE.detail.deleteOption);
      await click(PAGE.detail.deleteConfirm);
      // details update accordingly
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
      // updated toolbar options
      assertDeleteActions(assert, ['undelete', 'destroy']);
      // undelete flow
      await click(PAGE.detail.undelete);
      // details update accordingly
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version 2 timestamp not rendered');
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
      await click(PAGE.detail.destroy);
      assert.dom(PAGE.detail.deleteModalTitle).includesText('Destroy version?', 'modal has correct title');
      await click(PAGE.detail.deleteConfirm);
      // details update accordingly
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'Shows permissions message');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('does not show version timestamp');
      // updated toolbar options
      assertDeleteActions(assert, []);
    });
    test('cannot permanently delete all secret versions (mm)', async function (assert) {
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(PAGE.secretTab('Metadata'));
      assert.dom(PAGE.metadata.deleteMetadata).doesNotExist('does not show delete metadata button');
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
      const popupSelector = `${PAGE.list.item('bad-secret')} ${PAGE.popup}`;
      await click(popupSelector);
      assert.dom(PAGE.list.listMenuDelete).exists('shows the option to permanently delete');
    });
    test('can not delete all secret versions from root list view (snc)', async function (assert) {
      assert.expect(1);
      // go to root secret directory list view
      await visit(`/vault/secrets/${this.backend}/kv/list`);
      // shows overview card and not list view
      assert.dom(PAGE.list.overviewCard).exists('renders overview card');
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
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');

      // test with already deleted method
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('version timestamp not rendered');
      // updated toolbar options
      assertDeleteActions(assert, []);
    });
    test('cannot soft delete and undelete an older secret version (sc)', async function (assert) {
      assert.expect(4);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=2`);
      // correct toolbar options & details show
      assertDeleteActions(assert, []);
      assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
    });
    test('cannot destroy a secret version (sc)', async function (assert) {
      assert.expect(3);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/${this.secretPath}/details?version=3`);
      // correct toolbar options show
      assertDeleteActions(assert, []);
    });
    test('can permanently delete all secret versions (sc)', async function (assert) {
      await writeVersionedSecret(this.backend, 'nuke', 'foo', 'bar', 2);
      // go to secret details
      await visit(`/vault/secrets/${this.backend}/kv/nuke/details`);
      // Check metadata toolbar
      await click(PAGE.secretTab('Metadata'));
      assert.dom(PAGE.metadata.deleteMetadata).hasText('Permanently delete', 'shows delete metadata button');
      // delete flow
      await click(PAGE.metadata.deleteMetadata);
      assert
        .dom(PAGE.detail.deleteModalTitle)
        .includesText('Delete metadata and secret data?', 'modal has correct title');
      await click(PAGE.detail.deleteConfirm);
      await waitUntil(() => currentRouteName() === 'vault.cluster.secrets.backend.kv.list');
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
