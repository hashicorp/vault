/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | console/log json', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.codeMirror = this.owner.lookup('service:code-mirror');
    // TODO: Fix JSONEditor/CodeMirror
    setRunOptions({
      rules: {
        label: { enabled: false },
      },
    });
  });

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.on('myAction', function(val) { ... });
    const objectContent = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    const expectedText = JSON.stringify(objectContent, null, 2);

    this.set('content', objectContent);

    await render(hbs`{{console/log-json content=this.content}}`);
    const instance = find('[data-test-component=code-mirror-modifier]').innerText;
    assert.strictEqual(instance, expectedText);
  });
});
