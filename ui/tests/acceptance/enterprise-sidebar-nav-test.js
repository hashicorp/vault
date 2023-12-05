/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import authPage from 'vault/tests/pages/auth';

const link = (label) => `[data-test-sidebar-nav-link="${label}"]`;
const panel = (label) => `[data-test-sidebar-nav-panel="${label}"]`;

module('Acceptance | Enterprise | sidebar navigation', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
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

    await click(link('Disaster Recovery'));
    assert.strictEqual(currentURL(), '/vault/replication/dr', 'Replication dr route renders');

    await click(link('Client Count'));
    assert.strictEqual(currentURL(), '/vault/clients/dashboard', 'Client counts route renders');

    await click(link('License'));
    assert.strictEqual(currentURL(), '/vault/license', 'License route renders');

    await click(link('Access'));
    await click(link('Control Groups'));
    assert.strictEqual(currentURL(), '/vault/access/control-groups', 'Control groups route renders');

    await click(link('Namespaces'));
    assert.strictEqual(currentURL(), '/vault/access/namespaces?page=1', 'Replication route renders');

    await click(link('Back to main navigation'));
    await click(link('Policies'));
    await click(link('Role-Governing Policies'));
    assert.strictEqual(currentURL(), '/vault/policies/rgp', 'Role-Governing Policies route renders');

    await click(link('Endpoint Governing Policies'));
    assert.strictEqual(currentURL(), '/vault/policies/egp', 'Endpoint Governing Policies route renders');
  });
});
