import { registerAsyncHelper } from '@ember/test';

export default registerAsyncHelper('authLogout', function() {
  visit('/vault/logout');
});
