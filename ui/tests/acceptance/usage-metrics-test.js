import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import setupMirage from 'ember-cli-mirage/test-support/setup-mirage';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
// import metricsPage from 'vault/tests/pages/metrics';

module('Acceptance | usage metrics', function(hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function() {
    // this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  test('it shows message when disabled and no data available', async function(assert) {
    server.create('metrics/config');
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
  });

  test('it shows message when enabled and no data available', async function(assert) {
    server.create('metrics/config', { enabled: 'enable' });
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data is being received');
  });

  test('it shows message when disabled and data available but not returned', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    // server.create('metrics/activity');
    await visit('/vault/metrics');
    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data found');
  });

  test('it shows message when enabled and data available but not returned', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    // server.create('metrics/activity');
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data found');
  });

  test('it shows data when data available from query', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    server.create('metrics/activity');
    await visit('/vault/metrics');
    await this.pauseTest();
    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Empty state does not exist');
    assert.dom('[data-test-empty-state-title]').hasText('No data found');
  });
});
