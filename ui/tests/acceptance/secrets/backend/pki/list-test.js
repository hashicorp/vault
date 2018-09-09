import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/secrets/backend/list';

module('Acceptance | secrets/pki/list', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  const mountAndNav = assert => {
    const path = `pki-${new Date().getTime()}`;
    mountSupportedSecretBackend(assert, 'pki', path);
    page.visitRoot({ backend: path });
  };

  test('it renders an empty list', function(assert) {
    mountAndNav(assert);
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'redirects from the index');
    assert.ok(page.createIsPresent, 'create button is present');
    assert.ok(page.configureIsPresent, 'configure button is present');
    assert.equal(page.tabs.length, 2, 'shows 2 tabs');
    assert.ok(page.backendIsEmpty);
  });

  test('it navigates to the create page', function(assert) {
    mountAndNav(assert);
    page.create();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.create-root', 'links to the create page');
  });

  test('it navigates to the configure page', function(assert) {
    mountAndNav(assert);
    page.configure();
    assert.equal(
      currentRouteName(),
      'vault.cluster.settings.configure-secret-backend.section',
      'links to the configure page'
    );
  });
});
