/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, currentURL, fillIn, settled, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const unsealKeys = ['unseal-key-1', 'unseal-key-2', 'unseal-key-3'];

module('Acceptance | unseal', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.unsealCount = 0;
    this.sealed = false;
    return login();
  });

  test('seal then unseal', async function (assert) {
    this.server.get(`/sys/seal-status`, () => {
      return {
        type: 'shamir',
        initialized: true,
        sealed: this.sealed,
      };
    });
    this.server.put(`/sys/seal`, () => {
      this.sealed = true;
      return overrideResponse(204);
    });
    this.server.put(`/sys/unseal`, () => {
      const threshold = unsealKeys.length;
      const attemptCount = this.unsealCount + 1;
      if (attemptCount >= threshold) {
        this.sealed = false;
      }
      this.unsealCount = attemptCount;
      return {
        sealed: attemptCount < threshold,
        t: threshold,
        n: threshold,
        progress: attemptCount,
      };
    });
    await visit('/vault/settings/seal');

    assert.strictEqual(currentURL(), '/vault/settings/seal');

    // seal
    await click('[data-test-seal]');
    await click(GENERAL.confirmButton);

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
