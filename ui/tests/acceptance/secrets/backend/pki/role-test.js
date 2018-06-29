import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import editPage from 'vault/tests/pages/secrets/backend/pki/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/pki/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

moduleForAcceptance('Acceptance | secrets/pki/create', {
  beforeEach() {
    return authLogin();
  },
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
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    assert.ok(showPage.generateCertIsPresent, 'shows the generate button');
    assert.ok(showPage.signCertIsPresent, 'shows the sign button');
  });

  showPage.visit({ backend: path, id: 'role' });
  showPage.generateCert();
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );
  });

  showPage.visit({ backend: path, id: 'role' });
  showPage.signCert();
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );
  });

  listPage.visitRoot({ backend: path });
  andThen(() => {
    assert.equal(listPage.secrets.length, 1, 'shows role in the list');
  });
});

test('it deletes a role', function(assert) {
  mountAndNav(assert);
  editPage.createRole('role', 'example.com');
  showPage.edit();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.edit', 'navs to the edit page');
  });

  editPage.deleteRole();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'redirects to list page');
    assert.ok(listPage.backendIsEmpty, 'no roles listed');
  });
});
