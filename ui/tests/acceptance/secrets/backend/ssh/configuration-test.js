/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, currentURL, visit, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { configUrl } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | ssh | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it should show error if old url is entered', async function (assert) {
    // we are intentionally not redirecting from the old url to the new one
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await visit(`/vault/settings/secrets/configure/${sshPath}`);
    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(`Sorry, we were unable to find any content at settings/secrets/configure/${sshPath}.`);
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining('Double check the URL or return to the dashboard. Go to dashboard');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it should show a public key after saving default configuration and allows you to delete public key', async function (assert) {
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${sshPath}/configuration/edit`,
      'transitions to the configuration page'
    );
    // default has generate CA checked so we just submit the form
    await click(GENERAL.submitButton);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${sshPath}/configuration/plugin-settings`,
      'navigates to the details page.'
    );
    // There is a delay in the backend for the public key to be generated, wait for it to complete by checking that the public key is displayed
    await waitFor(GENERAL.infoRowLabel('Public key'));
    assert.dom(GENERAL.infoRowLabel('Public key')).exists('public key shown on the details screen');

    await click(SES.configure);
    assert
      .dom(SES.ssh.editConfigSection)
      .exists('renders the edit configuration section of the form and not the create part');
    // delete Public key
    await click(GENERAL.button('delete-public-key'));
    assert.dom(GENERAL.confirmMessage).hasText('Confirming will remove the CA certificate information.');
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${sshPath}/configuration/edit`,
      'after deleting public key stays on edit page'
    );
    assert.dom(GENERAL.inputByAttr('private_key')).hasNoText('Private key is empty and reset');
    assert.dom(GENERAL.inputByAttr('public_key')).hasNoText('Public key is empty and reset');
    assert.dom(GENERAL.inputByAttr('generate_signing_key')).isChecked('Generate signing key is checked');
    await click(SES.viewBackend);
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));

    await click(GENERAL.tabLink('general-settings'));
    await click(GENERAL.tabLink('plugin-settings'));

    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${sshPath}/configuration/edit`,
      'should redirect to edit page since public key is no longer configured'
    );
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it displays error if generate Signing key is not checked and no public and private keys', async function (assert) {
    const path = `ssh-configure-${this.uid}`;
    await enablePage.enable('ssh', path);
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    assert
      .dom(GENERAL.inputByAttr('generate_signing_key'))
      .isChecked('generate_signing_key defaults to true');
    await click(GENERAL.inputByAttr('generate_signing_key'));
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('generate_signing_key'))
      .hasText('Provide a Public and Private key or set "Generate Signing Key" to true.');
    // visit the details page and confirm the public key is not shown
    await visit(`/vault/secrets-engines/${path}/configuration/general-settings`);
    await click(GENERAL.tabLink('plugin-settings'));
    assert.dom(GENERAL.infoRowLabel('Public key')).doesNotExist('Public key label does not exist');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should show API error when SSH configuration read fails', async function (assert) {
    const path = `ssh-${this.uid}`;
    const type = 'ssh';
    await enablePage.enable(type, path);
    // interrupt get and return API error
    this.server.get(configUrl(type, path), () => {
      return overrideResponse(400, { errors: ['bad request'] });
    });
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.error',
      'it redirects to the secrets backend error route'
    );
    assert.dom(GENERAL.pageError.title(400)).hasText('ERROR 400 Error');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('A problem has occurred. Check the Vault logs or console for more details.');
    assert.dom(GENERAL.pageError.details).hasText('bad request');
  });
});
