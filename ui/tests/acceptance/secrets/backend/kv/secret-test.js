import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

moduleForAcceptance('Acceptance | secrets/secret/create', {
  beforeEach() {
    return authLogin();
  },
});

test('it creates a secret and redirects', function(assert) {
  const path = `kv-${new Date().getTime()}`;
  listPage.visitRoot({ backend: 'secret' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
  });

  listPage.create();
  editPage.createSecret(path, 'foo', 'bar');
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });
});
