import { registerAsyncHelper } from '@ember/test';
import { visit, fillIn, find, click } from '@ember/test-helpers';

export async function login(token) {
  await visit('/vault/auth?with=token');
  fillIn('[data-test-token]', token || 'root');
  await click('[data-test-auth-submit]');
  // get rid of the root warning flash
  if (find('[data-test-flash-message-body]')) {
    await click('[data-test-flash-message-body]');
  }
}

registerAsyncHelper('authLogin', function(app, token) {
  login(token);
});
