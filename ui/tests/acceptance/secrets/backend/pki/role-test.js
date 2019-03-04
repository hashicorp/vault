import { settled, waitFor, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/pki/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/pki/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/pki/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  const mountAndNav = async function() {
    const path = `pki-${new Date().getTime()}`;
    await enablePage.enable('pki', path);
    await editPage.visitRoot({ backend: path });
    await settled();
    return path;
  };

  test('it creates a role and redirects', async function(assert) {
    const path = await mountAndNav(assert);
    await editPage.createRole('role', 'example.com');
    await settled();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    await settled();
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    assert.ok(showPage.generateCertIsPresent, 'shows the generate button');
    assert.ok(showPage.signCertIsPresent, 'shows the sign button');

    await showPage.visit({ backend: path, id: 'role' });
    await showPage.generateCert();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await showPage.visit({ backend: path, id: 'role' });
    await waitFor('[data-test-sign-link]');
    await showPage.signCert();
    await settled();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await listPage.visitRoot({ backend: path });
    assert.equal(listPage.secrets.length, 1, 'shows role in the list');
    let secret = listPage.secrets.objectAt(0);
    await secret.menuToggle();
    // wait for permission check
    await settled();
    assert.ok(listPage.menuItems.length > 0, 'shows links in the menu');
  });

  test('it deletes a role', async function(assert) {
    await mountAndNav(assert);
    await editPage.createRole('role', 'example.com');
    // wait for permission check to render edit button
    await settled();
    await showPage.edit();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.edit', 'navs to the edit page');

    await waitFor('[data-test-role-delete]');
    await editPage.deleteRole();
    // wait for redirect
    await settled();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'redirects to list page');
    assert.ok(listPage.backendIsEmpty, 'no roles listed');
  });
});
