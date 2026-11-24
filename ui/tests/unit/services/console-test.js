/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Service | console', function (hooks) {
  setupTest(hooks);
  hooks.afterEach(function () {
    sinon.restore();
  });

  const testCases = [
    {
      method: 'read',
      args: ['/sys/health', {}],
      expectedURL: 'sys/health',
      expectedVerb: 'GET',
      expectedHeaders: undefined,
    },

    {
      method: 'read',
      args: ['/secrets/foo/bar', {}, { wrapTTL: '30m' }],
      expectedURL: 'secrets/foo/bar',
      expectedVerb: 'GET',
      expectedHeaders: { 'X-Vault-Wrap-TTL': '30m' },
    },

    {
      method: 'write',
      args: ['aws/roles/my-other-role', { arn: 'arn=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess' }],
      expectedURL: 'aws/roles/my-other-role',
      expectedVerb: 'POST',
      expectedData: { arn: 'arn=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess' },
      expectedHeaders: undefined,
    },

    {
      method: 'list',
      args: ['secret/mounts', {}],
      expectedURL: 'secret/mounts/?list=true',
      expectedVerb: 'GET',
      expectedHeaders: undefined,
    },

    {
      method: 'list',
      args: ['secret/mounts', {}, { wrapTTL: '1h' }],
      expectedURL: 'secret/mounts/?list=true',
      expectedVerb: 'GET',
      expectedHeaders: { 'X-Vault-Wrap-TTL': '1h' },
    },

    {
      method: 'delete',
      args: ['secret/secrets/kv'],
      expectedURL: 'secret/secrets/kv',
      expectedVerb: 'DELETE',
      expectedHeaders: undefined,
    },
  ];

  test('it reads, writes, lists, deletes', async function (assert) {
    const response = { text: async () => '' };
    const api = this.owner.lookup('service:api');
    const get = sinon.stub(api.request, 'get').resolves(response);
    const post = sinon.stub(api.request, 'post').resolves(response);
    const del = sinon.stub(api.request, 'delete').resolves(response);
    const stubs = { GET: get, POST: post, DELETE: del };
    const uiConsole = this.owner.lookup('service:console');

    for (const testCase of testCases) {
      await uiConsole[testCase.method](...testCase.args);
      const stub = stubs[testCase.expectedVerb];
      const [url, ...rest] = stub.lastCall.args;
      assert.strictEqual(url, testCase.expectedURL, `${testCase.method}: uses trimmed passed url`);
      const [data, headers] = testCase.expectedVerb === 'GET' ? [undefined, rest[0]] : rest;
      assert.deepEqual(headers, testCase.expectedHeaders, `${testCase.method}: uses the correct headers`);
      assert.deepEqual(data, testCase.expectedData, `${testCase.method}: sends the data as the body`);
    }
  });

  test('it parses json responses and tolerates empty ones', async function (assert) {
    const api = this.owner.lookup('service:api');
    const get = sinon.stub(api.request, 'get');
    get.onFirstCall().resolves({ text: async () => JSON.stringify({ data: { foo: 'bar' } }) });
    get.onSecondCall().resolves({ text: async () => '' });
    const uiConsole = this.owner.lookup('service:console');

    assert.deepEqual(await uiConsole.read('secret/foo'), { data: { foo: 'bar' } }, 'json body is parsed');
    assert.strictEqual(await uiConsole.read('secret/foo'), undefined, 'empty body resolves undefined');
  });

  const kvTestCases = [
    {
      method: 'kvGet',
      args: ['kv/foo'],
      expectedURL: 'kv/data/foo',
      expectedVerb: 'GET',
      expectedHeaders: undefined,
    },
    {
      method: 'kvGet',
      args: ['kv/foo', {}, { metadata: true }],
      expectedURL: 'kv/metadata/foo',
      expectedVerb: 'GET',
      expectedHeaders: undefined,
    },
    {
      method: 'kvGet',
      args: ['kv/foo', {}, { wrapTTL: '10m' }],
      expectedURL: 'kv/data/foo',
      expectedVerb: 'GET',
      expectedHeaders: { 'X-Vault-Wrap-TTL': '10m' },
    },
    {
      method: 'kvGet',
      args: ['kv/foo', {}, { metadata: true, wrapTTL: '10m' }],
      expectedURL: 'kv/metadata/foo',
      expectedVerb: 'GET',
      expectedHeaders: { 'X-Vault-Wrap-TTL': '10m' },
    },
  ];

  test('it reads kv secret and metadata', async function (assert) {
    const response = { text: async () => '' };
    const api = this.owner.lookup('service:api');
    const get = sinon.stub(api.request, 'get').resolves(response);
    const uiConsole = this.owner.lookup('service:console');

    for (const testCase of kvTestCases) {
      await uiConsole[testCase.method](...testCase.args);
      const [url, headers] = get.lastCall.args;
      assert.strictEqual(url, testCase.expectedURL, `${testCase.method}: uses correct url`);
      assert.deepEqual(headers, testCase.expectedHeaders, `${testCase.method}: uses the correct headers`);
    }
  });
});
