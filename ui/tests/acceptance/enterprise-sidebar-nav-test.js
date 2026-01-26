/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, visit } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';

const link = (label) => `[data-test-sidebar-nav-link="${label}"]`;
const panel = (label) => `[data-test-sidebar-nav-panel="${label}"]`;

module('Acceptance | Enterprise | sidebar navigation', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  // common links are tested in the sidebar-nav test and will not be covered here
  test('it should render enterprise only navigation links', async function (assert) {
    assert.dom(panel('Cluster')).exists('Cluster nav panel renders');

    await click(link('Secrets Sync'));
    assert.strictEqual(currentURL(), '/vault/sync/secrets/overview', 'Sync route renders');

    await click(link('Replication'));
    assert.strictEqual(currentURL(), '/vault/replication', 'Replication route renders');
    assert.dom(panel('Replication')).exists(`Replication nav panel renders`);
    assert.dom(link('Overview')).hasClass('active', 'Overview link is active');
    assert.dom(link('Performance')).exists('Performance link exists');
    assert.dom(link('Disaster Recovery')).exists('DR link exists');

    await click(link('Performance'));
    assert.strictEqual(
      currentURL(),
      '/vault/replication/performance',
      'Replication performance route renders'
    );

    // for some reason clicking this link would cause the testing browser locally
    // to navigate to 'vault/replication/dr' and halt the test runner
    assert
      .dom(link('Disaster Recovery'))
      .hasAttribute('href', '/ui/vault/replication/dr', 'Replication dr route renders');
    await visit('/vault');

    await click(link('Client Count'));
    assert.dom(panel('Client Count')).exists('Client Count nav panel renders');
    assert.dom(link('Client Usage')).hasClass('active', 'Client Usage link is active');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Client counts route renders');
    await click(link('Back to main navigation'));

    await click(link('License'));
    assert.strictEqual(currentURL(), '/vault/license', 'License route renders');
    await click(link('Access'));
    await click(link('Approval workflow'));
    assert.strictEqual(currentURL(), '/vault/access/control-groups', 'Approval workflow route renders');

    await click(link('Namespaces'));
    assert.strictEqual(currentURL(), '/vault/access/namespaces?page=1', 'Replication route renders');

    await click(link('Back to main navigation'));
    await click(link('Access'));
    await click(link('Role governing policies'));
    assert.strictEqual(currentURL(), '/vault/policies/rgp', 'Role governing policies route renders');

    await click(link('Endpoint governing policies'));
    assert.strictEqual(currentURL(), '/vault/policies/egp', 'Endpoint governing policies route renders');
  });

  test('it should link to correct routes at the access level', async function (assert) {
    assert.expect(12);

    await click(link('Access'));
    assert.dom(panel('Access')).exists('Access nav panel renders');

    const links = [
      { label: 'ACL policies', route: '/vault/policies/acl' },
      { label: 'Role governing policies', route: '/vault/policies/rgp' },
      { label: 'Endpoint governing policies', route: '/vault/policies/egp' },
      { label: 'Approval workflow', route: '/vault/access/control-groups' },
      { label: 'Leases', route: '/vault/access/leases/list' },
      { label: 'Authentication methods', route: '/vault/access' },
      { label: 'Multi-factor authentication', route: '/vault/access/mfa' },
      { label: 'OIDC provider', route: '/vault/access/oidc' },
      { label: 'Namespaces', route: '/vault/access/namespaces' },
      { label: 'Groups', route: '/vault/access/identity/groups' },
      { label: 'Entities', route: '/vault/access/identity/entities' },
    ];

    for (const l of links) {
      await click(link(l.label));
      assert.ok(currentURL().includes(l.route), `${l.label} route renders`);
    }
  });

  test('it should navigate to the correct links from Operational tools > Custom messages ember engine (enterprise)', async function (assert) {
    await click(link('Operational tools'));
    assert.strictEqual(currentURL(), '/vault/tools/wrap', 'Tool route renders');
    await click(link('Custom messages'));
    assert.strictEqual(currentURL(), '/vault/config-ui/messages', 'Custom messages route renders');
    await click(link('Lookup'));
    assert.strictEqual(currentURL(), '/vault/tools/lookup', 'Lookup route renders');
    await click(link('UI login settings'));
    assert.strictEqual(currentURL(), '/vault/config-ui/login-settings', 'UI login settings route renders');
  });
});
