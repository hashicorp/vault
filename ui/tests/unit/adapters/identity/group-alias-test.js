/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import testCases from './_test-cases';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | identity/group-alias', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  const cases = testCases('identity/group-alias');

  cases.forEach((testCase) => {
    test(`group-alias#${testCase.adapterMethod}`, function (assert) {
      assert.expect(3);
      const method = testCase.method.toLowerCase();
      const url = testCase.url.replace('/v1', '').split('?')[0];
      const queryParams = testCase.url.includes('?list=true') ? { list: 'true' } : {};
      this.server[method](url, (schema, req) => {
        assert.true(
          testCase.url.includes(url),
          `${testCase.adapterMethod} calls the correct url: ${testCase.url}`
        );
        assert.strictEqual(req.method, testCase.method, `usses the correct http verb: ${testCase.method}`);
        assert.deepEqual(req.queryParams, queryParams, 'calls with correct query params');
        return {};
      });
      const adapter = this.owner.lookup('adapter:identity/group-alias');
      adapter[testCase.adapterMethod](...testCase.args);
    });
  });
});
