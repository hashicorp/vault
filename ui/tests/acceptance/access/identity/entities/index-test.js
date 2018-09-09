import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/index';

module('Acceptance | /access/identity/entities', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it renders the entities page', function(assert) {
    page.visit({ item_type: 'entities' });
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.index', 'navigates to the correct route');
  });

  test('it renders the groups page', function(assert) {
    page.visit({ item_type: 'groups' });
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.index', 'navigates to the correct route');
  });
});
