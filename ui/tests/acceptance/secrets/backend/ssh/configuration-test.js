/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, waitUntil, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import authPage from 'vault/tests/pages/auth';
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
    const flash = this.owner.lookup('service:flash-messages');
    this.store = this.owner.lookup('service:store');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it should prompt configuration after mounting ssh engine', async function (assert) {
    const sshPath = `ssh-${this.uid}`;
    // in this test go through the full mount process. Bypass this step in later tests.
    await visit('/vault/settings/mount-secret-backend');
    await click(SES.mountType('ssh'));
    await fillIn(GENERAL.inputByAttr('path'), sshPath);
    await click(SES.mountSubmit);
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
      'after configuring it navigates to the details page'
    );

    // There is a delay in the backend for the public key to be generated, wait for it to complete and transition to configuration index route
    await waitUntil(() => currentURL() === `/vault/secrets/${sshPath}/configuration`, { timeout: 2000 });
    assert.dom(GENERAL.infoRowLabel('Public key')).exists('Public Key label exists');
    assert.dom(GENERAL.infoRowValue('Public key')).hasText('***********');
    assert
      .dom(GENERAL.infoRowValue('Generate signing key'))
      .hasText('Yes', 'value for Generate signing key displays default of true/yes.');
    // now confirm configure page shows public key and not the config create form
    await click(SES.configure);
    assert.dom(SES.ssh.editConfigSection).exists('renders the edit section');
    // delete Public key
    await click(SES.ssh.deletePublicKey);
    assert
      .dom('[data-test-confirm-action-message]')
      .hasText('This will remove the CA certificate information.');
    await click(GENERAL.confirmButton);
    // There is a delay in the backend for the public key to be generated, wait for it to complete and transition to configuration index route
    await waitUntil(() => currentURL() === `/vault/secrets/${sshPath}/configuration`, { timeout: 2000 });
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('SSH not configured', 'after deleting public key SSH is no longer configured');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it should throw validation errors if generate Signing key is not checked and no public and private keys', async function (assert) {
    const path = `ssh-configure-${this.uid}`;
    await enablePage.enable('ssh', path);
    await click(SES.configTab);
    await click(SES.configure);
    assert
      .dom(GENERAL.inputByAttr('generate-signing-key-checkbox'))
      .isChecked('generate_signing_key defaults to true');
    await click(GENERAL.inputByAttr('generate-signing-key-checkbox'));
    await click(SES.ssh.save);
    assert.true(this.flashDangerSpy.calledWith('missing public_key'), 'Danger flash message is displayed');
    // visit the details page and confirm the public key is not shown
    await visit(`/vault/secrets/${path}/configuration`);
    assert.dom(GENERAL.infoRowLabel('Public key')).doesNotExist('Public Key label does not exist');
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
