import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/policy/show';

moduleForAcceptance('Acceptance | policy/acl/:name', {
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

test('it navigates to edit when the toggle is clicked', function(assert) {
  page.visit({ type: 'acl', name: 'default' }).toggleEdit();
  andThen(() => {
    assert.equal(currentURL(), '/vault/policy/acl/default/edit', 'toggle navigates to edit page');
  });
});
