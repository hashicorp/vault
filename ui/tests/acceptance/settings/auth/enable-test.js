import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/auth/enable';
import listPage from 'vault/tests/pages/access/methods';

moduleForAcceptance('Acceptance | settings/auth/enable', {
  beforeEach() {
    return authLogin();
  },
});

test('it mounts and redirects', function(assert) {
  page.visit();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.enable');
  });
  // always force the new mount to the top of the list
  const path = `approle-${new Date().getTime()}`;
  const type = 'approle';
  page.enableAuth(type, path);
  andThen(() => {
    assert.equal(
      page.flash.latestMessage,
      `Successfully mounted ${type} auth method at ${path}.`,
      'success flash shows'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.methods',
      'redirects to the auth backend list page'
    );
    assert.ok(listPage.backendLinks().findById(path), 'mount is present in the list');
  });
});
