import { click, fillIn, findAll, currentURL, find, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';

const disableReplication = async (type, assert) => {
  // disable performance replication
  await visit(`/vault/replication/${type}`);
  if (findAll('[data-test-replication-link="manage"]').length) {
    await click('[data-test-replication-link="manage"]');
    await click('[data-test-disable-replication] button');
    await click('[data-test-confirm-button]');
    if (assert) {
      assert.equal(currentURL(), `/vault/replication/${type}`, 'redirects to the replication page');
      assert.equal(
        // TODO better test selectors for flash messages
        find('[data-test-flash-message-body]:contains(This cluster is having)').textContent.trim(),
        'This cluster is having replication disabled. Vault will be unavailable for a brief period and will resume service shortly.',
        'renders info flash when disabled'
      );
      await click('[data-test-flash-message-body]:contains(This cluster is having)');
    }
  } else {
    // do nothing, it's already off
  }
};

module('Acceptance | Enterprise | replication', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    await authPage.login();
    await disableReplication('dr');
    return disableReplication('performance');
  });

  hooks.afterEach(async function() {
    await disableReplication('dr');
    return disableReplication('performance');
  });

  test('replication', async function(assert) {
    const secondaryName = 'firstSecondary';
    const mode = 'blacklist';
    const mountType = 'kv';
    let mountPath;

    await visit('/vault/replication');
    assert.equal(currentURL(), '/vault/replication');

    // enable perf replication
    await click('[data-test-replication-type-select="performance"]');
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');

    await click('[data-test-replication-enable]');
    await pollCluster(this.owner);

    // add a secondary with a mount filter config
    await click('[data-test-replication-link="secondaries"]');
    await click('[data-test-secondary-add]');
    await fillIn('[data-test-replication-secondary-id]', secondaryName);
    //expand the config
    await click('[data-test-replication-secondary-token-options]');
    await fillIn('[data-test-replication-filter-mount-mode]', mode);
    await click(`[data-test-mount-filter="${mountType}"]:eq(0)`);
    mountPath = find(`[data-test-mount-filter-path-for-type="${mountType}"]`)
      .first()
      .text()
      .trim();
    await click('[data-test-secondary-add]');

    await pollCluster(this.owner);

    // click into the added secondary's mount filter config
    await click('[data-test-replication-link="secondaries"]');
    await click('[data-test-popup-menu-trigger]');

    await click('[data-test-replication-mount-filter-link]');
    assert.equal(currentURL(), `/vault/replication/performance/secondaries/config/show/${secondaryName}`);
    assert.ok(
      find('[data-test-mount-config-mode]')
        .textContent.trim()
        .toLowerCase()
        .includes(mode),
      'show page renders the correct mode'
    );
    assert
      .dom('[data-test-mount-config-paths]')
      .hasText(mountPath, 'show page renders the correct mount path');
    // click edit

    // delete config
    await click('[data-test-replication-link="edit-mount-config"]');
    await click('[data-test-delete-mount-config] button');
    await click('[data-test-confirm-button]');
    assert.equal(
      currentURL(),
      `/vault/replication/performance/secondaries`,
      'redirects to the secondaries page'
    );
    assert.equal(
      find('[data-test-flash-message-body]:contains(The performance mount filter)').textContent.trim(),
      `The performance mount filter config for the secondary ${secondaryName} was successfully deleted.`,
      'renders success flash upon deletion'
    );
    await click('[data-test-flash-message-body]:contains(The performance mount filter)');

    // nav to DR
    await visit('/vault/replication/dr');
    await fillIn('[data-test-replication-cluster-mode-select]', 'secondary');
    assert.ok(
      find('[data-test-replication-enable]').is(':disabled'),
      'dr secondary enable is disabled when other replication modes are on'
    );

    // disable performance replication
    disableReplication('replication', assert);

    // enable dr replication
    await visit('/vault/replication/dr');
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('button[type="submit"]');
    pollCluster();
    assert.ok(
      find('[data-test-replication-title]').textContent.includes('Disaster Recovery'),
      'it displays the replication type correctly'
    );
    assert.ok(
      find('[data-test-replication-mode-display]').textContent.includes('primary'),
      'it displays the cluster mode correctly'
    );

    // add dr secondary
    await click('[data-test-replication-link="secondaries"]');
    await click('[data-test-secondary-add]');
    await fillIn('[data-test-replication-secondary-id]', secondaryName);
    await click('[data-test-secondary-add]');
    pollCluster();
    await click('[data-test-replication-link="secondaries"]');
    assert
      .dom('[data-test-secondary-name]')
      .hasText(secondaryName, 'it displays the secondary in the list of known secondaries');

    // disable dr replication
    disableReplication('dr', assert);
    return wait();
  });

  test('disabling dr primary when perf replication is enabled', async function(assert) {
    await visit('/vault/replication/performance');
    // enable perf replication
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');

    await click('[data-test-replication-enable]');
    await pollCluster(this.owner);

    // enable dr replication
    await visit('/vault/replication/dr');
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('[data-test-replication-enable]');
    await pollCluster(this.owner);
    await visit('/vault/replication/dr/manage');
    assert.ok(findAll('[data-test-demote-warning]').length, 'displays the demotion warning');

    // disable replication
    disableReplication('performance', assert);
    disableReplication('dr', assert);
    return wait();
  });
});
