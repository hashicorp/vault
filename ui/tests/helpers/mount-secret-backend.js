import Ember from 'ember';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

export default Ember.Test.registerAsyncHelper('mountSupportedSecretBackend', function(_, assert, type, path) {
  mountSecrets.visit();
  andThen(() => {
    return mountSecrets.mount(type, path);
  });
  return andThen(() => {
    assert.equal(currentURL(), `/vault/secrets/${path}/list`, `redirects to ${path} index`);
    assert.ok(
      find('[data-test-flash-message]').text().trim(),
      `Successfully mounted '${type}' at '${path}'!`
    );
    click('[data-test-flash-message]');
  });
});
