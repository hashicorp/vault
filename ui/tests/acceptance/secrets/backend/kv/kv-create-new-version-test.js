import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/kv/kv-run-commands';
import { SELECTORS } from 'vault/tests/helpers/kv/kv-general-selectors';
import { PAGE } from 'vault/tests/helpers/kv/kv-page-selectors';

module('Acceptance | kv | creates a secret and a new version', function (hooks) {
  setupApplicationTest(hooks);
  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup KV engine
    this.mountPath = `kv-engine-${uuidv4()}`;
    await enablePage.enable('kv', this.mountPath);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
  });

  test('it creates a new secret then a new secret version and navigates to details route', async function (assert) {
    assert.expect(9);

    const secretPath = 'my-secret';
    await visit(`/vault/secrets/${this.mountPath}/kv/list`);
    assert.dom(SELECTORS.emptyStateTitle).hasText('No secrets yet');
    assert
      .dom(`${SELECTORS.emptyStateActions} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/kv/create`);
    await click(PAGE.list.createSecret);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/create?initialKey=`);

    await fillIn(PAGE.form.inputByAttr('path'), secretPath);
    await fillIn(PAGE.form.keyInput(), 'foo-1');
    await fillIn(PAGE.form.maskedValueInput(), 'bar-1');
    await click(PAGE.form.secretSave);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${secretPath}/details?version=1`);

    await click(PAGE.details.createNewVersion);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.mountPath}/kv/${secretPath}/details/edit?version=1`
    );
    assert.dom(PAGE.form.inputByAttr('path')).isDisabled();

    await fillIn(PAGE.form.keyInput(1), 'foo-2');
    await fillIn(PAGE.form.maskedValueInput(1), 'bar-2');
    await click(PAGE.form.secretSave);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${secretPath}/details?version=2`);

    await visit(`/vault/secrets/${this.mountPath}/kv/list`);
    await click(PAGE.list.item(secretPath));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.mountPath}/kv/${secretPath}/details?version=2`,
      'list view navigates to latest version'
    );
    assert.dom(SELECTORS.tooltipTrigger).hasTextContaining('Version 2');
  });
});
