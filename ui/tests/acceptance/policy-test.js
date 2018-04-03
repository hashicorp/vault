import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';

moduleForAcceptance('Acceptance | policies', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
});

test('it redirects to acls with unknown policy type', function(assert) {
  visit('/vault/policies/foo');
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.policies.index');
    assert.equal(currentURL(), '/vault/policies/acl');
  });
});
