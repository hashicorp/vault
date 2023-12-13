/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentRouteName } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';

module('Acceptance | chroot-namespace enterprise ui', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'chrootNamespace';
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should render normally when chroot namespace exists', async function (assert) {
    await authPage.login();
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'goes to dashboard page');
    assert.dom('[data-test-badge-namespace]').includesText('root', 'Shows root namespace badge');
  });
});
