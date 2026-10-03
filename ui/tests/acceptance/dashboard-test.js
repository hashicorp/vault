/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import Sinon from 'sinon';

module('Acceptance | landing page dashboard', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.version = this.owner.lookup('service:version');
    this.namespace = this.owner.lookup('service:namespace');
  });

  test('navigate to dashboard on login', async function (assert) {
    await login();
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  test('display the version number for the title', async function (assert) {
    await login();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(`Vault ${this.version.versionDisplay}`);
  });

  test('display the namespace badge for enterprise', async function (assert) {
    await login();
    assert.dom('.hds-badge').hasText('root', 'shows root namespace for enterprise');
  });

  test('hides the configuration details card on a non-root namespace enterprise version', async function (assert) {
    // The route checks `inRootNamespace` so stub that return
    const nsStub = Sinon.stub(this.namespace, 'inRootNamespace').get(() => false);
    await login();
    await visit('/vault/dashboard');
    assert.dom(DASHBOARD.cardName('configuration-details')).doesNotExist();
    nsStub.restore();
  });
});
