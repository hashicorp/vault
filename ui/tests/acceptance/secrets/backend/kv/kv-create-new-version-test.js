import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { deleteEngineCmd, runCmd } from 'vault/tests/helpers/commands';

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
    await runCmd(deleteEngineCmd(this.mountPath), false);
  });

  test('it creates a new secret then a new secret version and navigates to details route', async function (assert) {
    assert.expect(9);

    const secretPath = 'my-secret';
    await visit(`/vault/secrets/${this.mountPath}/kv/list`);
    assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
    assert
      .dom(`${PAGE.emptyStateActions} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/kv/create`);
    await click(PAGE.list.createSecret);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.mountPath}/kv/create?initialKey=`,
      'url is correct'
    );

    await fillIn(FORM.inputByAttr('path'), secretPath);
    await fillIn(FORM.keyInput(), 'foo-1');
    await fillIn(FORM.maskedValueInput(), 'bar-1');
    await click(FORM.saveBtn);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${secretPath}/details?version=1`);

    await click(PAGE.detail.createNewVersion);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.mountPath}/kv/${secretPath}/details/edit?version=1`
    );
    assert.dom(FORM.inputByAttr('path')).isDisabled();

    await fillIn(FORM.keyInput(1), 'foo-2');
    await fillIn(FORM.maskedValueInput(1), 'bar-2');
    await click(FORM.saveBtn);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${secretPath}/details?version=2`);

    await visit(`/vault/secrets/${this.mountPath}/kv/list`);
    await click(PAGE.list.item(secretPath));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.mountPath}/kv/${secretPath}/details?version=2`,
      'list view navigates to latest version'
    );
    assert.dom(PAGE.detail.versionTimestamp).hasTextContaining('Version 2');
  });
});
