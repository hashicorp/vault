/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { _getPathParam, getHelpUrlForModel, pathToHelpUrlSegment } from 'vault/utils/openapi-helpers';

module('Unit | Utility | OpenAPI helper utils', function () {
  test(`pathToHelpUrlSegment`, function (assert) {
    [
      { path: '/auth/{username}', result: '/auth/example' },
      { path: '{username}/foo', result: 'example/foo' },
      { path: 'foo/{username}/bar', result: 'foo/example/bar' },
      { path: '', result: '' },
      { path: undefined, result: '' },
    ].forEach((test) => {
      assert.strictEqual(pathToHelpUrlSegment(test.path), test.result, `translates ${test.path}`);
    });
  });

  test(`_getPathParam`, function (assert) {
    [
      { path: '/auth/{username}', result: 'username' },
      { path: '{unicorn}/foo', result: 'unicorn' },
      { path: 'foo/{bigfoot}/bar', result: 'bigfoot' },
      { path: '{alphabet}/bowl/{soup}', result: 'alphabet' },
      { path: 'no/params', result: false },
      { path: '', result: false },
      { path: undefined, result: false },
    ].forEach((test) => {
      assert.strictEqual(_getPathParam(test.path), test.result, `returns first param for ${test.path}`);
    });
  });

  test(`getHelpUrlForModel`, function (assert) {
    [
      { modelType: 'kmip/config', result: '/v1/foobar/config?help=1' },
      { modelType: 'does-not-exist', result: null },
      { modelType: 4, result: null },
      { modelType: '', result: null },
      { modelType: undefined, result: null },
    ].forEach((test) => {
      assert.strictEqual(
        getHelpUrlForModel(test.modelType, 'foobar'),
        test.result,
        `returns first param for ${test.path}`
      );
    });
  });
});
