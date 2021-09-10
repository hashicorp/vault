import { click, fillIn, settled, visit, triggerKeyEvent } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/auth/enable';

module('Acceptance | userpass secret backend', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  test('userpass backend', async function(assert) {
    let n = Math.random();
    const path1 = `userpass-${++n}`;
    const path2 = `userpass-${++n}`;
    const user1 = 'user1';
    const user2 = 'user2';

    // enable the first userpass method with one username
    await enablePage.enable('userpass', path1);
    await settled();
    await click('[data-test-save-config="true"]');
    await settled();
    await visit(`/vault/access/${path1}/item/user/create`);
    await settled();
    await fillIn('[data-test-input="username"]', user1);
    await triggerKeyEvent('[data-test-input="username"]', 'keyup', 65);
    await fillIn('[data-test-textarea]', user1);
    await triggerKeyEvent('[data-test-textarea]', 'keyup', 65);
    await click('[data-test-save-config="true"]');
    await settled();

    // enable the first userpass method with one username
    await visit(`/vault/settings/auth/enable`);
    await settled();
    await click('[data-test-mount-type="userpass"]');
    await settled();
    await click('[data-test-mount-next]');
    await settled();
    await fillIn('[data-test-input="path"]', path2);
    await settled();
    await click('[data-test-mount-submit="true"]');
    await settled();
    await click('[data-test-save-config="true"]');
    await settled();
    await click(`[data-test-auth-backend-link="${path2}"]`);
    await settled();
    await click('[data-test-create="user"]');
    await settled();
    await fillIn('[data-test-input="username"]', user2);
    await triggerKeyEvent('[data-test-input="username"]', 'keyup', 65);
    await fillIn('[data-test-textarea]', user2);
    await triggerKeyEvent('[data-test-textarea]', 'keyup', 65);
    await click('[data-test-save-config="true"]');
    await settled();

    //confirming that the user was created.  There was a bug where the apiPath was not being updated when toggling between auth routes
    assert
      .dom('[data-test-list-item-content]')
      .hasText(user2, 'user just created shows in current auth list');

    //confirm that the auth method 1 shows the user1.  There was a bug where it was not updated the list when toggling between auth routes
    await visit(`/vault/access/${path1}/item/user`);
    await settled();
    assert
      .dom('[data-test-list-item-content]')
      .hasText(user1, 'first user created shows in current auth list');
  });
});
