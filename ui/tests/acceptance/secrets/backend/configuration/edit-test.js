/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { runCmd } from 'vault/tests/helpers/commands';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { create } from 'ember-cli-page-object';
import fm from 'vault/tests/pages/components/flash-message';
const flashMessage = create(fm);

module('Acceptance | secrets configuration | edit', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it configures ssh ca', async function (assert) {
    const path = `ssh-configure-${this.uid}`;
    await enablePage.enable('ssh', path);
    await click(SES.configTab);
    await click(SES.configure);
    assert
      .dom(SES.ssh.sshInput('generate-signing-key-checkbox'))
      .isChecked('generate_signing_key defaults to true');
    await click(SES.ssh.sshInput('generate-signing-key-checkbox'));
    await click(SES.ssh.sshInput('configure-submit'));
    assert.strictEqual(
      flashMessage.latestMessage,
      'missing public_key',
      'renders warning flash message for failed save'
    );
    await click(SES.ssh.sshInput('generate-signing-key-checkbox'));
    await click(SES.ssh.sshInput('configure-submit'));
    assert.dom(SES.ssh.sshInput('public-key')).exists('renders public key after saving config');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });
});
