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
import {
  createTokenCmd,
  deleteEngineCmd,
  mountEngineCmd,
  runCmd,
  tokenWithPolicyCmd,
} from 'vault/tests/helpers/commands';
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
    await visit('/vault/secrets-engines');
    await page.enableEngine();
    await click(GENERAL.cardContainer('aws'));
    await fillIn(GENERAL.inputByAttr('path'), 'aws_engine');
    await click(GENERAL.submitButton);
    await click(GENERAL.breadcrumbLink('Secrets engines'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'breadcrumb navigates to the list page'
    );
    assert.dom(GENERAL.tableData('aws_engine/', 'path')).hasTextContaining('aws_engine/');
    // cleanup
    await runCmd(deleteEngineCmd('aws_engine'));
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
    await visit('/vault/secrets-engines');
    await page.enableEngine();
    await click(GENERAL.cardContainer('aws'));
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
    await runCmd(deleteEngineCmd('aws'));
  });

  test('enterprise: cannot view list without permissions inside a namespace', async function (assert) {
    this.namespace = `ns-${this.uid}`;
    const enginePath1 = `kv-t1-${this.uid}`;
    const userDefault = await runCmd(createTokenCmd()); // creates a default user token

    await runCmd([`write sys/namespaces/${this.namespace} -force`]); // creates a namespace
    await loginNs(this.namespace); //logs into namespace with root token
    await runCmd(mountEngineCmd('kv', enginePath1)); // mounts a kv engine in namespace

    await loginNs(this.namespace, userDefault); // logs into that same namespace with a default user token

    await visit(`/vault/secrets-engines?namespace=${this.namespace}`); // nav to specified namespace list
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines?namespace=${this.namespace}`,
      'Should be on main secret engines list page within namespace.'
    );
    assert.dom(GENERAL.tableData(`${enginePath1}/`, 'path')).doesNotExist(); // without permissions, engine should not show for this user

    // cleanup namespace
    await login();
    await runCmd(`delete sys/namespaces/${this.namespace}`);
  });

  test('enterprise: can view list with permissions inside a namespace', async function (assert) {
    this.namespace = `ns-${this.uid}`;
    const enginePath1 = `kv-t2-${this.uid}`;
    const userToken = await runCmd(
      tokenWithPolicyCmd(
        'policy',
        `path "${this.namespace}/sys/*" {
          capabilities = ["create", "read", "update", "delete", "list"]
        }`
      )
    );

    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    await loginNs(this.namespace, userToken); // logs into namespace with user token
    await runCmd(mountEngineCmd('kv', enginePath1)); // mount kv engine as user

    await loginNs(this.namespace); // logs into namespace with root token

    await visit(`/vault/secrets-engines?namespace=${this.namespace}`); // nav to specified namespace list
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines?namespace=${this.namespace}`,
      'Should be on main secret engines list page within namespace.'
    );

    assert.dom(GENERAL.tableData(`${enginePath1}/`, 'path')).exists(); // with permissions, able to see the engine in list

    // cleanup namespace
    await login();
    await runCmd(`delete sys/namespaces/${this.namespace}`);
  });

  test('enterprise: it should navigate to cubbyhole list view in child namespace', async function (assert) {
    this.namespace = `ns-${this.uid}`;

    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    await loginNs(this.namespace);
    localStorage.setItem('dismissed-wizards', ['secret-engines']);
    await visit(`/vault/secrets-engines?namespace=${this.namespace}`);
    await click(`${GENERAL.tableData('cubbyhole/', 'path')} a`);

    assert.dom(GENERAL.emptyStateTitle).hasText('No secrets in this backend');

    // cleanup namespace
    await login();
    await runCmd(`delete sys/namespaces/${this.namespace}`);
  });

  test('after disabling it stays on the list view', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-disable-${this.uid}`;
    await runCmd(mountEngineCmd('alicloud', enginePath));
    await visit('/vault/secrets-engines');
    // to reduce flakiness, searching by engine name first in case there are pagination issues
    await fillIn(GENERAL.inputSearch('secret-engine-path'), enginePath);
    assert.dom(GENERAL.tableData(`${enginePath}/`, 'path')).exists('the alicloud engine is mounted');

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
    await click(`${GENERAL.tableData(`${enginePath1}/`, 'path')} a`);
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
    await click(`${GENERAL.tableData(`${enginePath1}/`, 'path')} a`);
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
});
