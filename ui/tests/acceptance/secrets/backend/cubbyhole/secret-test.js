/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';

import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import authPage from 'vault/tests/pages/auth';
import assertSecretWrap from 'vault/tests/helpers/secret-edit-toolbar';

module('Acceptance | secrets/cubbyhole/create', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it creates and can view a secret with the cubbyhole backend', async function (assert) {
    assert.expect(5);
    const kvPath = `cubbyhole-kv-${this.uid}`;
    const requestPath = `cubbyhole/${kvPath}`;
    await listPage.visitRoot({ backend: 'cubbyhole' });
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'navigates to the list page'
    );

    await listPage.create();
    await settled();
    await editPage.createSecret(kvPath, 'foo', 'bar');
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.dom('[data-test-created-time]').hasText('', 'it does not render created time if blank');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await assertSecretWrap(assert, this.server, requestPath);
  });
});
