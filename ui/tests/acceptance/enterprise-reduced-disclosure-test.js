/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentRouteName, currentURL, fillIn, settled, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { createTokenCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import ENV from 'vault/config/environment';

const { unsealKeys } = VAULT_KEYS;
const SELECTORS = {
  footerVersion: `[data-test-footer-version]`,
  dashboardTitle: `[data-test-dashboard-card-header="Vault version"]`,
};

module('Acceptance | Enterprise | reduced disclosure test', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'reducedDisclosure';
  });
  hooks.beforeEach(function () {
    this.versionSvc = this.owner.lookup('service:version');
    return authPage.logout();
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it works when reduced disclosure enabled', async function (assert) {
    const namespace = 'reduced-disclosure';
    assert.dom(SELECTORS.footerVersion).hasText(`Vault`, 'shows Vault without version when logged out');
    await authPage.login();

    // Ensure it shows version on dashboard
    assert.dom(SELECTORS.dashboardTitle).includesText(`Vault v1.`);
    assert
      .dom(SELECTORS.footerVersion)
      .hasText(`Vault ${this.versionSvc.version}`, 'shows Vault version after login');

    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    await authPage.loginNs(namespace);

    assert
      .dom(SELECTORS.footerVersion)
      .hasText(`Vault ${this.versionSvc.version}`, 'shows Vault version within namespace');

    const token = await runCmd(createTokenCmd('default'));

    await authPage.logout();
    assert.dom(SELECTORS.footerVersion).hasText(`Vault`, 'no vault version after logout');

    await authPage.loginNs(namespace, token);
    assert
      .dom(SELECTORS.footerVersion)
      .hasText(`Vault ${this.versionSvc.version}`, 'shows Vault version for default policy in namespace');
  });

  test('it works for user accessing child namespace', async function (assert) {
    const namespace = 'reduced-disclosure';
    await authPage.login();

    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    const token = await runCmd(
      tokenWithPolicyCmd(
        'child-ns-access',
        `
    path "${namespace}/sys/*" {
      capabilities = ["read"]
    }
    `
      )
    );

    await authPage.logout();
    await authPage.login(token);
    assert
      .dom(SELECTORS.footerVersion)
      .hasText(`Vault ${this.versionSvc.version}`, 'shows Vault version for default policy in namespace');

    // navigate to child namespace
    await visit(`/vault/dashboard?namespace=${namespace}`);
    assert
      .dom(SELECTORS.footerVersion)
      .hasText(
        `Vault ${this.versionSvc.version}`,
        'shows Vault version for default policy in child namespace'
      );
    assert.dom(SELECTORS.dashboardTitle).includesText('Vault v1.');
  });

  test('shows correct version on unseal flow', async function (assert) {
    await authPage.login();

    const versionSvc = this.owner.lookup('service:version');
    await visit('/vault/settings/seal');
    assert
      .dom('[data-test-footer-version]')
      .hasText(`Vault ${versionSvc.version}`, 'shows version on seal page');
    assert.strictEqual(currentURL(), '/vault/settings/seal');

    // seal
    await click('[data-test-seal]');

    await click('[data-test-confirm-button]');

    await pollCluster(this.owner);
    await settled();
    assert.strictEqual(currentURL(), '/vault/unseal', 'vault is on the unseal page');
    assert.dom('[data-test-footer-version]').hasText(`Vault`, 'Clears version on unseal');

    // unseal
    for (const key of unsealKeys) {
      await fillIn('[data-test-shamir-key-input]', key);

      await click('button[type="submit"]');

      await pollCluster(this.owner);
      await settled();
    }

    assert.dom('[data-test-cluster-status]').doesNotExist('ui does not show sealed warning');
    assert.strictEqual(currentRouteName(), 'vault.cluster.auth', 'vault is ready to authenticate');
    assert.dom('[data-test-footer-version]').hasText(`Vault`, 'Version is still not shown before auth');
    await authPage.login();
    assert
      .dom('[data-test-footer-version]')
      .hasText(`Vault ${versionSvc.version}`, 'Version is shown after login');
  });

  test('does not allow access to replication pages', async function (assert) {
    await authPage.login();
    assert.dom('[data-test-sidebar-nav-link="Replication"]').doesNotExist('hides replication nav item');

    await visit(`/vault/replication/dr`);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.dashboard',
      'redirects to dashboard if replication access attempted'
    );
    assert.dom('[data-test-card="replication"]').doesNotExist('hides replication card on dashboard');
  });
});
