import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';

const disableReplication = (type, assert) => {
  // disable performance replication
  visit(`/vault/replication/${type}`);
  return andThen(() => {
    if (find('[data-test-replication-link="manage"]').length) {
      click('[data-test-replication-link="manage"]');
      click('[data-test-disable-replication] button');
      click('[data-test-confirm-button]');
      if (assert) {
        andThen(() => {
          assert.equal(currentURL(), `/vault/replication/${type}`, 'redirects to the replication page');
          assert.equal(
            // TODO better test selectors for flash messages
            find('[data-test-flash-message-body]:contains(This cluster is having)').text().trim(),
            'This cluster is having replication disabled. Vault will be unavailable for a brief period and will resume service shortly.',
            'renders info flash when disabled'
          );
          click('[data-test-flash-message-body]:contains(This cluster is having)');
        });
      }
    } else {
      // do nothing, it's already off
    }
  });
};
moduleForAcceptance('Acceptance | Enterprise | replication', {
  beforeEach() {
    authLogin();
    disableReplication('dr');
    return disableReplication('performance');
  },
  afterEach() {
    disableReplication('dr');
    return disableReplication('performance');
  },
});

test('replication', function(assert) {
  const secondaryName = 'firstSecondary';
  const mode = 'blacklist';
  const mountType = 'kv';
  let mountPath;

  visit('/vault/replication');
  andThen(function() {
    assert.equal(currentURL(), '/vault/replication');
  });

  // enable perf replication
  click('[data-test-replication-type-select="performance"]');
  fillIn('[data-test-replication-cluster-mode-select]', 'primary');

  click('[data-test-replication-enable]');
  andThen(() => {
    pollCluster();
  });

  // add a secondary with a mount filter config
  click('[data-test-replication-link="secondaries"]');
  click('[data-test-secondary-add]');
  fillIn('[data-test-replication-secondary-id]', secondaryName);
  //expand the config
  click('[data-test-replication-secondary-token-options]');
  fillIn('[data-test-replication-filter-mount-mode]', mode);
  click(`[data-test-mount-filter="${mountType}"]:eq(0)`);
  andThen(() => {
    mountPath = find(`[data-test-mount-filter-path-for-type="${mountType}"]`).first().text().trim();
  });
  click('[data-test-secondary-add]');

  // fetch new secondaries
  andThen(() => {
    pollCluster();
  });

  // click into the added secondary's mount filter config
  click('[data-test-replication-link="secondaries"]');
  click('[data-test-popup-menu-trigger]');

  click('[data-test-replication-mount-filter-link]');
  andThen(() => {
    assert.equal(currentURL(), `/vault/replication/performance/secondaries/config/show/${secondaryName}`);
    assert.ok(
      find('[data-test-mount-config-mode]').text().trim().toLowerCase().includes(mode),
      'show page renders the correct mode'
    );
    assert
      .dom('[data-test-mount-config-paths]')
      .hasText(mountPath, 'show page renders the correct mount path');
  });
  // click edit

  // delete config
  click('[data-test-replication-link="edit-mount-config"]');
  click('[data-test-delete-mount-config] button');
  click('[data-test-confirm-button]');
  andThen(() => {
    assert.equal(
      currentURL(),
      `/vault/replication/performance/secondaries`,
      'redirects to the secondaries page'
    );
    assert.equal(
      find('[data-test-flash-message-body]:contains(The performance mount filter)').text().trim(),
      `The performance mount filter config for the secondary ${secondaryName} was successfully deleted.`,
      'renders success flash upon deletion'
    );
    click('[data-test-flash-message-body]:contains(The performance mount filter)');
  });

  // nav to DR
  visit('/vault/replication/dr');
  fillIn('[data-test-replication-cluster-mode-select]', 'secondary');
  andThen(() => {
    assert.ok(
      find('[data-test-replication-enable]').is(':disabled'),
      'dr secondary enable is disabled when other replication modes are on'
    );
  });

  // disable performance replication
  disableReplication('replication', assert);

  // enable dr replication
  visit('/vault/replication/dr');
  fillIn('[data-test-replication-cluster-mode-select]', 'primary');
  click('button[type="submit"]');
  andThen(() => {
    pollCluster();
  });
  andThen(() => {
    assert.ok(
      find('[data-test-replication-title]').text().includes('Disaster Recovery'),
      'it displays the replication type correctly'
    );
    assert.ok(
      find('[data-test-replication-mode-display]').text().includes('primary'),
      'it displays the cluster mode correctly'
    );
  });

  // add dr secondary
  click('[data-test-replication-link="secondaries"]');
  click('[data-test-secondary-add]');
  fillIn('[data-test-replication-secondary-id]', secondaryName);
  click('[data-test-secondary-add]');
  andThen(() => {
    pollCluster();
  });
  click('[data-test-replication-link="secondaries"]');
  andThen(() => {
    assert
      .dom('[data-test-secondary-name]')
      .hasText(secondaryName, 'it displays the secondary in the list of known secondaries');
  });

  // disable dr replication
  disableReplication('dr', assert);
  return wait();
});

test('disabling dr primary when perf replication is enabled', function(assert) {
  visit('/vault/replication/performance');
  // enable perf replication
  fillIn('[data-test-replication-cluster-mode-select]', 'primary');

  click('[data-test-replication-enable]');
  andThen(() => {
    pollCluster();
  });

  // enable dr replication
  visit('/vault/replication/dr');
  fillIn('[data-test-replication-cluster-mode-select]', 'primary');
  click('[data-test-replication-enable]');
  andThen(() => {
    pollCluster();
  });
  visit('/vault/replication/dr/manage');
  andThen(() => {
    assert.ok(find('[data-test-demote-warning]').length, 'displays the demotion warning');
  });

  // disable replication
  disableReplication('performance', assert);
  disableReplication('dr', assert);
  return wait();
});
