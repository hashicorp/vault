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
    let results = breadcrumbsForSecret('beep/bop/boop');
    assert.deepEqual(
      results,
      [
        { label: 'beep', route: 'list-directory', model: 'beep/' },
        { label: 'bop', route: 'list-directory', model: 'beep/bop/' },
        { label: 'boop', route: 'secret.details', model: 'beep/bop/boop' },
      ],
      'correct when full nested path to secret'
    );

    results = breadcrumbsForSecret('beep/bop/boop', true);
    assert.deepEqual(
      results,
      [
        { label: 'beep', route: 'list-directory', model: 'beep/' },
        { label: 'bop', route: 'list-directory', model: 'beep/bop/' },
        { label: 'boop' },
      ],
      'correct when full nested path to secret and last item current'
    );

    results = breadcrumbsForSecret('beep');
    assert.deepEqual(
      results,
      [{ label: 'beep', route: 'secret.details', model: 'beep' }],
      'correct when non-nested secret path'
    );

    results = breadcrumbsForSecret('beep', true);
    assert.deepEqual(
      results,
      [{ label: 'beep' }],
      'correct when non-nested secret path and last item current'
    );

    results = breadcrumbsForSecret('beep/bop/');
    assert.deepEqual(
      results,
      [
        { label: 'beep', route: 'list-directory', model: 'beep/' },
        { label: 'bop', route: 'list-directory', model: 'beep/bop/' },
      ],
      'correct when path is directory'
    );

    results = breadcrumbsForSecret('beep/bop/', true);
    assert.deepEqual(
      results,
      [{ label: 'beep', route: 'list-directory', model: 'beep/' }, { label: 'bop' }],
      'correct when path is directory and last item current'
    );

    results = breadcrumbsForSecret();
    assert.deepEqual(results, [], 'fails gracefully if secretPath is undefined');
  });
});
