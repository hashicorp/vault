import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import setupMirage from 'ember-cli-mirage/test-support/setup-mirage';
import { create } from 'ember-cli-page-object';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

const tokenWithPolicy = async function(name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);

  return consoleComponent.lastLogOutput;
};

module('Acceptance | client count history', function(hooks) {
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
    await visit('/vault/metrics?tab=history');

    assert.equal(currentURL(), '/vault/metrics?tab=history');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('Data tracking is disabled');
  });

  test('it shows empty state when enabled and no data available', async function(assert) {
    server.create('metrics/config', { enabled: 'enable' });
    await visit('/vault/metrics?tab=history');

    assert.equal(currentURL(), '/vault/metrics?tab=history');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No monthly history');
  });

  test('it shows empty state when data available but not returned', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    await visit('/vault/metrics?tab=history');

    assert.equal(currentURL(), '/vault/metrics?tab=history');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-pricing-result-dates]').doesNotExist('Pricing metric result dates are not shown');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data received');
  });

  test('it shows warning when disabled and data available', async function(assert) {
    server.create('metrics/config', { queries_available: true, enabled: 'disable' });
    server.create('metrics/activity');
    await visit('/vault/metrics?tab=history');

    assert.equal(currentURL(), '/vault/metrics?tab=history');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-tracking-disabled]').exists('Flash message exists');
    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
  });

  test('it shows data when available from query', async function(assert) {
    server.create('metrics/config', { queries_available: true });
    server.create('metrics/activity');
    await visit('/vault/metrics?tab=history');

    assert.equal(currentURL(), '/vault/metrics?tab=history');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-configuration-tab]').exists('Metrics config tab exists');
    assert.dom('[data-test-tracking-disabled]').doesNotExist('Flash message does not exists');
    assert.dom('[data-test-client-count-stats]').exists('Client count metrics exists');
  });

  test('it shows metrics even if config endpoint not allowed', async function(assert) {
    server.create('metrics/activity');
    const deny_config_policy = `
    path "sys/internal/counters/config" {
      capabilities = ["deny"]
    },
    `;

    const userToken = await tokenWithPolicy('no-metrics-config', deny_config_policy);
    await logout.visit();
    await authPage.login(userToken);

    await visit('/vault/metrics?tab=history');

    assert.equal(currentURL(), '/vault/metrics?tab=history');
    assert.dom('[data-test-pricing-metrics-form]').exists('Pricing metrics date form exists');
    assert.dom('[data-test-configuration-tab]').doesNotExist('Metrics config tab does not exist');
    assert.dom('[data-test-tracking-disabled]').doesNotExist('Flash message does not exists');
    assert.dom('[data-test-client-count-stats]').exists('Client count metrics exists');
  });
});
