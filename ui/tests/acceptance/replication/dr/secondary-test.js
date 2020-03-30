import { module, test } from 'qunit';
import { later } from '@ember/runloop';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, visit, pauseTest } from '@ember/test-helpers';

import authPage from 'vault/tests/pages/auth';

module('Acceptance | DR secondary details', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('visiting DR Recovery Details page', async function(assert) {
    await visit('/vault/replication-dr-promote/details');
    assert.equal(currentURL(), '/vault/replication-dr-promote/details');
  });
});
