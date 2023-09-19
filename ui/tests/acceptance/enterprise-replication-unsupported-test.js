/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import authPage from 'vault/tests/pages/auth';
import { visit } from '@ember/test-helpers';

module('Acceptance | Enterprise | replication unsupported', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.get('/sys/replication/status', function () {
      return {
        data: {
          mode: 'unsupported',
        },
      };
    });
    return authPage.login();
  });

  test('replication page when unsupported', async function (assert) {
    await visit('/vault/replication');
    assert
      .dom('[data-test-replication-title]')
      .hasText('Replication unsupported', 'it shows the unsupported view');
  });
});
