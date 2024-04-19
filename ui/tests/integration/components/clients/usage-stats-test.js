/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/usage-stats', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(
      hbs`
        <Clients::UsageStats @title="My stats">
         yielded content!
        </Clients::UsageStats>`
    );

    assert.dom('[data-test-usage-stats="My stats"]').exists();
    assert
      .dom('[data-test-usage-stats="My stats"]')
      .hasTextContaining('yielded content!', 'it renders yielded content');
    assert
      .dom('a')
      .hasAttribute('href', 'https://developer.hashicorp.com/vault/tutorials/monitoring/usage-metrics');
  });
});
