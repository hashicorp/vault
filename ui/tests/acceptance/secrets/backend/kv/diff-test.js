/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, settled, fillIn } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import listPage from 'vault/tests/pages/secrets/backend/list';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

module('Acceptance | kv2 diff view', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it shows correct diff status based on versions', async function (assert) {
    const secretPath = `my-secret`;

    await consoleComponent.runCommands([
      `write sys/mounts/secret type=kv options=version=2`,
      // delete any kv previously written here so that tests can be re-run
      `delete secret/metadata/${secretPath}`,
      'write -field=client_token auth/token/create policies=kv-v2-degrade',
    ]);

    await listPage.visitRoot({ backend: 'secret' });
    await settled();
    await listPage.create();
    await settled();
    await editPage.createSecret(secretPath, 'version1', 'hello');
    await settled();
    await click('[data-test-popup-menu-trigger="version"]');

    assert.dom('[data-test-view-diff]').doesNotExist('does not show diff view with only one version');
    // add another version
    await click('[data-test-secret-edit="true"]');

    const secondKey = document.querySelectorAll('[data-test-secret-key]')[1];
    const secondValue = document.querySelectorAll('.masked-value')[1];
    await fillIn(secondKey, 'version2');
    await fillIn(secondValue, 'world!');
    await click('[data-test-secret-save]');

    await click('[data-test-popup-menu-trigger="version"]');

    assert.dom('[data-test-view-diff]').exists('does show diff view with two versions');

    await click('[data-test-view-diff]');

    const diffBetweenVersion2and1 = document.querySelector('.jsondiffpatch-added').innerText;
    assert.strictEqual(diffBetweenVersion2and1, 'version2"world!"', 'shows the correct added part');

    await click('[data-test-popup-menu-trigger="right-version"]');

    await click('[data-test-rightSide-version="2"]');

    assert.dom('.diff-status').exists('shows States Match');
  });
});
