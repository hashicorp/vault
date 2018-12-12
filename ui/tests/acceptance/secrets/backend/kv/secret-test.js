import { currentURL, currentRouteName } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import withFlash from 'vault/tests/helpers/with-flash';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

let writeSecret = async function(backend, path, key, val) {
  await listPage.visitRoot({ backend });
  await listPage.create();
  return editPage.createSecret(path, key, val);
};

module('Acceptance | secrets/secret/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    this.server = apiStub({ usePassthrough: true });
    await logout.visit();
    return authPage.login();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('it creates a secret and redirects', async function(assert) {
    const path = `kv-path-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'secret' });
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');

    await listPage.create();
    assert.ok(editPage.hasMetadataFields, 'shows the metadata form');
    await editPage.createSecret(path, 'foo', 'bar');

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  test('it can create a secret when check-and-set is required', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = 'foo/bar';
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await consoleComponent.runCommands(`write ${enginePath}/config cas_required=true`);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  test('version 1 performs the correct capabilities lookup', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = 'foo/bar';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await withFlash(
      mountSecrets
        .next()
        .path(enginePath)
        .version(1)
        .submit()
    );

    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  test('it redirects to the path ending in / for list pages', async function(assert) {
    const path = `foo/bar/kv-path-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'secret' });
    await listPage.create();
    await editPage.createSecret(path, 'foo', 'bar');
    await listPage.visit({ backend: 'secret', id: 'foo/bar' });
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list');
    assert.ok(currentURL().endsWith('/'), 'redirects to the path ending in a slash');
  });

  test('it can edit via the JSON input', async function(assert) {
    let content = JSON.stringify({ foo: 'fa', bar: 'boo' });
    const path = `kv-path-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'secret' });
    await listPage.create();
    await editPage.path(path).toggleJSON();
    await editPage.editor.fillIn(this, content);
    await editPage.save();

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    assert.equal(
      showPage.editor.content(this),
      JSON.stringify({ bar: 'boo', foo: 'fa' }, null, 2),
      'saves the content'
    );
  });

  test('version 2 with restricted policy still allows creation', async function(assert) {
    let backend = 'kv-v2';
    const V2_POLICY = `'
      path "kv-v2/metadata/*" {
        capabilities = ["list"]
      }
      path "kv-v2/data/secret" {
        capabilities = ["create", "read", "update"]
      }
    '`;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=kv options=version=2`,
      `write sys/policies/acl/kv-v2-degrade policy=${V2_POLICY}`,
      // delete any kv previously written here so that tests can be re-run
      'delete kv-v2/metadata/secret',
      'write -field=client_token auth/token/create policies=kv-v2-degrade',
    ]);

    let userToken = consoleComponent.lastLogOutput;
    await logout.visit();
    await authPage.login(userToken);

    await writeSecret(backend, 'secret', 'foo', 'bar');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    await logout.visit();
  });

  test('version 2 with restricted policy still allows edit', async function(assert) {
    let backend = 'kv-v2';
    const V2_POLICY = `'
      path "kv-v2/metadata/*" {
        capabilities = ["list"]
      }
      path "kv-v2/data/secret" {
        capabilities = ["create", "read", "update"]
      }
    '`;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=kv options=version=2`,
      `write sys/policies/acl/kv-v2-degrade policy=${V2_POLICY}`,
      // delete any kv previously written here so that tests can be re-run
      'delete kv-v2/metadata/secret',
      'write -field=client_token auth/token/create policies=kv-v2-degrade',
    ]);

    let userToken = consoleComponent.lastLogOutput;
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await editPage.visitEdit({ backend, id: 'secret' });
    assert.notOk(editPage.hasMetadataFields, 'hides the metadata form');
    await editPage.editSecret('bar', 'baz');

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    await logout.visit();
  });
});
