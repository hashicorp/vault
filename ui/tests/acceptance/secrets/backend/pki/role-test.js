import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/pki/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/pki/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

module('Acceptance | secrets/pki/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  const mountAndNav = assert => {
    const path = `pki-${new Date().getTime()}`;
    mountSupportedSecretBackend(assert, 'pki', path);
    editPage.visitRoot({ backend: path });
    return path;
  };

  test('it creates a role and redirects', function(assert) {
    const path = mountAndNav(assert);
    editPage.createRole('role', 'example.com');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    assert.ok(showPage.generateCertIsPresent, 'shows the generate button');
    assert.ok(showPage.signCertIsPresent, 'shows the sign button');

    showPage.visit({ backend: path, id: 'role' });
    showPage.generateCert();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    showPage.visit({ backend: path, id: 'role' });
    showPage.signCert();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    listPage.visitRoot({ backend: path });
    assert.equal(listPage.secrets.length, 1, 'shows role in the list');
  });

  test('it deletes a role', function(assert) {
    mountAndNav(assert);
    editPage.createRole('role', 'example.com');
    showPage.edit();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.edit', 'navs to the edit page');

    editPage.deleteRole();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'redirects to list page');
    assert.ok(listPage.backendIsEmpty, 'no roles listed');
  });
});
