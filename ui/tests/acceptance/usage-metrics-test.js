import { module, test } from 'qunit';
import { visit, currentURL, findAll } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import setupMirage from 'ember-cli-mirage/test-support/setup-mirage';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

module('Acceptance | usage metrics', function(hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  test('it shows empty state when disabled and no data available', async function(assert) {
    server.create('metrics/config', { enabled: 'disable' });
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data is being received');
  });

  test('it shows empty state when enabled and no data available', async function(assert) {
    server.create('metrics/config', { enabled: 'enable' });
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data is being received');
  });

  test('it shows empty state when data available but not returned', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-pricing-result-dates]').doesNotExist('Pricing metric result dates are not shown');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data found');
  });

  test('it shows warning when disabled and data available', async function(assert) {
    server.create('metrics/config', { queries_available: true, enabled: 'disable' });
    server.create('metrics/activity');
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-tracking-disabled]').exists('Flash message exists');
    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
  });

  test('it shows data when available from query', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    server.create('metrics/activity');
    await visit('/vault/metrics');

    assert.equal(currentURL(), '/vault/metrics');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-tracking-disabled]').doesNotExist('Flash message does not exists');
    assert.ok(findAll('.selectable-card').length === 3, 'renders the counts');
  });
});
