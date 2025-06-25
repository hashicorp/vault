/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentRouteName, visit, currentURL } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { MOUNT_BACKEND_FORM } from '../helpers/components/mount-backend-form-selectors';
import page from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | secret-engine list view', function (hooks) {
  setupApplicationTest(hooks);

  const createSecret = async (path, key, value, enginePath) => {
    await click(SES.createSecretLink);
    await fillIn(SES.secretPath('create'), path);
    await fillIn(SES.secretKey('create'), key);
    await fillIn(GENERAL.inputByAttr(key), value);
    await click(GENERAL.saveButton);
    await click(SES.crumb(enginePath));
  };

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  // the new API service camelizes response keys, so this tests is to assert that does NOT happen when we re-implement it
  test('it does not camelize the secret mount path', async function (assert) {
    await visit('/vault/secrets');
    await page.enableEngine();
    await click(MOUNT_BACKEND_FORM.mountType('aws'));
    await fillIn(GENERAL.inputByAttr('path'), 'aws_engine');
    await click(GENERAL.saveButton);
    await click(GENERAL.breadcrumbLink('Secrets'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'breadcrumb navigates to the list page'
    );
    assert.dom(SES.secretsBackendLink('aws_engine')).hasTextContaining('aws_engine/');
    // cleanup
    await runCmd(deleteEngineCmd('aws_engine'));
  });

  test('after enabling an unsupported engine it takes you to list page', async function (assert) {
    await visit('/vault/secrets');
    await page.enableEngine();
    await click(MOUNT_BACKEND_FORM.mountType('nomad'));
    await click(GENERAL.saveButton);

    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backends', 'navigates to the list page');
    // cleanup
    await runCmd(deleteEngineCmd('nomad'));
  });

  test('after enabling a supported engine it takes you to mount page, can see configure and clicking breadcrumb takes you back to list page', async function (assert) {
    await visit('/vault/secrets');
    await page.enableEngine();
    await click(MOUNT_BACKEND_FORM.mountType('aws'));
    await click(GENERAL.saveButton);

    assert.dom(SES.configTab).exists();

    await click(GENERAL.breadcrumbLink('Secrets'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'breadcrumb navigates to the list page'
    );
    // cleanup
    await runCmd(deleteEngineCmd('aws'));
  });

  test('enterprise: cannot view list without permissions inside namespace', async function (assert) {
    this.version = 'enterprise';
    this.backend = `bk-${this.uid}`;
    this.namespace = `ns-${this.uid}`;
    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    await loginNs(this.namespace, ' ');

    await visit('/vault/secrets');
    assert.dom(SES.secretsBackendLink('cubbyhole')).doesNotExist();

    await logout();
  });

  test('enterprise: can view list with permissions inside namespace', async function (assert) {
    this.version = 'enterprise';
    this.backend = `bk-${this.uid}`;
    this.namespace = `ns-${this.uid}`;
    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    await loginNs(this.namespace);
    await visit('/vault/secrets');

    assert.dom(SES.secretsBackendLink('cubbyhole')).exists();
  });

  test('after disabling it stays on the list view', async function (assert) {
    // first mount an engine so we can disable it.
    const enginePath = `alicloud-disable-${this.uid}`;
    await runCmd(mountEngineCmd('alicloud', enginePath));
    await visit('/vault/secrets');
    // to reduce flakiness, searching by engine name first in case there are pagination issues
    await selectChoose(GENERAL.searchSelect.trigger('filter-by-engine-name'), enginePath);
    assert.dom(SES.secretsBackendLink(enginePath)).exists('the alicloud engine is mounted');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('disable-engine'));
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
    await visit('/vault/secrets');
    await click(SES.secretsBackendLink(enginePath1));
    for (let i = 0; i <= 15; i++) {
      await createSecret(`secret-${i}`, 'foo', 'bar', enginePath1);
    }

    // navigate and check that details view is shown from non-nested secrets
    await click(GENERAL.pagination.next);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/list?page=2`,
      'After clicking next page in navigates to the second page.'
    );
    await click(SES.secretLink(secretPath));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/show/${secretPath}`,
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
    await visit('/vault/secrets');
    await click(SES.secretsBackendLink(enginePath1));
    for (let i = 0; i <= 15; i++) {
      await createSecret(`${parentPath}/secret-${i}`, 'foo', 'bar', enginePath1);
    }

    // navigate and check that the children list view is shown from nested secrets
    await click(SES.secretLink(`${parentPath}/`));

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/list/${parentPath}/`,
      'After clicking a nested secret it navigates to the children list view.'
    );

    await click(GENERAL.pagination.next);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath1}/list/${parentPath}/?page=2`,
      'After clicking next page it navigates to the second page.'
    );

    // cleanup
    await runCmd(deleteEngineCmd(enginePath1));
  });
});
