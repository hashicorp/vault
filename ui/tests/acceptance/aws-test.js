import { click, fillIn, findAll, currentURL, find, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | aws secret backend', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  const POLICY = {
    Version: '2012-10-17',
    Statement: [
      {
        Effect: 'Allow',
        Action: 'iam:*',
        Resource: '*',
      },
    ],
  };
  test('aws backend', async function(assert) {
    const now = new Date().getTime();
    const path = `aws-${now}`;
    const roleName = 'awsrole';

    await enablePage.enable('aws', path);

    await click('[data-test-secret-backend-configure]');
    assert.equal(currentURL(), `/vault/settings/secrets/configure/${path}`);
    assert.ok(findAll('[data-test-aws-root-creds-form]').length, 'renders the empty root creds form');
    assert.ok(findAll('[data-test-aws-link="root-creds"]').length, 'renders the root creds link');
    assert.ok(findAll('[data-test-aws-link="leases"]').length, 'renders the leases config link');

    await fillIn('[data-test-aws-input="accessKey"]', 'foo');
    await fillIn('[data-test-aws-input="secretKey"]', 'bar');

    await withFlash(click('[data-test-aws-input="root-save"]'), () => {
      assert.ok(
        find('[data-test-flash-message]').textContent.trim(),
        `The backend configuration saved successfully!`
      );
    });

    await click('[data-test-aws-link="leases"]');
    await withFlash(click('[data-test-aws-input="lease-save"]'), () => {
      assert.ok(
        find('[data-test-flash-message]').textContent.trim(),
        `The backend configuration saved successfully!`
      );
    });

    await click('[data-test-backend-view-link]');
    assert.equal(currentURL(), `/vault/secrets/${path}/list`, `navigates to the roles list`);

    await click('[ data-test-secret-create]');
    assert.ok(
      find('[data-test-secret-header]').textContent.includes('AWS Role'),
      `aws: renders the create page`
    );

    await fillIn('[data-test-input="name"]', roleName);
    findAll('.CodeMirror')[0].CodeMirror.setValue(JSON.stringify(POLICY));

    // save the role
    await click('[data-test-role-aws-create]');
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      `$aws: navigates to the show page on creation`
    );

    await click('[data-test-secret-root-link]');

    assert.equal(currentURL(), `/vault/secrets/${path}/list`);
    assert.ok(findAll(`[data-test-secret-link="${roleName}"]`).length, `aws: role shows in the list`);

    //and delete
    await click(`[data-test-secret-link="${roleName}"] [data-test-popup-menu-trigger]`);
    await click(`[data-test-aws-role-delete="${roleName}"] button`);

    await withFlash(click(`[data-test-confirm-button]`));
    await settled();
    assert.dom(`[data-test-secret-link="${roleName}"]`).doesNotExist(`aws: role is no longer in the list`);
  });
});
