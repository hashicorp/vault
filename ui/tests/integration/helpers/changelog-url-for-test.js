/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { changelogUrlFor } from '../../../helpers/changelog-url-for';

const CHANGELOG_URL = 'https://www.github.com/hashicorp/vault/blob/main/CHANGELOG.md#';

module('Integration | Helper | changelog-url-for', function (hooks) {
  setupRenderingTest(hooks);

  test('it builds an enterprise URL', function (assert) {
    const result = changelogUrlFor(['1.5.0+prem']);
    assert.strictEqual(result, CHANGELOG_URL.concat('150'));
  });

  test('it builds an OSS URL', function (assert) {
    const result = changelogUrlFor(['1.4.3']);
    assert.strictEqual(result, CHANGELOG_URL.concat('143'));
  });

  test('it returns the base changelog URL if the version is less than 1.4.3', function (assert) {
    const result = changelogUrlFor(['1.4.0']);
    assert.strictEqual(result, CHANGELOG_URL);
  });

  test('it returns the base changelog URL if version cannot be found', function (assert) {
    const result = changelogUrlFor(['']);
    assert.strictEqual(result, CHANGELOG_URL);
  });
});
