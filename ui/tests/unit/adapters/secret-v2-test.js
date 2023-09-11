/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

module('Unit | Adapter | secret-v2', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.server = apiStub();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  [
    ['query', null, {}, { id: '', backend: 'secret' }, 'GET', '/v1/secret/metadata/?list=true'],
    ['queryRecord', null, {}, { id: 'foo', backend: 'secret' }, 'GET', '/v1/secret/metadata/foo'],
    [
      'updateRecord',
      {
        serializerFor() {
          return {
            serializeIntoHash() {},
          };
        },
      },
      {},
      {
        id: 'foo',
        belongsTo() {
          return 'secret';
        },
      },
      'PUT',
      '/v1/secret/metadata/foo',
    ],
    [
      'deleteRecord',
      {
        serializerFor() {
          return {
            serializeIntoHash() {},
          };
        },
      },
      {},
      {
        id: 'foo',
        belongsTo() {
          return 'secret';
        },
      },
      'DELETE',
      '/v1/secret/metadata/foo',
    ],
  ].forEach(([adapterMethod, store, type, queryOrSnapshot, expectedHttpVerb, expectedURL]) => {
    test(`secret-v2: ${adapterMethod}`, function (assert) {
      const adapter = this.owner.lookup('adapter:secret-v2');
      adapter[adapterMethod](store, type, queryOrSnapshot);
      const { url, method } = this.server.handledRequests[0];
      assert.strictEqual(url, expectedURL, `${adapterMethod} calls the correct url: ${expectedURL}`);
      assert.strictEqual(
        method,
        expectedHttpVerb,
        `${adapterMethod} uses the correct http verb: ${expectedHttpVerb}`
      );
    });
  });
});
