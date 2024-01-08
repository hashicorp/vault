/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, currentURL, fillIn, settled, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import authPage from 'vault/tests/pages/auth';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';

const { unsealKeys } = VAULT_KEYS;

module.skip('Acceptance | unseal', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('seal then unseal', async function (assert) {
    await visit('/vault/settings/seal');

    assert.strictEqual(currentURL(), '/vault/settings/seal');

    // seal
    await click('[data-test-seal]');

    await click('[data-test-confirm-button]');

    await pollCluster(this.owner);
    await settled();
    assert.strictEqual(currentURL(), '/vault/unseal', 'vault is on the unseal page');

    // unseal
    for (const key of unsealKeys) {
      await fillIn('[data-test-shamir-key-input]', key);

      await click('button[type="submit"]');

      await pollCluster(this.owner);
      await settled();
    }

    assert.dom('[data-test-cluster-status]').doesNotExist('ui does not show sealed warning');
    assert.strictEqual(currentRouteName(), 'vault.cluster.auth', 'vault is ready to authenticate');
  });
});
