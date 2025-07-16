/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | console/log error', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    const errorText = 'Error deleting at: sys/foo. URL: v1/sys/foo Code: 404';
    this.set('content', errorText);
    await render(hbs`{{console/log-error content=this.content}}`);
    assert.dom('pre').includesText(errorText);
  });
});
