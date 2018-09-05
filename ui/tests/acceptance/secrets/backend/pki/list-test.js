import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/secrets/backend/list';

moduleForAcceptance('Acceptance | secrets/pki/list', {
  beforeEach() {
    return authLogin();
  },
});

const mountAndNav = assert => {
  const path = `pki-${new Date().getTime()}`;
  mountSupportedSecretBackend(assert, 'pki', path);
  page.visitRoot({ backend: path });
};

test('it renders an empty list', function(assert) {
  mountAndNav(assert);
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'redirects from the index');
    assert.ok(page.createIsPresent, 'create button is present');
    assert.ok(page.configureIsPresent, 'configure button is present');
    assert.equal(page.tabs.length, 2, 'shows 2 tabs');
    assert.ok(page.backendIsEmpty);
  });
});

test('it navigates to the create page', function(assert) {
  mountAndNav(assert);
  page.create();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.create-root', 'links to the create page');
  });
});

test('it navigates to the configure page', function(assert) {
  mountAndNav(assert);
  page.configure();
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.settings.configure-secret-backend.section',
      'links to the configure page'
    );
  });
});
