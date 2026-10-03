/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import modifyPassthroughResponse from 'vault/mirage/helpers/modify-passthrough-response';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const link = (label) => `[data-test-sidebar-nav-link="${label}"]`;
const panel = (label) => `[data-test-sidebar-nav-panel="${label}"]`;

module('Acceptance | sidebar navigation', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    // set storage_type to raft to test link
    this.server.get('/sys/seal-status', (schema, req) => {
      return modifyPassthroughResponse(req, { storage_type: 'raft' });
    });
    this.server.get('/sys/storage/raft/configuration', () => this.server.create('configuration', 'withRaft'));
    setRunOptions({
      rules: {
        // TODO: fix use Dropdown on user-menu
        'nested-interactive': { enabled: false },
      },
    });
    return login();
  });

  test('it should navigate back to the dashboard when logo is clicked in app header', async function (assert) {
    await click('[data-test-app-header-logo]');
    assert.strictEqual(currentURL(), '/vault/dashboard', 'dashboard route renders');
  });

  test('it should link to correct routes at the cluster level', async function (assert) {
    assert.expect(9);

    assert.dom(panel('Cluster')).exists('Cluster nav panel renders');

    const subNavs = [
      { label: 'Access', route: 'policies/acl' },
      { label: 'Operational tools', route: 'tools/wrap' },
    ];

    for (const subNav of subNavs) {
      const { label, route } = subNav;
      await click(link(label));
      assert.strictEqual(currentURL(), `/vault/${route}`, `${label} route renders`);
      assert.dom(panel(label)).exists(`${label} nav panel renders`);
      await click(link('Back to main navigation'));
    }

    const links = [
      { label: 'Raft Storage', route: '/vault/storage/raft' },
      { label: 'Seal Vault', route: '/vault/settings/seal' },
      { label: 'Secrets Engines', route: '/vault/secrets-engines' },
      { label: 'Dashboard', route: '/vault/dashboard' },
    ];

    for (const l of links) {
      await click(link(l.label));
      assert.strictEqual(currentURL(), l.route, `${l.label} route renders`);
    }
  });

  test('it should link to correct routes at the access level', async function (assert) {
    assert.expect(8);

    await click(link('Access'));
    assert.dom(panel('Access')).exists('Access nav panel renders');

    const links = [
      { label: 'ACL policies', route: '/vault/policies/acl' },
      { label: 'Leases', route: '/vault/access/leases/list' },
      { label: 'Authentication methods', route: '/vault/access' },
      { label: 'Multi-factor authentication', route: '/vault/access/mfa' },
      { label: 'OIDC provider', route: '/vault/access/oidc' },
      { label: 'Groups', route: '/vault/access/identity/groups' },
      { label: 'Entities', route: '/vault/access/identity/entities' },
    ];

    for (const l of links) {
      await click(link(l.label));
      assert.ok(currentURL().includes(l.route), `${l.label} route renders`);
    }
  });

  test('it should link to correct routes at the tools level', async function (assert) {
    assert.expect(7);

    await click(link('Operational tools'));
    assert.dom(panel('Operational tools')).exists('Operational tools nav panel renders');

    const links = [
      { label: 'Wrap', route: '/vault/tools/wrap' },
      { label: 'Lookup', route: '/vault/tools/lookup' },
      { label: 'Unwrap', route: '/vault/tools/unwrap' },
      { label: 'Rewrap', route: '/vault/tools/rewrap' },
      { label: 'Random', route: '/vault/tools/random' },
      { label: 'Hash', route: '/vault/tools/hash' },
    ];

    for (const l of links) {
      await click(link(l.label));
      assert.strictEqual(currentURL(), l.route, `${l.label} route renders`);
    }
  });

  test('it should link to correct routes at the client counts level', async function (assert) {
    assert.expect(7);
    await click(link('Client Count'));
    assert.dom(panel('Client Count')).exists('Client counts nav panel renders');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Top level nav link renders overview');
    assert.dom(link('Client Usage')).hasClass('active');
    await click(link('Configuration'));
    assert.strictEqual(currentURL(), '/vault/clients/config', 'Clients configuration renders');
    assert.dom(link('Configuration')).hasClass('active');
    await click(link('Client Usage'));
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Sub nav link navigates to overview');
    assert.dom(link('Client Usage')).hasClass('active');
  });

  test('it should display access nav when mounting and configuring auth methods', async function (assert) {
    await click(link('Access'));
    await click('[data-test-sidebar-nav-link="Authentication methods"]');
    await click('[data-test-auth-enable]');
    assert.dom('[data-test-sidebar-nav-panel="Access"]').exists('Access nav panel renders');
    await click(link('Authentication methods'));
    await click(GENERAL.linkedBlock('token'));
    await click('[data-test-configure-link]');
    assert.dom('[data-test-sidebar-nav-panel="Access"]').exists('Access nav panel renders');
  });
});
