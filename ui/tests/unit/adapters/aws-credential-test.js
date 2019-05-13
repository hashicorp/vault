import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

module('Unit | Adapter | aws credential', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  const storeStub = {
    pushPayload() {},
    serializerFor() {
      return {
        serializeIntoHash() {},
      };
    },
  };

  let makeSnapshot = obj => {
    obj.role = {
      backend: 'aws',
      name: 'foo',
    };
    obj.attr = attr => obj[attr];
    return obj;
  };

  const type = {
    modelName: 'aws-credential',
  };

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
      'assumed_role type no arn, no ttl',
      [storeStub, type, makeSnapshot({ credentialType: 'assumed_role' })],
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
  cases.forEach(([description, args, expectedMethod, expectedRequestBody]) => {
    test(`aws-credential: ${description}`, function(assert) {
      assert.expect(3);
      let adapter = this.owner.lookup('adapter:aws-credential');
      adapter.createRecord(...args);
      let { method, url, requestBody } = this.server.handledRequests[0];
      assert.equal(url, '/v1/aws/creds/foo', `calls the correct url`);
      assert.equal(method, expectedMethod, `${description} uses the correct http verb: ${expectedMethod}`);
      assert.equal(requestBody, JSON.stringify(expectedRequestBody));
    });
  });
});
