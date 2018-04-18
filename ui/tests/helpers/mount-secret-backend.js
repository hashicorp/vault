import Ember from 'ember';

export default Ember.Test.registerAsyncHelper('mountSupportedSecretBackend', function(_, assert, type, path) {
  visit('/vault/settings/mount-secret-backend');
  andThen(function() {
    assert.equal(currentURL(), '/vault/settings/mount-secret-backend');
  });

  fillIn('[data-test-secret-backend-type]', type);
  fillIn('[data-test-secret-backend-path]', path);
  click('[data-test-secret-backend-submit]');
  return andThen(() => {
    assert.equal(currentURL(), `/vault/secrets/${path}/list`, `redirects to ${path} index`);
    assert.ok(
      find('[data-test-flash-message]').text().trim(),
      `Successfully mounted '${type}' at '${path}'!`
    );
    click('[data-test-flash-message]');
  });
});
