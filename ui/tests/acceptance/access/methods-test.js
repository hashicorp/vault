import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/access/methods';

moduleForAcceptance('Acceptance | /access/', {
  beforeEach() {
    return authLogin();
  },
});

test('it navigates', function(assert) {
  page.visit();
  andThen(() => {
    assert.ok(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.ok(page.navLinks(0).isActive, 'the first link is active');
    assert.equal(page.navLinks(0).text, 'Auth Methods');
  });
});
