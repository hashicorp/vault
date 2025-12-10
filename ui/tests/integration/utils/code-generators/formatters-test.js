/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { formatEot } from 'core/utils/code-generators/formatters';

module('Integration | Util | code-generators/formatters', function (hooks) {
  setupTest(hooks);

  test('formatEot: it wraps content in EOT heredoc block', async function (assert) {
    const content = 'path "secret/*" {\n  capabilities = ["read"]\n}';
    const formatted = formatEot(content);
    const expected = `<<EOT
path "secret/*" {
  capabilities = ["read"]
}
EOT`;
    assert.strictEqual(formatted, expected, 'it wraps content with line breaks "\n"');
  });

  test('formatEot: it handles single line content', async function (assert) {
    const formatted = formatEot('single line content');
    const expected = `<<EOT
single line content
EOT`;
    assert.strictEqual(formatted, expected, 'it wraps single line');
  });

  test('formatEot: it handles empty content', async function (assert) {
    const formatted = formatEot('');
    const expected = `<<EOT

EOT`;
    assert.strictEqual(formatted, expected, 'it wraps empty content');
  });

  test('formatEot: it handles multi-line content', async function (assert) {
    const content = `line 1
line 2
line 3`;
    const formatted = formatEot(content);
    const expected = `<<EOT
line 1
line 2
line 3
EOT`;
    assert.strictEqual(formatted, expected, 'it wraps multi-line content');
  });
});
