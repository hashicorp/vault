import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';

moduleForAcceptance('Acceptance | policies', {
  beforeEach() {
    return authLogin();
  },
});

test('it redirects to acls on unknown policy type', function(assert) {
  visit('/vault/policy/foo/default');
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.policies.index');
    assert.equal(currentURL(), '/vault/policies/acl');
  });

  visit('/vault/policy/foo/default/edit');
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.policies.index');
    assert.equal(currentURL(), '/vault/policies/acl');
  });
});

test('it redirects to acls on index navigation', function(assert) {
  visit('/vault/policy/acl');
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.policies.index');
    assert.equal(currentURL(), '/vault/policies/acl');
  });
});
