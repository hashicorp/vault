/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { loginNs } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { MOUNT_BACKEND_FORM } from '../helpers/components/mount-backend-form-selectors';
import page from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | secret-engine list view', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('after enabling an unsupported engine it takes you to list page', async function (assert) {
    await visit('/vault/secrets');
    await page.enableEngine();
    await click(MOUNT_BACKEND_FORM.mountType('nomad'));
    await click(GENERAL.saveButton);

    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backends', 'navigates to the list page');
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
  });

  test('cannot view list without permissions inside namespace', async function (assert) {
    const uid = uuidv4();
    this.backend = `bk-${uid}`;
    this.namespace = `ns-${uid}`;
    // mount engine within namespace
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    // await login();
    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    await loginNs(this.namespace, ' ');

    await visit('/vault/secrets');
    assert.dom('[data-test-secrets-backend-link="kv"]').doesNotExist();
  });

  test('can view list with permissions inside namespace', async function (assert) {
    const uid = uuidv4();
    this.backend = `bk-${uid}`;
    this.namespace = `ns-${uid}`;
    // await login();
    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    await loginNs(this.namespace);
    // mount engine within namespace
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await visit('/vault/secrets');
    assert.dom('[data-test-secrets-backend-link="kv"]').exists();
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
});
