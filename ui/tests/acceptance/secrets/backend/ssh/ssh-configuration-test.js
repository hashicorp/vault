/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCmd } from 'vault/tests/helpers/commands';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

module('Acceptance | ssh | configuration', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it should show a public key after saving default configuration', async function (assert) {
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await click(SES.configTab);
    await click(SES.configure);
    assert.strictEqual(
      currentURL(),
      `/vault/settings/secrets/configure/${sshPath}`,
      'transitions to the configuration page'
    );
    assert.dom(SES.ssh.configureForm).exists('renders ssh configuration form');

    // default has generate CA checked so we just submit the form
    await click(SES.ssh.sshInput('configure-submit'));
    assert.strictEqual(
      currentURL(),
      `/vault/settings/secrets/configure/${sshPath}`,
      'stays on configuration form page.'
    );

    await waitFor(SES.ssh.sshInput('public-key'));
    assert.dom(SES.ssh.sshInput('public-key')).exists('renders the public key input on form page');
    assert.dom(SES.ssh.sshInput('public-key')).hasClass('masked-input', 'public key is masked');
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });
});
