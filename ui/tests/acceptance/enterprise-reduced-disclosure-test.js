/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { createTokenCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import ENV from 'vault/config/environment';

const SELECTORS = {
  footerVersion: `[data-test-footer-version]`,
  dashboardTitle: `[data-test-dashboard-card-header="Vault version"]`,
};

module('Acceptance | Enterprise | reduced disclosure test', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'mfaConfig';
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
});
