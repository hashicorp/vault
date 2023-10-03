/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { sanitizePath, ensureTrailingSlash } from 'vault/services/console';
import sinon from 'sinon';

module('Unit | Service | console', function (hooks) {
  setupTest(hooks);
  hooks.beforeEach(function () {});
  hooks.afterEach(function () {});

  test('#sanitizePath', function (assert) {
    assert.strictEqual(
      sanitizePath(' /foo/bar/baz/ '),
      'foo/bar/baz',
      'removes spaces and slashs on either side'
    );
    assert.strictEqual(sanitizePath('//foo/bar/baz/'), 'foo/bar/baz', 'removes more than one slash');
  });

  test('#ensureTrailingSlash', function (assert) {
    assert.strictEqual(ensureTrailingSlash('foo/bar'), 'foo/bar/', 'adds trailing slash');
    assert.strictEqual(ensureTrailingSlash('baz/'), 'baz/', 'keeps trailing slash if there is one');
  });

  const testCases = [
    {
      method: 'read',
      args: ['/sys/health', {}],
      expectedURL: 'sys/health',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: undefined },
    },

    {
      method: 'read',
      args: ['/secrets/foo/bar', {}, { wrapTTL: '30m' }],
      expectedURL: 'secrets/foo/bar',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: '30m' },
    },

    {
      method: 'write',
      args: ['aws/roles/my-other-role', { arn: 'arn=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess' }],
      expectedURL: 'aws/roles/my-other-role',
      expectedVerb: 'POST',
      expectedOptions: {
        data: { arn: 'arn=arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess' },
        wrapTTL: undefined,
      },
    },

    {
      method: 'list',
      args: ['secret/mounts', {}],
      expectedURL: 'secret/mounts/',
      expectedVerb: 'GET',
      expectedOptions: { data: { list: true }, wrapTTL: undefined },
    },

    {
      method: 'list',
      args: ['secret/mounts', {}, { wrapTTL: '1h' }],
      expectedURL: 'secret/mounts/',
      expectedVerb: 'GET',
      expectedOptions: { data: { list: true }, wrapTTL: '1h' },
    },

    {
      method: 'delete',
      args: ['secret/secrets/kv'],
      expectedURL: 'secret/secrets/kv',
      expectedVerb: 'DELETE',
      expectedOptions: { data: undefined, wrapTTL: undefined },
    },
  ];

  test('it reads, writes, lists, deletes', function (assert) {
    assert.expect(18);
    const ajax = sinon.stub();
    const uiConsole = this.owner.factoryFor('service:console').create({
      adapter() {
        return {
          buildURL(url) {
            return url;
          },
          ajax,
        };
      },
    });

    testCases.forEach((testCase) => {
      uiConsole[testCase.method](...testCase.args);
      const [url, verb, options] = ajax.lastCall.args;
      assert.strictEqual(url, testCase.expectedURL, `${testCase.method}: uses trimmed passed url`);
      assert.strictEqual(verb, testCase.expectedVerb, `${testCase.method}: uses the correct verb`);
      assert.deepEqual(options, testCase.expectedOptions, `${testCase.method}: uses the correct options`);
    });
  });

  const kvTestCases = [
    {
      method: 'kvGet',
      args: ['kv/foo'],
      expectedURL: 'kv/data/foo',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: undefined },
    },
    {
      method: 'kvGet',
      args: ['kv/foo', {}, { metadata: true }],
      expectedURL: 'kv/metadata/foo',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: undefined },
    },
    {
      method: 'kvGet',
      args: ['kv/foo', {}, { wrapTTL: '10m' }],
      expectedURL: 'kv/data/foo',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: '10m' },
    },
    {
      method: 'kvGet',
      args: ['kv/foo', {}, { metadata: true, wrapTTL: '10m' }],
      expectedURL: 'kv/metadata/foo',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: '10m' },
    },
  ];

  test('it reads kv secret and metadata', function (assert) {
    assert.expect(12);
    const ajax = sinon.stub();
    const uiConsole = this.owner.factoryFor('service:console').create({
      adapter() {
        return {
          buildURL(url) {
            return url;
          },
          ajax,
        };
      },
    });

    kvTestCases.forEach((testCase) => {
      uiConsole[testCase.method](...testCase.args);
      const [url, verb, options] = ajax.lastCall.args;
      assert.strictEqual(url, testCase.expectedURL, `${testCase.method}: uses correct url`);
      assert.strictEqual(verb, testCase.expectedVerb, `${testCase.method}: uses the correct verb`);
      assert.deepEqual(options, testCase.expectedOptions, `${testCase.method}: uses the correct options`);
    });
  });
});
