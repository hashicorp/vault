/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';

import personaAnalyticsValue from 'vault/utils/persona-analytics-value';

module('Unit | Utils | persona-analytics-value', function () {
  test('returns predefined personas unchanged', function (assert) {
    assert.strictEqual(personaAnalyticsValue('developer'), 'developer');
    assert.strictEqual(personaAnalyticsValue('platform-engineer'), 'platform-engineer');
    assert.strictEqual(personaAnalyticsValue('security-analyst'), 'security-analyst');
  });

  test('ignores customRole for non-other personas', function (assert) {
    assert.strictEqual(personaAnalyticsValue('developer', 'Something Else'), 'developer');
  });

  test('returns "other" when no custom role is provided', function (assert) {
    assert.strictEqual(personaAnalyticsValue('other'), 'other');
    assert.strictEqual(personaAnalyticsValue('other', ''), 'other');
  });

  test('returns "other" when the custom role normalises to nothing', function (assert) {
    assert.strictEqual(personaAnalyticsValue('other', '   '), 'other');
  });

  test('appends the kebab-cased custom role to "other"', function (assert) {
    assert.strictEqual(personaAnalyticsValue('other', 'DevOps SRE'), 'other-devops-sre');
    assert.strictEqual(personaAnalyticsValue('other', '  Platform   Engineer  '), 'other-platform-engineer');
  });
});
