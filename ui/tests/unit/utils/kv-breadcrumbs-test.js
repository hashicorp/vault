/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { breadcrumbsForSecret, pathIsDirectory, pathIsFromDirectory } from 'kv/utils/kv-breadcrumbs';

module('Unit | Utility | kv-breadcrumbs', function () {
  test('pathIsDirectory works', function (assert) {
    assert.expect(5);
    [
      { path: 'some/path', expect: false },
      { path: 'some/path/', expect: true },
      { path: 'some', expect: false },
      { path: 'some/', expect: true },
      { path: '/some', expect: false },
    ].forEach((scenario) => {
      assert.strictEqual(
        pathIsDirectory(scenario.path),
        scenario.expect,
        `correct for path ${scenario.path}`
      );
    });
  });
  test('pathIsFromDirectory works', function (assert) {
    assert.expect(5);
    [
      { path: 'some/path', expect: true },
      { path: 'some/path/', expect: true },
      { path: 'some', expect: false },
      { path: 'some/', expect: true },
      { path: '/some', expect: true },
    ].forEach((scenario) => {
      assert.strictEqual(
        pathIsFromDirectory(scenario.path),
        scenario.expect,
        `correct for path ${scenario.path}`
      );
    });
  });

  test('breadcrumbsForSecret works', function (assert) {
    let results = breadcrumbsForSecret('kv-mount', 'beep/bop/boop');
    assert.deepEqual(
      results,
      [
        { label: 'beep', route: 'list-directory', models: ['kv-mount', 'beep/'] },
        { label: 'bop', route: 'list-directory', models: ['kv-mount', 'beep/bop/'] },
        { label: 'boop', route: 'secret.details', models: ['kv-mount', 'beep/bop/boop'] },
      ],
      'correct when full nested path to secret'
    );

    results = breadcrumbsForSecret('kv-mount', 'beep/bop/boop', true);
    assert.deepEqual(
      results,
      [
        { label: 'beep', route: 'list-directory', models: ['kv-mount', 'beep/'] },
        { label: 'bop', route: 'list-directory', models: ['kv-mount', 'beep/bop/'] },
        { label: 'boop' },
      ],
      'correct when full nested path to secret and last item current'
    );

    results = breadcrumbsForSecret('kv-mount', 'beep');
    assert.deepEqual(
      results,
      [{ label: 'beep', route: 'secret.details', models: ['kv-mount', 'beep'] }],
      'correct when non-nested secret path'
    );

    results = breadcrumbsForSecret('kv-mount', 'beep', true);
    assert.deepEqual(
      results,
      [{ label: 'beep' }],
      'correct when non-nested secret path and last item current'
    );

    results = breadcrumbsForSecret('kv-mount', 'beep/bop/');
    assert.deepEqual(
      results,
      [
        { label: 'beep', route: 'list-directory', models: ['kv-mount', 'beep/'] },
        { label: 'bop', route: 'list-directory', models: ['kv-mount', 'beep/bop/'] },
      ],
      'correct when path is directory'
    );

    results = breadcrumbsForSecret('kv-mount', 'beep/bop/', true);
    assert.deepEqual(
      results,
      [{ label: 'beep', route: 'list-directory', models: ['kv-mount', 'beep/'] }, { label: 'bop' }],
      'correct when path is directory and last item current'
    );

    results = breadcrumbsForSecret();
    assert.deepEqual(results, [], 'fails gracefully if backend is undefined');

    results = breadcrumbsForSecret('backend-only');
    assert.deepEqual(results, [], 'fails gracefully if secretPath is undefined');
  });
});
