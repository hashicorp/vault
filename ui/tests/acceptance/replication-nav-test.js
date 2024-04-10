/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import authPage from 'vault/tests/pages/auth';
import ENV from 'vault/config/environment';
import { click } from '@ember/test-helpers';

const SELECTORS = {
  navReplication: '[data-test-sidebar-nav-link="Replication"]',
  navPerformance: '[data-test-sidebar-nav-link="Performance"]',
  navDR: '[data-test-sidebar-nav-link="Disaster Recovery"]',
  title: '[data-test-replication-title]',
  primaryCluster: '[data-test-value-div="primary_cluster_addr"]',
  replicationSet: '[data-test-row-value="Replication set"]',
  knownSecondariesTitle: '.secondaries h3',
};
module('Acceptance | Enterprise | replication navigation', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'replication';
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('navigate between replication types updates page', async function (assert) {
    await click(SELECTORS.navReplication);
    assert.dom('[data-test-replication-header] h1').hasText('Disaster Recovery & Performance primary');
    await click(SELECTORS.navPerformance);

    // Ensure data is expected for performance
    assert.dom(SELECTORS.title).hasText('Performance primary');
    assert.dom(SELECTORS.primaryCluster).hasText('perf-foobar');
    assert.dom(SELECTORS.replicationSet).hasText('perf-cluster-id');
    assert.dom(SELECTORS.knownSecondariesTitle).hasText('0 Known secondaries');

    // Nav to DR and see updated data
    await click(SELECTORS.navDR);
    assert.dom(SELECTORS.title).hasText('Disaster Recovery primary');
    assert.dom(SELECTORS.primaryCluster).hasText('dr-foobar');
    assert.dom(SELECTORS.replicationSet).hasText('dr-cluster-id');
    assert.dom(SELECTORS.knownSecondariesTitle).hasText('1 Known secondaries');
  });
});
