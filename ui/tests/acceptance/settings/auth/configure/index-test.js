import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/index';

moduleForAcceptance('Acceptance | settings/auth/configure', {
  beforeEach() {
    return authLogin();
  },
});

test('it redirects to section options when there are no other sections', function(assert) {
  const path = `approle-${new Date().getTime()}`;
  const type = 'approle';
  enablePage.visit().enableAuth(type, path);
  page.visit({ path });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.equal(currentURL(), `/vault/settings/auth/configure/${path}/options`, 'loads the options route');
  });
});

test('it redirects to the first section', function(assert) {
  const path = `aws-${new Date().getTime()}`;
  const type = 'aws';
  enablePage.visit().enableAuth(type, path);
  page.visit({ path });

  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.equal(
      currentURL(),
      `/vault/settings/auth/configure/${path}/client`,
      'loads the first section for the type of auth method'
    );
  });
});
