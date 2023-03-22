import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consolePanel = create(consoleClass);

module('Acceptance | kv | breadcrumbs', function (hooks) {
  setupApplicationTest(hooks);

  test('it should route back to parent path from metadata tab', async function (assert) {
    await authPage.login();
    await consolePanel.runCommands(['delete sys/mounts/kv', 'write sys/mounts/kv type=kv-v2']);
    await visit('/vault/secrets/kv/list');
    await click('[data-test-secret-create]');
    await fillIn('[data-test-secret-path]', 'foo/bar');
    await click('[data-test-secret-save]');
    await click('[data-test-secret-metadata-tab]');
    await click('[data-test-secret-breadcrumb="foo"]');
    assert.strictEqual(
      currentURL(),
      '/vault/secrets/kv/list/foo/',
      'Routes back to list view on breadcrumb click'
    );
  });
});
