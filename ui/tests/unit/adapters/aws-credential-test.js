/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

const storeStub = {
  pushPayload() {},
  serializerFor() {
    return {
      serializeIntoHash() {},
    };
  },
};

const makeSnapshot = (obj) => {
  obj.role = {
    backend: 'aws',
    name: 'foo',
  };
  obj.attr = (attr) => obj[attr];
  return obj;
};

const type = {
  modelName: 'aws-credential',
};
module('Unit | Adapter | aws credential', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.getAssertion = () => {};
    this.postAssertion = () => {};
    this.server.get('/aws/creds/foo', (schema, req) => {
      this.getAssertion(req);
      return {};
    });
    this.server.post('/aws/creds/foo', (schema, req) => {
      this.postAssertion(req);
      return {};
    });
  });

  const cases = [
    ['iam_user type', [storeStub, type, makeSnapshot({ credentialType: 'iam_user', ttl: '3h' })], 'GET'],
    [
      'federation_token type with ttl',
      [storeStub, type, makeSnapshot({ credentialType: 'federation_token', ttl: '3h', roleArn: 'arn' })],
      'POST',
      { ttl: '3h' },
    ],
    [
      'federation_token type no ttl',
      [storeStub, type, makeSnapshot({ credentialType: 'federation_token', roleArn: 'arn' })],
      'POST',
    ],
    [
      'session_token type with ttl',
      [storeStub, type, makeSnapshot({ credentialType: 'session_token', ttl: '3h' })],
      'POST',
      { ttl: '3h' },
    ],
    [
      'session_token type no ttl',
      [storeStub, type, makeSnapshot({ credentialType: 'session_token' })],
      'POST',
    ],
    [
      'assumed_role type no arn, no ttl',
      [storeStub, type, makeSnapshot({ credentialType: 'assumed_role' })],
      'POST',
    ],
    [
      'assumed_role type no arn, ttl empty',
      [storeStub, type, makeSnapshot({ credentialType: 'assumed_role', ttl: '' })],
      'POST',
    ],
    [
      'assumed_role type no arn',
      [storeStub, type, makeSnapshot({ credentialType: 'assumed_role', ttl: '3h' })],
      'POST',
      { ttl: '3h' },
    ],
    [
      'assumed_role type',
      [storeStub, type, makeSnapshot({ credentialType: 'assumed_role', roleArn: 'arn', ttl: '3h' })],
      'POST',
      { ttl: '3h', role_arn: 'arn' },
    ],
  ];

  cases.forEach(([description, args, method, expectedRequestBody]) => {
    test(`aws-credential: ${description}`, function (assert) {
      assert.expect(2);
      const assertionName = method === 'GET' ? 'getAssertion' : 'postAssertion';
      this.set(assertionName, (req) => {
        assert.strictEqual(req.method, method, `query calls the correct url with method ${method}`);
        const body = JSON.parse(req.requestBody);
        const expected = expectedRequestBody ? expectedRequestBody : null;
        assert.deepEqual(body, expected);
      });
      const adapter = this.owner.lookup('adapter:aws-credential');
      adapter.createRecord(...args);
    });
  });
});
