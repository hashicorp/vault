import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import backendListPage from 'vault/tests/pages/secrets/backends';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

moduleForAcceptance('Acceptance | settings', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
});

test('settings', function(assert) {
  const now = new Date().getTime();
  const type = 'consul';
  const path = `path-${now}`;

  // mount unsupported backend
  visit('/vault/settings/mount-secret-backend');
  andThen(function() {
    assert.equal(currentURL(), '/vault/settings/mount-secret-backend');

    mountSecrets
      .selectType(type)
      .next()
      .path(path)
      .toggleOptions()
      .defaultTTLVal(100)
      .defaultTTLUnit('s')
      .submit();
  });

  andThen(() => {
    assert.equal(currentURL(), `/vault/secrets`, 'redirects to secrets page');
    assert.ok(
      find('[data-test-flash-message]').text().trim(),
      `Successfully mounted '${type}' at '${path}'!`
    );
    let row = backendListPage.rows().findByPath(path);
    row.menu();
  });

  andThen(() => {
    backendListPage.configLink();
  });

  andThen(() => {
    assert.ok(currentURL(), '/vault/secrets/${path}/configuration', 'navigates to the config page');
  });
});
