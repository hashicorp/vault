/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, settled, click, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';

import { createSecret } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { assertSecretWrap } from 'vault/tests/helpers/components/secret-edit-toolbar';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | secrets/cubbyhole/create', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it creates and can view a secret with the cubbyhole backend', async function (assert) {
    assert.expect(4);
    const kvPath = `cubbyhole-kv-${this.uid}`;
    const requestPath = `cubbyhole/${kvPath}`;
    await listPage.visitRoot({ backend: 'cubbyhole' });
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'navigates to the list page'
    );

    await click(SES.createSecretLink);
    await createSecret(kvPath, 'foo', 'bar');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await assertSecretWrap(assert, this.server, requestPath);
  });

  test('it does not show the option to configure', async function (assert) {
    await visit(`/vault/secrets-engines/cubbyhole/list`);
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    assert.dom(GENERAL.tab('plugin-settings')).doesNotExist('does not show the configure button');
    // try to force it by visiting the URL
    await visit(`/vault/secrets-engines/cubbyhole/configuration/edit`);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(
        'Sorry, we were unable to find any content at /vault/secrets-engines/cubbyhole/configuration/edit.'
      );
  });
});
