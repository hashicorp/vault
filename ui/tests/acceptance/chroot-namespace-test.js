/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentRouteName } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { createTokenCmd, runCmd, tokenWithPolicyCmd } from '../helpers/commands';

const navLink = (title) => `[data-test-sidebar-nav-link="${title}"]`;
// Matches the chroot namespace on the mirage handler
const namespace = 'my-ns';

module('Acceptance | chroot-namespace enterprise ui', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'chrootNamespace';
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should render normally when chroot namespace exists', async function (assert) {
    await authPage.login();
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'goes to dashboard page');
    assert.dom('[data-test-badge-namespace]').includesText('root', 'Shows root namespace badge');
  });

  test('a user with default policy should see nav items', async function (assert) {
    await authPage.login();
    // Create namespace
    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    // Create user within the namespace
    await authPage.loginNs(namespace);
    const userDefault = await runCmd(createTokenCmd());

    await authPage.loginNs(namespace, userDefault);
    ['Dashboard', 'Secrets Engines', 'Access', 'Tools'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item for user with default policy`);
    });
    ['Policies', 'Client Count', 'Replication', 'Raft Storage', 'License', 'Seal Vault'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item for user with default policy`);
    });

    // cleanup namespace
    await authPage.login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });

  test('a user with read policy should see nav items', async function (assert) {
    await authPage.login();
    // Create namespace
    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    // Create user within the namespace
    await authPage.loginNs(namespace);
    const reader = await runCmd(
      tokenWithPolicyCmd(
        'read-all',
        `
    path "sys/*" {
      capabilities = ["read"]
    }
    `
      )
    );

    await authPage.loginNs(namespace, reader);
    ['Dashboard', 'Secrets Engines', 'Access', 'Policies', 'Tools', 'Client Count'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item for user with read access policy`);
    });
    ['Replication', 'Raft Storage', 'License', 'Seal Vault'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item for user with read access policy`);
    });

    // cleanup namespace
    await authPage.login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });

  test('it works within a child namespace', async function (assert) {
    await authPage.login();
    // Create namespace
    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    // Create user within the namespace
    await authPage.loginNs(namespace);
    const childReader = await runCmd(
      tokenWithPolicyCmd(
        'read-child',
        `
        path "child/sys/*" {
          capabilities = ["read"]
        }
        `
      )
    );
    // Create child namespace
    await runCmd(`write sys/namespaces/child -f`, false);

    await authPage.loginNs(namespace, childReader);
    ['Dashboard', 'Secrets Engines', 'Access', 'Tools'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item`);
    });
    ['Policies', 'Client Count', 'Replication', 'Raft Storage', 'License', 'Seal Vault'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item`);
    });

    await authPage.loginNs(`${namespace}/child`, childReader);
    ['Dashboard', 'Secrets Engines', 'Access', 'Policies', 'Tools', 'Client Count'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item within child namespace`);
    });
    ['Replication', 'Raft Storage', 'License', 'Seal Vault'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item within child namespace`);
    });

    // cleanup namespaces
    await authPage.loginNs(namespace);
    await runCmd(`delete sys/namespaces/child`);
    await authPage.login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });
});
