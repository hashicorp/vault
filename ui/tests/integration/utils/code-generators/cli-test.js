/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { cliTemplate } from 'core/utils/code-generators/cli';

module('Integration | Util | code-generators/cli', function (hooks) {
  setupTest(hooks);

  test('cliTemplate: it formats CLI command with content', async function (assert) {
    const content = `- <<EOT
  path "secret/*" {
    capabilities = ["read"]
  }
EOT`;
    const formatted = cliTemplate({ command: 'policy write my-policy', content });
    const expected = `vault policy write my-policy ${content}`;
    assert.strictEqual(formatted, expected, 'it formats CLI command with content');
  });

  test('cliTemplate: it handles empty content', async function (assert) {
    const formatted = cliTemplate({ command: 'policy list', content: '' });
    const expected = 'vault policy list';
    assert.strictEqual(formatted, expected, 'it formats CLI command with empty content');
  });

  test('cliTemplate: it only supplies placeholder "[args]" when there is no command', async function (assert) {
    let formatted = cliTemplate();
    let expected = 'vault <command> [args]';
    assert.strictEqual(formatted, expected, 'it formats CLI command with undefined args');

    formatted = cliTemplate({ command: 'read my-secret' });
    expected = 'vault read my-secret';
    assert.strictEqual(formatted, expected, 'it formats CLI command without "[args]" placeholder');

    formatted = cliTemplate({ command: '' });
    expected = 'vault <command> [args]';
    assert.strictEqual(
      formatted,
      expected,
      'it formats CLI command with placeholders content is an empty string'
    );
  });
});
