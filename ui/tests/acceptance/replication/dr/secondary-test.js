import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { pauseTest, currentURL, visit } from '@ember/test-helpers';

import authPage from 'vault/tests/pages/auth';

module('Acceptance | DR secondary details', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it visits the Secondary Disaster Recovery Details page', async function(assert) {
    await visit('/vault/replication-dr-promote/details');
    assert.equal(currentURL(), '/vault/replication-dr-promote/details');
  });

  test('I should see the error message if I am not a DR secondary cluster', async function(assert) {
    // not passing any data, so it should always show empty state
    await visit('/vault/replication-dr-promote/details');
    assert.dom('[data-test-component="empty-state"]').exists();
  });
});
