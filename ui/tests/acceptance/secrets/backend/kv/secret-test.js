import {
  click,
  visit,
  settled,
  currentURL,
  currentRouteName,
  fillIn,
  triggerKeyEvent,
} from '@ember/test-helpers';
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

  test('it disables save when validation errors occur', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await click('[data-test-secret-create="true"]');
    await fillIn('[data-test-secret-path="true"]', '123');
    await triggerKeyEvent('[data-test-secret-path="true"]', 'keyup', 65);
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'A secret with this path already exists.',
        'when duplicate path it shows correct error message'
      );

    document.querySelector('#maxVersions').value = 'abc';
    await triggerKeyEvent('[data-test-input="maxVersions"]', 'keyup', 65);
    assert
      .dom('[data-test-input="maxVersions"]')
      .hasClass('has-error-border', 'shows border error on input with error');
    assert.dom('[data-test-secret-save="true"]').isDisabled('Save button is disabled');

    await fillIn('[data-test-input="maxVersions"]', 20);
    await triggerKeyEvent('[data-test-input="maxVersions"]', 'keyup', 65);
    await fillIn('[data-test-secret-path="true"]', '1234');
    await triggerKeyEvent('[data-test-secret-path="true"]', 'keyup', 65);
    await fillIn('[data-test-secret-key]', 'someKey');
    await triggerKeyEvent('[data-test-secret-key]', 'keyup', 65);
    await click('[data-test-secret-save="true"]');
    assert.equal(currentURL(), `/vault/secrets/${enginePath}/show/1234`, 'navigates to show secret');
    // using a number for path to check for bug where paths as numbers wouldn't render properly
    assert.dom('[data-test-row-label="someKey"]').exists({ count: 1 }, 'renders the key');
  });

  test('version 1 performs the correct capabilities lookup', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = 'foo/bar';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets
      .next()
      .path(enginePath)
      .toggleOptions()
      .version(1)
      .submit();
    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  // https://github.com/hashicorp/vault/issues/5960
  test('version 1: nested paths creation maintains ability to navigate the tree', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = '1/2/3/4';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets
      .next()
      .path(enginePath)
      .toggleOptions()
      .version(1)
      .submit();
    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');

    // setup an ancestor for when we delete
    await listPage.visitRoot({ backend: enginePath });
    await listPage.secrets.filterBy('text', '1/')[0].click();
    await listPage.create();
    await editPage.createSecret('1/2', 'foo', 'bar');

    // lol we have to do this because ember-cli-page-object doesn't like *'s in visitable
    await listPage.visitRoot({ backend: enginePath });
    await listPage.secrets.filterBy('text', '1/')[0].click();
    await listPage.secrets.filterBy('text', '2/')[0].click();
    await listPage.secrets.filterBy('text', '3/')[0].click();
    await listPage.create();

    await editPage.createSecret(secretPath + 'a', 'foo', 'bar');
    await listPage.visitRoot({ backend: enginePath });
    await listPage.secrets.filterBy('text', '1/')[0].click();
    await listPage.secrets.filterBy('text', '2/')[0].click();
    let secretLink = listPage.secrets.filterBy('text', '3/')[0];
    assert.ok(secretLink, 'link to the 3/ branch displays properly');

    await listPage.secrets.filterBy('text', '3/')[0].click();
    await listPage.secrets.objectAt(0).menuToggle();
    await settled();
    await listPage.delete();
    await listPage.confirmDelete();
    await settled();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list');
    assert.equal(currentURL(), `/vault/secrets/${enginePath}/list/1/2/3/`, 'remains on the page');

    await listPage.secrets.objectAt(0).menuToggle();
    await listPage.delete();
    await listPage.confirmDelete();
    await settled();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list');
    assert.equal(
      currentURL(),
      `/vault/secrets/${enginePath}/list/1/`,
      'navigates to the ancestor created earlier'
    );
  });
  test('first level secrets redirect properly upon deletion', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = 'test';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets
      .next()
      .path(enginePath)
      .toggleOptions()
      .version(1)
      .submit();
    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');
    await showPage.deleteSecretV1();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'redirected to the list page on delete'
    );
  });

  // https://github.com/hashicorp/vault/issues/5994
  test('version 1: key named keys', async function(assert) {
    await consoleComponent.runCommands([
      'vault write sys/mounts/test type=kv',
      'refresh',
      'vault write test/a keys=a keys=b',
    ]);
    await showPage.visit({ backend: 'test', id: 'a' });
    assert.ok(showPage.editIsPresent, 'renders the page properly');
  });

  test('it redirects to the path ending in / for list pages', async function(assert) {
    const path = `foo/bar/kv-path-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'secret' });
    await listPage.create();
    await editPage.createSecret(path, 'foo', 'bar');
    await settled();
    // use visit helper here because ids with / in them get encoded
    await visit('/vault/secrets/secret/list/foo/bar');
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
    const V2_POLICY = `
      path "kv-v2/metadata/*" {
        capabilities = ["list"]
      }
      path "kv-v2/data/secret" {
        capabilities = ["create", "read", "update"]
      }
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=kv options=version=2`,
      `write sys/policies/acl/kv-v2-degrade policy=${btoa(V2_POLICY)}`,
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
  });

  test('version 2 with restricted policy still allows edit', async function(assert) {
    let backend = 'kv-v2';
    const V2_POLICY = `
      path "kv-v2/metadata/*" {
        capabilities = ["list"]
      }
      path "kv-v2/data/secret" {
        capabilities = ["create", "read", "update"]
      }
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=kv options=version=2`,
      `write sys/policies/acl/kv-v2-degrade policy=${btoa(V2_POLICY)}`,
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
  });

  test('version 2 with policy with destroy capabilities shows modal', async function(assert) {
    let backend = 'kv-v2';
    const V2_POLICY = `
      path "kv-v2/destroy/*" {
        capabilities = ["update"]
      }
      path "kv-v2/metadata/*" {
        capabilities = ["list", "update", "delete"]
      }
      path "kv-v2/data/secret" {
        capabilities = ["create", "read", "update"]
      }
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=kv options=version=2`,
      `write sys/policies/acl/kv-v2-degrade policy=${btoa(V2_POLICY)}`,
      // delete any kv previously written here so that tests can be re-run
      'delete kv-v2/metadata/secret',
      'write -field=client_token auth/token/create policies=kv-v2-degrade',
    ]);

    let userToken = consoleComponent.lastLogOutput;
    await logout.visit();
    await authPage.login(userToken);

    await writeSecret(backend, 'secret', 'foo', 'bar');
    await click('[data-test-delete-open-modal]');
    await settled();
    assert.dom('[data-test-delete-modal="destroy-version"]').exists('destroy this version option shows');
    assert.dom('[data-test-delete-modal="destroy-all-versions"]').exists('destroy all versions option shows');
    assert.dom('[data-test-delete-modal="delete-version"]').doesNotExist('delete version does not show');
  });

  test('version 2 with policy with only delete option does not show modal and undelete is an option', async function(assert) {
    let backend = 'kv-v2';
    const V2_POLICY = `
      path "kv-v2/delete/*" {
        capabilities = ["update"]
      }
      path "kv-v2/undelete/*" {
        capabilities = ["update"]
      }
      path "kv-v2/metadata/*" {
        capabilities = ["list","read","create","update"]
      }
      path "kv-v2/data/secret" {
        capabilities = ["create", "read"]
      }
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=kv options=version=2`,
      `write sys/policies/acl/kv-v2-degrade policy=${btoa(V2_POLICY)}`,
      // delete any kv previously written here so that tests can be re-run
      'delete kv-v2/metadata/secret',
      'write -field=client_token auth/token/create policies=kv-v2-degrade',
    ]);

    let userToken = consoleComponent.lastLogOutput;
    await logout.visit();
    await authPage.login(userToken);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    assert.dom('[data-test-delete-open-modal]').doesNotExist('delete version does not show');
    assert.dom('[data-test-secret-v2-delete="true"]').exists('drop down delete shows');
    await showPage.deleteSecretV2();
    // unable to reload page in test scenario so going to list and back to secret to confirm deletion
    let url = `/vault/secrets/${backend}/list`;
    await visit(url);
    await click('[data-test-secret-link="secret"]');
    assert.dom('[data-test-component="empty-state"]').exists('secret has been deleted');
    assert.dom('[data-test-secret-undelete]').exists('undelete button shows');
  });

  test('paths are properly encoded', async function(assert) {
    let backend = 'kv';
    let paths = [
      '(',
      ')',
      '"',
      //"'",
      '!',
      '#',
      '$',
      '&',
      '*',
      '+',
      '@',
      '{',
      '|',
      '}',
      '~',
      '[',
      '\\',
      ']',
      '^',
      '_',
    ].map(char => `${char}some`);
    assert.expect(paths.length * 2);
    let secretName = '2';
    let commands = paths.map(path => `write '${backend}/${path}/${secretName}' 3=4`);
    await consoleComponent.runCommands(['write sys/mounts/kv type=kv', ...commands]);
    for (let path of paths) {
      await listPage.visit({ backend, id: path });
      assert.ok(listPage.secrets.filterBy('text', '2')[0], `${path}: secret is displayed properly`);
      await listPage.secrets.filterBy('text', '2')[0].click();
      assert.equal(
        currentRouteName(),
        'vault.cluster.secrets.backend.show',
        `${path}: show page renders correctly`
      );
    }
  });

  test('create secret with space shows version data', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = 'space space';
    // mount version 2
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets
      .next()
      .path(enginePath)
      .submit();
    await settled();
    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');
    await settled();
    await click('[data-test-popup-menu-trigger="version"]');
    await settled();
    await click('[data-test-version-history]');
    await settled();
    assert.dom('[data-test-list-item-content]').exists('renders the version and not an error state');
    // click on version
    await click('[data-test-popup-menu-trigger="true"]');
    await click('[data-test-version]');
    await settled();
    // perform encode function that should be done by the encodePath
    let encodedSecretPath = secretPath.replace(/ /g, '%20');
    assert.equal(currentURL(), `/vault/secrets/${enginePath}/show/${encodedSecretPath}?version=1`);
  });

  // the web cli does not handle a quote as part of a path, so we test it here via the UI
  test('creating a secret with a single or double quote works properly', async function(assert) {
    await consoleComponent.runCommands('write sys/mounts/kv type=kv');
    let paths = ["'some", '"some'];
    for (let path of paths) {
      await listPage.visitRoot({ backend: 'kv' });
      await listPage.create();
      await editPage.createSecret(`${path}/2`, 'foo', 'bar');
      await listPage.visit({ backend: 'kv', id: path });
      assert.ok(listPage.secrets.filterBy('text', '2')[0], `${path}: secret is displayed properly`);
      await listPage.secrets.filterBy('text', '2')[0].click();
      assert.equal(
        currentRouteName(),
        'vault.cluster.secrets.backend.show',
        `${path}: show page renders correctly`
      );
    }
  });

  test('filter clears on nav', async function(assert) {
    await consoleComponent.runCommands([
      'vault write sys/mounts/test type=kv',
      'refresh',
      'vault write test/filter/foo keys=a keys=b',
      'vault write test/filter/foo1 keys=a keys=b',
      'vault write test/filter/foo2 keys=a keys=b',
    ]);
    await listPage.visit({ backend: 'test', id: 'filter' });
    assert.equal(listPage.secrets.length, 3, 'renders three secrets');
    await listPage.filterInput('filter/foo1');
    assert.equal(listPage.secrets.length, 1, 'renders only one secret');
    await listPage.secrets.objectAt(0).click();
    await showPage.breadcrumbs.filterBy('text', 'filter')[0].click();
    assert.equal(listPage.secrets.length, 3, 'renders three secrets');
    assert.equal(listPage.filterInputValue, 'filter/', 'pageFilter has been reset');
  });

  let setupNoRead = async function(backend, canReadMeta = false) {
    const V2_WRITE_ONLY_POLICY = `
      path "${backend}/+/+" {
        capabilities = ["create", "update", "list"]
      }
      path "${backend}/+" {
        capabilities = ["list"]
      }
    `;

    const V2_WRITE_WITH_META_READ_POLICY = `
      path "${backend}/+/+" {
        capabilities = ["create", "update", "list"]
      }
      path "${backend}/metadata/+" {
        capabilities = ["read"]
      }
      path "${backend}/+" {
        capabilities = ["list"]
      }
    `;
    const V1_WRITE_ONLY_POLICY = `
     path "${backend}/+" {
        capabilities = ["create", "update", "list"]
      }
    `;

    let policy;
    if (backend === 'kv-v2' && canReadMeta) {
      policy = V2_WRITE_WITH_META_READ_POLICY;
    } else if (backend === 'kv-v2') {
      policy = V2_WRITE_ONLY_POLICY;
    } else if (backend === 'kv-v1') {
      policy = V1_WRITE_ONLY_POLICY;
    }
    await consoleComponent.runCommands([
      // disable any kv previously enabled kv
      `delete sys/mounts/${backend}`,
      `write sys/mounts/${backend} type=kv options=version=${backend === 'kv-v2' ? 2 : 1}`,
      `write sys/policies/acl/${backend} policy=${btoa(policy)}`,
      `write -field=client_token auth/token/create policies=${backend}`,
    ]);

    return consoleComponent.lastLogOutput;
  };
  test('write without read: version 2', async function(assert) {
    let backend = 'kv-v2';
    let userToken = await setupNoRead(backend);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await showPage.visit({ backend, id: 'secret' });
    assert.ok(showPage.noReadIsPresent, 'shows no read empty state');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await editPage.visitEdit({ backend, id: 'secret' });
    assert.notOk(editPage.hasMetadataFields, 'hides the metadata form');
    assert.ok(editPage.showsNoCASWarning, 'shows no CAS write warning');

    await editPage.editSecret('bar', 'baz');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
  });

  test('write without read: version 2 with metadata read', async function(assert) {
    let backend = 'kv-v2';
    let userToken = await setupNoRead(backend, true);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await showPage.visit({ backend, id: 'secret' });
    assert.ok(showPage.noReadIsPresent, 'shows no read empty state');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await editPage.visitEdit({ backend, id: 'secret' });
    assert.notOk(editPage.hasMetadataFields, 'hides the metadata form');
    assert.ok(editPage.showsV2WriteWarning, 'shows v2 warning');

    await editPage.editSecret('bar', 'baz');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
  });

  test('write without read: version 1', async function(assert) {
    let backend = 'kv-v1';
    let userToken = await setupNoRead(backend);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await showPage.visit({ backend, id: 'secret' });
    assert.ok(showPage.noReadIsPresent, 'shows no read empty state');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await editPage.visitEdit({ backend, id: 'secret' });
    assert.ok(editPage.showsV1WriteWarning, 'shows v1 warning');

    await editPage.editSecret('bar', 'baz');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
  });
});
