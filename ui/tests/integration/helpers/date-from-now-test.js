/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { subMinutes } from 'date-fns';
import { setupRenderingTest } from 'ember-qunit';
import { dateFromNow } from '../../../helpers/date-from-now';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | date-from-now', function (hooks) {
  setupRenderingTest(hooks);

  test('it works', function (assert) {
    const result = dateFromNow([1481022124443]);
    assert.strictEqual(typeof result, 'string', 'it is a string');
  });

  test('you can include a suffix', function (assert) {
    const result = dateFromNow([1481022124443], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });

  test('you can pass in UTC timestamp', function (assert) {
    const result = dateFromNow(['Fri, 11 Oct 2019 18:56:08 GMT'], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });

  test('you can pass in ISO timestamp', function (assert) {
    const result = dateFromNow(['2019-10-11T18:56:08.984Z'], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });

  test('you can include a suffix using date class', function (assert) {
    const now = Date.now();
    const pastDate = subMinutes(now, 30);
    const result = dateFromNow([pastDate], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });

  test('you can include a suffix using ISO 8601 format', function (assert) {
    const result = dateFromNow(['2021-02-05T20:43:09+00:00'], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });

  test('you can include a suffix in the helper', async function (assert) {
    await render(hbs`<p data-test-date-from-now>Date: {{date-from-now 1481022124443 addSuffix=true}}</p>`);
    assert.dom('[data-test-date-from-now]').includesText(' years ago');
  });
});
