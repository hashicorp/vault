/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL, visit, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { configUrl } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | ssh | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it should prompt configuration after mounting ssh engine', async function (assert) {
    const sshPath = `ssh-${this.uid}`;
    // in this test go through the full mount process. Bypass this step in later tests.
    await visit('/vault/settings/mount-secret-backend');
    await mountBackend('ssh', sshPath);
    await click(SES.configTab);
    assert.dom(GENERAL.emptyStateTitle).hasText('SSH not configured');
    assert.dom(GENERAL.emptyStateActions).hasText('Configure SSH');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it should show error if old url is entered', async function (assert) {
    // we are intentionally not redirecting from the old url to the new one
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await click(SES.configTab);
    await visit(`/vault/settings/secrets/configure/${sshPath}`);
    assert.dom(GENERAL.notFound).exists('shows page-error');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it should show a public key after saving default configuration and allows you to delete public key', async function (assert) {
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await click(SES.configTab);
    await click(SES.configure);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${sshPath}/configuration/edit`,
      'transitions to the configuration page'
    );
    // default has generate CA checked so we just submit the form
    await click(SES.ssh.save);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${sshPath}/configuration`,
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
    await click(SES.ssh.delete);
    assert.dom(GENERAL.confirmMessage).hasText('Confirming will remove the CA certificate information.');
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${sshPath}/configuration/edit`,
      'after deleting public key stays on edit page'
    );
    assert.dom(GENERAL.inputByAttr('privateKey')).hasNoText('Private key is empty and reset');
    assert.dom(GENERAL.inputByAttr('publicKey')).hasNoText('Public key is empty and reset');
    assert.dom(GENERAL.inputByAttr('generateSigningKey')).isChecked('Generate signing key is checked');
    await click(SES.viewBackend);
    await click(SES.configTab);
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('SSH not configured', 'after deleting public key SSH is no longer configured');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it displays error if generate Signing key is not checked and no public and private keys', async function (assert) {
    const path = `ssh-configure-${this.uid}`;
    await enablePage.enable('ssh', path);
    await click(SES.configTab);
    await click(SES.configure);
    assert.dom(GENERAL.inputByAttr('generateSigningKey')).isChecked('generate_signing_key defaults to true');
    await click(GENERAL.inputByAttr('generateSigningKey'));
    await click(SES.ssh.save);
    assert
      .dom(GENERAL.inlineError)
      .hasText('Provide a Public and Private key or set "Generate Signing Key" to true.');
    // visit the details page and confirm the public key is not shown
    await visit(`/vault/secrets/${path}/configuration`);
    assert.dom(GENERAL.infoRowLabel('Public key')).doesNotExist('Public key label does not exist');
    assert.dom(GENERAL.emptyStateTitle).hasText('SSH not configured', 'SSH not configured');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should show API error when SSH configuration read fails', async function (assert) {
    assert.expect(1);
    const path = `ssh-${this.uid}`;
    const type = 'ssh';
    await enablePage.enable(type, path);
    // interrupt get and return API error
    this.server.get(configUrl(type, path), () => {
      return overrideResponse(400, { errors: ['bad request'] });
    });
    await click(SES.configTab);
    assert.dom(SES.error.title).hasText('Error', 'shows the secrets backend error route');
  });
});
