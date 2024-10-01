/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | console/log command', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    const commandText = 'list this/path';
    this.set('content', commandText);

    await render(hbs`<Console::LogCommand @content={{this.content}} />`);
    assert.dom('.console-ui-command').includesText(commandText);
  });
});
