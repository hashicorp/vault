/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { formatCli } from 'core/utils/code-generators/cli';

module('Integration | Util | code-generators/cli', function (hooks) {
  setupTest(hooks);

  test('formatCli: it formats CLI command with content', async function (assert) {
    const content = `- <<EOT
  path "secret/*" {
    capabilities = ["read"]
  }
EOT`;
    const formatted = formatCli({ command: 'policy write my-policy', content: content });
    const expected = `vault policy write my-policy ${content}`;
    assert.strictEqual(formatted, expected, 'it formats CLI command with content');
  });

  test('formatCli: it handles empty content', async function (assert) {
    const formatted = formatCli({ command: 'policy list', content: '' });
    const expected = 'vault policy list';
    assert.strictEqual(formatted, expected, 'it formats CLI command with empty content');
  });
});
