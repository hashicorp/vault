/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { newLineSplit } from 'core/helpers/new-line-split';

module('Integration | Helper | new-line-split', function (hooks) {
  setupRenderingTest(hooks);

  test('it splits the string by new line characters', async function (assert) {
    const lines = newLineSplit(['First new line.\nSecond new line.\nThird new line.']);
    assert.deepEqual(lines, ['First new line.', 'Second new line.', 'Third new line.']);
  });
});
