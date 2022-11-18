import { currentRouteName, settled, visit, waitUntil } from '@ember/test-helpers';
import { module, test, skip } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/pki/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/pki/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/pki/create', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  skip('it creates a role and redirects', async function (assert) {
    const path = `pki-${new Date().getTime()}`;
    await enablePage.enable('pki', path);
    await settled();
    await editPage.visitRoot({ backend: path });
    await settled();
    await editPage.createRole('role', 'example.com');
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.dom('[data-test-edit-link="true"]').exists('shows the edit button');
    assert.dom('[data-test-credentials-link="true"]').exists('shows the generate button');
    assert.dom('[data-test-sign-link="true"]').exists('shows the sign button');

    await showPage.visit({ backend: path, id: 'role' });
    await settled();
    await showPage.generateCert();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await showPage.visit({ backend: path, id: 'role' });
    await settled();
    await visit(`/vault/secrets/${path}/credentials/role?action=sign`);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await listPage.visitRoot({ backend: path });
    await settled();
    assert.strictEqual(listPage.secrets.length, 1, 'shows role in the list');
    const secret = listPage.secrets.objectAt(0);
    await secret.menuToggle();
    await settled();
    assert.ok(listPage.menuItems.length > 0, 'shows links in the menu');
  });

  test('it deletes a role', async function (assert) {
    const path = `pki-${new Date().getTime()}`;
    await enablePage.enable('pki', path);
    await settled();
    await editPage.visitRoot({ backend: path });
    await settled();
    await editPage.createRole('role', 'example.com');
    await settled();
    await showPage.visit({ backend: path, id: 'role' });
    await settled();
    await showPage.deleteRole();
    await settled();
    assert.ok(
      await waitUntil(() => currentRouteName() === 'vault.cluster.secrets.backend.list-root'),
      'redirects to list page'
    );
    assert.ok(listPage.backendIsEmpty, 'no roles listed');
  });
});
