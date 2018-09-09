import { registerAsyncHelper } from '@ember/test';

export default registerAsyncHelper('authLogin', function(app, token) {
  visit('/vault/auth?with=token');
  fillIn('[data-test-token]', token || 'root');
  click('[data-test-auth-submit]');
  // get rid of the root warning flash
  if (find('[data-test-flash-message-body]').length) {
    return click('[data-test-flash-message-body]');
  }
});
