import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { sanitizePath, ensureTrailingSlash } from 'vault/services/console';
import sinon from 'sinon';

module('Unit | Service | console', function(hooks) {
  setupTest(hooks);
  hooks.beforeEach(function() {});
  hooks.afterEach(function() {});

  test('#sanitizePath', function(assert) {
    assert.equal(sanitizePath(' /foo/bar/baz/ '), 'foo/bar/baz', 'removes spaces and slashs on either side');
    assert.equal(sanitizePath('//foo/bar/baz/'), 'foo/bar/baz', 'removes more than one slash');
  });

  test('#ensureTrailingSlash', function(assert) {
    assert.equal(ensureTrailingSlash('foo/bar'), 'foo/bar/', 'adds trailing slash');
    assert.equal(ensureTrailingSlash('baz/'), 'baz/', 'keeps trailing slash if there is one');
  });

  let testCases = [
    {
      method: 'read',
      args: ['/sys/health', {}],
      expectedURL: 'sys/health',
      expectedVerb: 'GET',
      expectedOptions: { data: undefined, wrapTTL: undefined },
    },

    {
      method: 'read',
      args: ['/secrets/foo/bar', {}, '30m'],
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
      args: ['secret/mounts', {}, '1h'],
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

  test('it reads, writes, lists, deletes', function(assert) {
    let ajax = sinon.stub();
    let uiConsole = this.owner.factoryFor('service:console').create({
      adapter() {
        return {
          buildURL(url) {
            return url;
          },
          ajax,
        };
      },
    });

    testCases.forEach(testCase => {
      uiConsole[testCase.method](...testCase.args);
      let [url, verb, options] = ajax.lastCall.args;
      assert.equal(url, testCase.expectedURL, `${testCase.method}: uses trimmed passed url`);
      assert.equal(verb, testCase.expectedVerb, `${testCase.method}: uses the correct verb`);
      assert.deepEqual(options, testCase.expectedOptions, `${testCase.method}: uses the correct options`);
    });
  });
});
