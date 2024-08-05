/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, waitFor, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

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
    assert.dom('[data-test-not-found]').exists('shows page-error');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });

  test('it should show a public key after saving default configuration', async function (assert) {
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await click(SES.configTab);
    await click(SES.configure);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${sshPath}/configuration/edit`,
      'transitions to the configuration page'
    );
    assert.dom(SES.ssh.configureForm).exists('renders ssh configuration form');

    // default has generate CA checked so we just submit the form
    await click(SES.ssh.sshInput('configure-submit'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${sshPath}/configuration/edit`,
      'stays on configuration form page.'
    );

    await waitFor(SES.ssh.sshInput('public-key'));
    assert.dom(SES.ssh.sshInput('public-key')).exists('renders the public key input on form page');
    assert.dom(SES.ssh.sshInput('public-key')).hasClass('masked-input', 'public key is masked');

    await click(SES.viewBackend);
    await click(SES.configTab);
    assert
      .dom(`[data-test-value-div="Public key"] [data-test-masked-input]`)
      .hasText('***********', 'value for Public key is on config details and is masked');
    assert
      .dom(GENERAL.infoRowValue('Generate signing key'))
      .hasText('Yes', 'value for Generate signing key displays default of true/yes.');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
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
