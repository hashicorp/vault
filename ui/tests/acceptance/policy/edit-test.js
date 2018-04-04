import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/policy/edit';

moduleForAcceptance('Acceptance | policy/acl/:name/edit', {
  beforeEach() {
    return authLogin();
  },
});

test('it redirects to list if navigating to root', function(assert) {
  page.visit({ type: 'acl', name: 'root' });
  andThen(function() {
    assert.equal(currentURL(), '/vault/policies/acl', 'navigation to root show redirects you to policy list');
  });
});

test('it does not show delete for default policy', function(assert) {
  page.visit({ type: 'acl', name: 'default' });
  andThen(function() {
    assert.notOk(page.deleteIsPresent, 'there is no delete button');
  });
});

test('it navigates to show when the toggle is clicked', function(assert) {
  page.visit({ type: 'acl', name: 'default' }).toggleEdit();
  andThen(() => {
    assert.equal(currentURL(), '/vault/policy/acl/default', 'toggle navigates from edit to show');
  });
});
