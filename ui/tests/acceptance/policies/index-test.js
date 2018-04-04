import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/policies/index';

moduleForAcceptance('Acceptance | policies/acl', {
  beforeEach() {
    return authLogin();
  },
});

test('it lists default and root acls', function(assert) {
  page.visit({ type: 'acl' });
  andThen(() => {
    let policies = page.policies();
    assert.equal(currentURL(), '/vault/policies/acl');
    assert.ok(policies.findByName('root'), 'root policy shown in the list');
    assert.ok(policies.findByName('default'), 'default policy shown in the list');
  });
});

test('it navigates to show when clicking on the link', function(assert) {
  page.visit({ type: 'acl' });
  andThen(() => {
    page.policies().findByName('default').click();
  });

  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.policy.show');
    assert.equal(currentURL(), '/vault/policy/acl/default');
  });
});
