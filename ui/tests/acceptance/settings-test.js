import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';

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
  });

  fillIn('[data-test-secret-backend-type]', type);
  fillIn('[data-test-secret-backend-path]', path);
  click('[data-test-secret-backend-options]');

  // set a ttl of 100s
  fillIn('[data-test-secret-backend-default-ttl] input', 100);
  fillIn('[data-test-secret-backend-default-ttl] select', 's');

  click('[data-test-secret-backend-submit]');
  andThen(() => {
    assert.equal(currentURL(), `/vault/secrets`, 'redirects to secrets page');
    assert.ok(
      find('[data-test-flash-message]').text().trim(),
      `Successfully mounted '${type}' at '${path}'!`
    );
  });

  //show mount details
  click(`[data-test-secret-backend-row="${path}"] [data-test-secret-backend-detail]`);
  andThen(() => {
    assert.ok(
      find('[data-test-secret-backend-details="default-ttl"]').text().match(/100/),
      'displays the input ttl'
    );
  });
});
