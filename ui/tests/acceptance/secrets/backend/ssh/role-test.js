import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/ssh/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/ssh/show';
import generatePage from 'vault/tests/pages/secrets/backend/ssh/generate-otp';
import listPage from 'vault/tests/pages/secrets/backend/list';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/ssh', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  const mountAndNav = async () => {
    const path = `ssh-${new Date().getTime()}`;
    await enablePage.enable('ssh', path);
    await editPage.visitRoot({ backend: path });
    return path;
  };

  test('it creates a role and redirects', async function(assert) {
    const path = await mountAndNav(assert);
    await editPage.createOTPRole('role');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.generateIsPresent, 'shows the generate button');

    await showPage.visit({ backend: path, id: 'role' });
    await showPage.generate();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await listPage.visitRoot({ backend: path });
    assert.equal(listPage.secrets.length, 1, 'shows role in the list');
    let secret = listPage.secrets.objectAt(0);
    await secret.menuToggle();
    assert.ok(listPage.menuItems.length > 0, 'shows links in the menu');
  });

  test('it deletes a role', async function(assert) {
    const path = await mountAndNav(assert);
    await editPage.createOTPRole('role');
    await showPage.visit({ backend: path, id: 'role' });
    await showPage.deleteRole();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'redirects to list page');
    assert.ok(listPage.backendIsEmpty, 'no roles listed');
  });

  test('it generates an OTP', async function(assert) {
    const path = await mountAndNav(assert);
    await editPage.createOTPRole('role');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.generateIsPresent, 'shows the generate button');

    await showPage.visit({ backend: path, id: 'role' });
    await showPage.generate();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await generatePage.generateOTP();
    assert.ok(generatePage.warningIsPresent, 'shows warning');
    await generatePage.back();
    assert.ok(generatePage.userIsPresent, 'clears generate, shows user input');
    assert.ok(generatePage.ipIsPresent, 'clears generate, shows ip input');
  });
});
