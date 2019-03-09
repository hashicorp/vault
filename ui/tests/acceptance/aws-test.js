import { click, fillIn, findAll, currentURL, find, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import confirmAction from 'vault/tests/pages/components/confirm-action';
import { create } from 'ember-cli-page-object';

const popup = create(confirmAction);

module('Acceptance | aws secret backend', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    await authPage.login();
  });

  hooks.afterEach(async function() {
    await logout.visit();
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
    await settled();
    assert.equal(currentURL(), `/vault/settings/secrets/configure/${path}`);
    assert.ok(findAll('[data-test-aws-root-creds-form]').length, 'renders the empty root creds form');
    assert.ok(findAll('[data-test-aws-link="root-creds"]').length, 'renders the root creds link');
    assert.ok(findAll('[data-test-aws-link="leases"]').length, 'renders the leases config link');

    await fillIn('[data-test-aws-input="accessKey"]', 'foo');
    await fillIn('[data-test-aws-input="secretKey"]', 'bar');

    await click('[data-test-aws-input="root-save"]');
    assert.ok(
      find('[data-test-flash-message]').textContent.trim(),
      `The backend configuration saved successfully!`
    );

    await click('[data-test-aws-link="leases"]');
    await click('[data-test-aws-input="lease-save"]');
    assert.ok(
      find('[data-test-flash-message]').textContent.trim(),
      `The backend configuration saved successfully!`
    );

    await click('[data-test-backend-view-link]');
    await settled();
    assert.equal(currentURL(), `/vault/secrets/${path}/list`, `navigates to the roles list`);

    await click('[ data-test-secret-create]');
    await settled();
    assert.ok(
      find('[data-test-secret-header]').textContent.includes('AWS Role'),
      `aws: renders the create page`
    );

    await fillIn('[data-test-input="name"]', roleName);
    findAll('.CodeMirror')[0].CodeMirror.setValue(JSON.stringify(POLICY));

    // save the role
    await click('[data-test-role-aws-create]');
    // wait for redirect
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      `$aws: navigates to the show page on creation`
    );

    await click('[data-test-secret-root-link]');
    await settled();

    assert.equal(currentURL(), `/vault/secrets/${path}/list`);
    assert.ok(findAll(`[data-test-secret-link="${roleName}"]`).length, `aws: role shows in the list`);

    //and delete
    await click(`[data-test-secret-link="${roleName}"] [data-test-popup-menu-trigger]`);
    // wait for permissions checks
    await settled();
    await popup.delete();
    await popup.confirmDelete();

    //wait for redirect
    await settled();
    assert.dom(`[data-test-secret-link="${roleName}"]`).doesNotExist(`aws: role is no longer in the list`);
  });
});
