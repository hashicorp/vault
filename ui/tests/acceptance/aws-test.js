import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';

moduleForAcceptance('Acceptance | aws secret backend', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
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
test('aws backend', function(assert) {
  const now = new Date().getTime();
  const path = `aws-${now}`;
  const roleName = 'awsrole';

  mountSupportedSecretBackend(assert, 'aws', path);
  click('[data-test-secret-backend-configure]');
  andThen(() => {
    assert.equal(currentURL(), `/vault/settings/secrets/configure/${path}`);
    assert.ok(find('[data-test-aws-root-creds-form]').length, 'renders the empty root creds form');
    assert.ok(find('[data-test-aws-link="root-creds"]').length, 'renders the root creds link');
    assert.ok(find('[data-test-aws-link="leases"]').length, 'renders the leases config link');
  });

  fillIn('[data-test-aws-input="accessKey"]', 'foo');
  fillIn('[data-test-aws-input="secretKey"]', 'bar');
  click('[data-test-aws-input="root-save"]');
  andThen(() => {
    assert.ok(
      find('[data-test-flash-message]').text().trim(),
      `The backend configuration saved successfully!`
    );
    click('[data-test-flash-message]');
  });
  click('[data-test-aws-link="leases"]');
  click('[data-test-aws-input="lease-save"]');

  andThen(() => {
    assert.ok(
      find('[data-test-flash-message]').text().trim(),
      `The backend configuration saved successfully!`
    );
    click('[data-test-flash-message]');
  });

  click('[data-test-backend-view-link]');
  //back at the roles list
  andThen(() => {
    assert.equal(currentURL(), `/vault/secrets/${path}/list`, `navigates to the roles list`);
  });

  click('[ data-test-secret-create]');
  andThen(() => {
    assert.ok(find('[data-test-secret-header]').text().includes('AWS Role'), `aws: renders the create page`);
  });

  fillIn('[data-test-input="name"]', roleName);
  andThen(function() {
    find('.CodeMirror').get(0).CodeMirror.setValue(JSON.stringify(POLICY));
  });

  // save the role
  click('[data-test-role-aws-create]');
  andThen(() => {
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      `$aws: navigates to the show page on creation`
    );
  });

  click('[data-test-secret-root-link]');

  //back at the roles list
  andThen(() => {
    assert.equal(currentURL(), `/vault/secrets/${path}/list`);
    assert.ok(find(`[data-test-secret-link="${roleName}"]`).length, `aws: role shows in the list`);
  });

  //and delete
  click(`[data-test-secret-link="${roleName}"] [data-test-popup-menu-trigger]`);
  andThen(() => {
    click(`[data-test-aws-role-delete="${roleName}"] button`);
  });
  click(`[data-test-confirm-button]`);

  andThen(() => {
    assert.dom(`[data-test-secret-link="${roleName}"]`).doesNotExist(`aws: role is no longer in the list`);
  });
});
