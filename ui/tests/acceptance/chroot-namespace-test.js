/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentRouteName } from '@ember/test-helpers';
import { login, loginNs } from 'vault/tests/helpers/auth/auth-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import chrootNamespaceHandlers from 'vault/mirage/handlers/chroot-namespace';
import { createTokenCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';

const navLink = (title) => `[data-test-sidebar-nav-link="${title}"]`;
// Matches the chroot namespace on the mirage handler
const namespace = 'my-ns';

module('Acceptance | chroot-namespace enterprise ui', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    chrootNamespaceHandlers(this.server);
  });

  test('it should render normally when chroot namespace exists', async function (assert) {
    await login();
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'goes to dashboard page');
    assert.dom('[data-test-badge-namespace]').includesText('root', 'Shows root namespace badge');
  });

  test('root-only nav items are unavailable', async function (assert) {
    await login();

    ['Dashboard', 'Secrets Engines', 'Access', 'Operational tools'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item in chroot listener`);
    });
    // Client count is not root-only, but it is hidden for chroot
    ['Replication', 'Raft Storage', 'License', 'Seal Vault', 'Client Count'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item in chroot listener`);
    });

    // cleanup namespace
    await login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });

  test('a user with default policy should see nav items', async function (assert) {
    await login();
    // Create namespace
    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    // Create user within the namespace
    await loginNs(namespace);
    const userDefault = await runCmd(createTokenCmd());

    await loginNs(namespace, userDefault);
    [('Dashboard', 'Secrets Engines', 'Access', 'Operational tools')].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item for user with default policy`);
    });
    ['Client Count', 'Replication', 'Raft Storage', 'License', 'Seal Vault'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item for user with default policy`);
    });

    // cleanup namespace
    await login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });

  test('a user with read policy should see nav items', async function (assert) {
    await login();
    // Create namespace
    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    // Create user within the namespace
    await loginNs(namespace);
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

    await loginNs(namespace, reader);
    ['Dashboard', 'Secrets Engines', 'Access', 'Operational tools'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item for user with read access policy`);
    });
    ['Replication', 'Raft Storage', 'License', 'Seal Vault', 'Client Count'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item for user with read access policy`);
    });

    // cleanup namespace
    await login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });

  test('it works within a child namespace', async function (assert) {
    await login();
    // Create namespace
    await runCmd(`write sys/namespaces/${namespace} -f`, false);
    // Create user within the namespace
    await loginNs(namespace);
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

    await loginNs(namespace, childReader);
    ['Dashboard', 'Secrets Engines', 'Access', 'Operational tools'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item`);
    });
    ['Client Count', 'Replication', 'Raft Storage', 'License', 'Seal Vault'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item`);
    });

    await loginNs(`${namespace}/child`, childReader);

    ['Dashboard', 'Secrets Engines', 'Access', 'Operational tools'].forEach((nav) => {
      assert.dom(navLink(nav)).exists(`Shows ${nav} nav item within child namespace`);
    });
    ['Replication', 'Raft Storage', 'License', 'Seal Vault', 'Client Count'].forEach((nav) => {
      assert.dom(navLink(nav)).doesNotExist(`Does not show ${nav} nav item within child namespace`);
    });

    // cleanup namespaces
    await loginNs(namespace);
    await runCmd(`delete sys/namespaces/child`);
    await login();
    await runCmd(`delete sys/namespaces/${namespace}`);
  });
});
