import Ember from 'ember';

export default Ember.Test.registerAsyncHelper('authLogin', function() {
  visit('/vault/auth?with=token');
  fillIn('[data-test-token]', 'root');
  click('[data-test-auth-submit]');
  // get rid of the root warning flash
  if (find('[data-test-flash-message-body]').length) {
    click('[data-test-flash-message-body]');
  }
});
