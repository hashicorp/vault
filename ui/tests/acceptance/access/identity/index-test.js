import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/access/identity/index';

moduleForAcceptance('Acceptance | /access/identity/entities', {
  beforeEach() {
    return authLogin();
  },
});

test('it renders the entities page', function(assert) {
  page.visit({ item_type: 'entities' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.index', 'navigates to the correct route');
  });
});

test('it renders the groups page', function(assert) {
  page.visit({ item_type: 'groups' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.index', 'navigates to the correct route');
  });
});
