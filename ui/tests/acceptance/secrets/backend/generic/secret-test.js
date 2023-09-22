/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import consolePanel from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';
import { writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

import { create } from 'ember-cli-page-object';

import apiStub from 'vault/tests/helpers/noop-all-api-requests';

const cli = create(consolePanel);

module('Acceptance | secrets/generic/create', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it creates and can view a secret with the generic backend', async function (assert) {
    const path = `generic-${this.uid}`;
    const kvPath = `generic-kv-${this.uid}`;
    await cli.runCommands([`write sys/mounts/${path} type=generic`, `write ${path}/foo bar=baz`]);
    await listPage.visitRoot({ backend: path });
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'navigates to the list page'
    );
    assert.strictEqual(listPage.secrets.length, 1, 'lists one secret in the backend');

    await listPage.create();
    await editPage.createSecret(kvPath, 'foo', 'bar');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  test('upgrading generic to version 2 lists all existing secrets, and CRUD continues to work', async function (assert) {
    const path = `generic-${this.uid}`;
    await cli.runCommands([
      `write sys/mounts/${path} type=generic`,
      `write ${path}/foo bar=baz`,
      // upgrade to version 2 generic mount
      `write sys/mounts/${path}/tune options=version=2`,
    ]);
    await visit(`/vault/secrets/${path}/kv/list`);

    assert
      .dom(PAGE.list.item('foo'))
      .exists('lists secret created under kv1 engine as secret in the kv2 list view');

    await writeSecret(path, 'bar', 'key', 'value');
    await visit(`/vault/secrets/${path}/kv/list`);

    ['foo', 'bar'].forEach((secret) => {
      assert.dom(PAGE.list.item(secret.path)).exists('lists both records');
    });
    assert.dom(PAGE.list.item()).exists({ count: 2 }, 'lists only the two secrets');
  });
});
