import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { visit, currentURL, waitFor } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Acceptance | vault-reporting', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
  });

  test('it visits the usage reporting dashboard and renders the header', async function (assert) {
    await visit('/vault/usage-reporting');
    assert.strictEqual(currentURL(), '/vault/usage-reporting', 'navigates to usage reporting dashboard');
    assert.dom('.hds-page-header').includesText('Vault Usage', 'renders the "Vault Usage" header');
  });

  test('it renders the counters dashboard block with all expected counters', async function (assert) {
    await visit('/vault/usage-reporting');

    await waitFor('[data-test-dashboard-counters]');
    assert.dom('[data-test-dashboard-counters]').exists('renders the counters dashboard block');

    const expectedCounters = ['Child namespaces', 'KV secrets', 'Secrets sync', 'PKI roles'];

    expectedCounters.forEach((counterLabel) => {
      assert.dom(`[data-test-counter="${counterLabel}"]`).exists(`counter "${counterLabel}" is rendered`);
    });
  });

  test('dashboard card links work correctly', async function (assert) {
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-secret-engines]');

    assert.strictEqual(currentURL(), '/vault/usage-reporting', 'landed on reporting dashboard');

    // Secret Engines
    const secrets = document.querySelector('[data-test-dashboard-secret-engines]');
    const secretsLink = secrets.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(secretsLink, 'secret engines card title link exists');
    assert.strictEqual(secretsLink.getAttribute('href'), 'secrets', 'link points to secrets');

    // Auth Methods
    const auth = document.querySelector('[data-test-dashboard-auth-methods]');
    const authLink = auth.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(authLink, 'auth methods card title link exists');
    assert.strictEqual(authLink.getAttribute('href'), 'access', 'link points to access');

    // Lease Count Quotas
    const lease = document.querySelector('[data-test-dashboard-lease-count]');
    const leaseLink = lease.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(leaseLink, 'lease count quota card title link exists');
    assert.strictEqual(
      leaseLink.getAttribute('href'),
      'https://developer.hashicorp.com/vault/docs/enterprise/lease-count-quotas',
      'link points to external lease count docs'
    );

    // Cluster Replication
    const replication = document.querySelector('[data-test-dashboard-cluster-replication]');
    const replicationLink = replication.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(replicationLink, 'replication card title link exists');
    assert.strictEqual(replicationLink.getAttribute('href'), 'replication', 'link points to replication');
  });

  //
});
