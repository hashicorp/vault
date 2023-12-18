/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | transform-advanced-templating', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render', async function (assert) {
    setRunOptions({
      rules: {
        // TODO: fix JSONEditor/CodeMirror
        label: { enabled: false },
      },
    });
    this.model = {
      pattern: '(\\d{3})-(\\d{2})-(?<last>\\d{5})',
      encodeFormat: null,
      decodeFormats: {},
    };
    await render(hbs`<TransformAdvancedTemplating @model={{this.model}} />`);

    assert.dom('.box').doesNotExist('Form is hidden when not toggled');
    await click('[data-test-toggle-advanced]');
    await fillIn('[data-test-input="regex-test-val"]', '123-45-67890');
    [
      ['$1', '123'],
      ['$2', '45'],
      ['$3', '67890'],
      ['$last', '67890'],
    ].forEach(([p, v]) => {
      assert.dom(`[data-test-regex-group-position="${p}"]`).hasText(`${p}`, `Capture group ${p} renders`);
      assert.dom(`[data-test-regex-group-value="${p}"]`).hasText(v, `Capture group value ${v} renders`);
    });
    // need to simulate InputEvent
    await triggerEvent('[data-test-encode-format] input', 'input', { data: '$' });
    const options = this.element.querySelectorAll('.autocomplete-input-option');
    ['$1: 123', '$2: 45', '$3: 67890', '$last: 67890'].forEach((val, index) => {
      assert.dom(options[index]).hasText(val, 'Autocomplete option renders');
    });

    assert.dom('[data-test-kv-object-editor]').exists('KvObjectEditor renders for decode formats');
    assert.dom('[data-test-decode-format]').exists('AutocompleteInput renders for decode format value');
    await fillIn('[data-test-kv-key]', 'last');
    await fillIn('[data-test-decode-format] input', '$last');
    assert.deepEqual(this.model.decodeFormats, { last: '$last' }, 'Decode formats updates correctly');
  });
});
