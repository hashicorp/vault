/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { find, fillIn, visit, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

import authPage from 'vault/tests/pages/auth';

module('Acceptance | API Explorer', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it filters paths after swagger-ui is loaded', async function (assert) {
    await visit('/vault/api-explorer');
    await waitUntil(() => {
      return find('[data-test-filter-input]').disabled === false;
    });
    await fillIn('[data-test-filter-input]', 'sys/health');
    assert.dom('.opblock').exists({ count: 1 }, 'renders a single opblock for sys/health');
  });
});
