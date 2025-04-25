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

  test('dashboard card: Secret engines', async function (assert) {
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-secret-engines]');

    const card = document.querySelector('[data-test-dashboard-secret-engines]');
    assert.ok(card, 'renders Secret engines card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Secret engines', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(link.getAttribute('href'), 'secrets', 'link points to secrets');

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Breakdown of secret engines for this namespace(s)',
      'description is correct'
    );
  });

  test('dashboard card: Authentication methods', async function (assert) {
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-auth-methods]');

    const card = document.querySelector('[data-test-dashboard-auth-methods]');
    assert.ok(card, 'renders Authentication methods card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Authentication methods', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(link.getAttribute('href'), 'access', 'link points to access');

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Breakdown of authentication methods',
      'description is correct'
    );
  });

  test('dashboard card: Global lease count quota', async function (assert) {
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-lease-count]');

    const card = document.querySelector('[data-test-dashboard-lease-count]');
    assert.ok(card, 'renders Global lease count quota card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Global lease count quota', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(
      link.getAttribute('href'),
      'https://developer.hashicorp.com/vault/docs/enterprise/lease-count-quotas',
      'link points to lease count docs'
    );

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Snapshot of global lease count quota consumption',
      'description is correct'
    );
  });

  test('dashboard card: Cluster replication status', async function (assert) {
    await visit('/vault/usage-reporting');
    await waitFor('[data-test-dashboard-cluster-replication]');

    const card = document.querySelector('[data-test-dashboard-cluster-replication]');
    assert.ok(card, 'renders Cluster replication status card');

    const title = card.querySelector('[data-test-dashboard-card-title]');
    assert.ok(title, 'title is present');
    assert.strictEqual(title.textContent.trim(), 'Cluster replication status', 'title is correct');

    const link = card.querySelector('[data-test-dashboard-card-title-link]');
    assert.ok(link, 'title link is present');
    assert.strictEqual(link.getAttribute('href'), 'replication', 'link points to replication');

    const desc = card.querySelector('[data-test-dashboard-card-description]');
    assert.ok(desc, 'description is present');
    assert.strictEqual(
      desc.textContent.trim(),
      'Check the status and health of Vault clusters',
      'description is correct'
    );
  });

  //
});
