/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentRouteName, visit, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { login, loginNs } from 'vault/tests/helpers/auth/auth-helpers';
import page from 'vault/tests/pages/settings/mount-secret-backend';
import localStorage from 'vault/lib/local-storage';

module('Acceptance | secret-engine list view', function (hooks) {
  setupApplicationTest(hooks);

  const createSecret = async (path, key, value, enginePath) => {
    await click(SES.createSecretLink);
    await fillIn(SES.secretPath('create'), path);
    await fillIn(SES.secretKey('create'), key);
    await fillIn(GENERAL.inputByAttr(key), value);
    await click(GENERAL.submitButton);
    await click(SES.crumb(enginePath));
  };

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    await login();
    // dismiss wizard
    localStorage.setItem('dismissed-wizards', ['secret-engines']);
  });

  // the new API service camelizes response keys, so this tests is to assert that does NOT happen when we re-implement it
  test('it does not camelize the secret mount path', async function (assert) {
    const path = `aws_${this.uid}`;
    await visit('/vault/secrets-engines');
    await page.enableEngine();
    await click(GENERAL.cardContainer('aws'));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.submitButton);
    await click(GENERAL.breadcrumbLink('Secrets engines'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'breadcrumb navigates to the list page'
    );
    await fillIn(GENERAL.inputSearch('secret-engine-path'), path);
    assert.dom(GENERAL.tableData(0, 'path')).hasText(`${path}/`);
    await runCmd(deleteEngineCmd(path));
  });

  test('after enabling an unsupported engine it takes you to list page', async function (assert) {
    await visit('/vault/secrets-engines');
    await page.enableEngine();
    await click(GENERAL.cardContainer('nomad'));
    await click(GENERAL.submitButton);

    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backends', 'navigates to the list page');
    // cleanup
    await runCmd(deleteEngineCmd('nomad'));
  });

  test('after enabling a supported engine it takes you to mount page, can see configure and clicking breadcrumb takes you back to list page', async function (assert) {
    const path = `aws-${this.uid}`;
    await visit('/vault/secrets-engines');
    await page.enableEngine();
    await click(GENERAL.cardContainer('aws'));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.submitButton);

    await click(GENERAL.dropdownToggle('Manage'));
    assert.dom(GENERAL.menuItem('Configure')).exists();

    await click(GENERAL.breadcrumbLink('Secrets engines'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'breadcrumb navigates to the list page'
    );
    // cleanup
    await runCmd(deleteEngineCmd(path));
  });

  test('after disabling it stays on the list view', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-disable-${this.uid}`;
    await runCmd(mountEngineCmd('alicloud', enginePath));
    await visit('/vault/secrets-engines');
    // to reduce flakiness, searching by engine name first in case there are pagination issues
    await fillIn(GENERAL.inputSearch('secret-engine-path'), enginePath);
    assert
      .dom(GENERAL.tableData(0, 'path'))
      .hasTextContaining(`${enginePath}/`, 'the alicloud engine is mounted');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends list page'
    );
  });

  test('it allows navigation to a non-nested secret with pagination', async function (assert) {
    assert.expect(2);

    const enginePath1 = `kv-v1-${this.uid}`;
    const secretPath = 'secret-9';
    await runCmd(mountEngineCmd('kv', enginePath1));

    // check kv1
    await visit('/vault/secrets-engines');
    await fillIn(GENERAL.inputSearch('secret-engine-path'), enginePath1);
    await click(GENERAL.linkTo(`${enginePath1}/`));
    for (let i = 0; i <= 15; i++) {
      await createSecret(`secret-${i}`, 'foo', 'bar', enginePath1);
    }

    // navigate and check that details view is shown from non-nested secrets
    await click(GENERAL.nextPage);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${enginePath1}/list?page=2`,
      'After clicking next page in navigates to the second page.'
    );
    await click(SES.secretLink(secretPath));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${enginePath1}/show/${secretPath}`,
      'After clicking a non-nested secret, it navigates to the details view.'
    );

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
  });

  test('it allows navigation to a nested secret with pagination', async function (assert) {
    assert.expect(2);

    const enginePath1 = `kv-v1-${this.uid}`;
    const parentPath = 'nested';

    await runCmd(mountEngineCmd('kv', enginePath1));

    // check kv1
    await visit('/vault/secrets-engines');
    await fillIn(GENERAL.inputSearch('secret-engine-path'), enginePath1);
    await click(GENERAL.linkTo(`${enginePath1}/`));
    for (let i = 0; i <= 15; i++) {
      await createSecret(`${parentPath}/secret-${i}`, 'foo', 'bar', enginePath1);
    }

    // navigate and check that the children list view is shown from nested secrets
    await click(SES.secretLink(`${parentPath}/`));

    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${enginePath1}/list/${parentPath}/`,
      'After clicking a nested secret it navigates to the children list view.'
    );

    await click(GENERAL.nextPage);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${enginePath1}/list/${parentPath}/?page=2`,
      'After clicking next page it navigates to the second page.'
    );

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
  });

  module('enterprise | namespaces', function (hooks) {
    hooks.beforeEach(async function () {
      await login();
      this.namespace = `ns-${this.uid}`;
      await runCmd([`write sys/namespaces/${this.namespace} -force`]);
      await loginNs(this.namespace); // log into namespace with root token
      // dismiss wizard
      localStorage.setItem('dismissed-wizards', ['secret-engines']);
    });

    // Ember route models won't refresh within a namespace when this.router.transitionTo() is called
    // because ?namespace is a query param that remains the same so the app doesn't detect any changes
    // and therefore does not refire the model hook.
    // this.router.refresh() must be called to refire model hooks and request fresh data.
    test('list refreshes after deleting an engine in a namespace', async function (assert) {
      const enginePath1 = `kv-t2-${this.uid}`;
      await runCmd(mountEngineCmd('kv', enginePath1)); // mount kv engine in the namespace
      await visit(`/vault/secrets-engines?namespace=${this.namespace}`); // nav to specified namespace list

      assert.dom(GENERAL.linkTo(`${enginePath1}/`)).exists();
      assert.dom(GENERAL.tableRow()).exists({ count: 2 }, 'only 2 secret engines are listed');
      // Delete the engine
      await click(`${GENERAL.listItem(`${enginePath1}/`)} ${GENERAL.menuTrigger}`);
      await click(GENERAL.menuItem('Delete'));
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backends',
        'redirects to the backends list page'
      );
      assert.dom(GENERAL.linkTo(enginePath1)).doesNotExist('deleted engine is no longer in list');
      assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only 1 secret engine is listed');
      // cleanup namespace
      await login();
      await runCmd(`delete sys/namespaces/${this.namespace}`);
    });

    test('it should navigate to cubbyhole list view in child namespace', async function (assert) {
      await visit(`/vault/secrets-engines?namespace=${this.namespace}`);
      await click(GENERAL.linkTo('cubbyhole/'));
      assert.dom(GENERAL.emptyStateTitle).hasText('No secrets in this backend');

      // cleanup namespace
      await login();
      await runCmd(`delete sys/namespaces/${this.namespace}`);
    });
  });
});
