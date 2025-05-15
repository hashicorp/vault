/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit, currentURL } from '@ember/test-helpers';

module('Acceptance | config-ui/login-settings', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('Login settings index page', async function (assert) {
    assert.expect(1);

    // Mock API response
    this.server.get('/v1/sys/config/ui/login/default-auth', () => {
      return {
        data: {
          key_info: 'key_info',
          keys: ['key1', 'key2'],
        },
      };
    });

    // Visit the login settings index page
    await visit('/login-settings');

    // Check if the page is rendered correctly
    assert.strictEqual(currentURL(), '/login-settings', 'Navigated to login settings index page');
  });

  // keep a render test

  // fetched login rules list should display on index page

  // clicking a rule should navigate to the details page (verify route & render)

  // from login list index page, clicking delete should have a popup modal, should be removed from list on confirm
});
