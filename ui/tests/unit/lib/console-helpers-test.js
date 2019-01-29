import { module, test } from 'qunit';
import {
  parseCommand,
  extractDataAndFlags,
  logFromResponse,
  logFromError,
  logErrorFromInput,
} from 'vault/lib/console-helpers';

module('Unit | Lib | console helpers', function() {
  const testCommands = [
    {
      name: 'write with data',
      command: `vault write aws/config/root \
      access_key=AKIAJWVN5Z4FOFT7NLNA \
      secret_key=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 \
      region=us-east-1`,
      expected: [
        'write',
        [],
        'aws/config/root',
        [
          'access_key=AKIAJWVN5Z4FOFT7NLNA',
          'secret_key=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',
          'region=us-east-1',
        ],
      ],
    },
    {
      name: 'read with field',
      command: `vault read -field=access_key aws/creds/my-role`,
      expected: ['read', ['-field=access_key'], 'aws/creds/my-role', []],
    },
  ];

  testCommands.forEach(function(testCase) {
    test(`#parseCommand: ${testCase.name}`, function(assert) {
      let result = parseCommand(testCase.command);
      assert.deepEqual(result, testCase.expected);
    });
  });

  test('#parseCommand: invalid commands', function(assert) {
    let command = 'vault kv get foo';
    let result = parseCommand(command);
    assert.equal(result, false, 'parseCommand returns false by default');

    assert.throws(
      () => {
        parseCommand(command, true);
      },
      /invalid command/,
      'throws on invalid command when `shouldThrow` is true'
    );
  });

  const testExtractCases = [
    {
      name: 'data fields',
      input: [
        [
          'access_key=AKIAJWVN5Z4FOFT7NLNA',
          'secret_key=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',
          'region=us-east-1',
        ],
        [],
      ],
      expected: {
        data: {
          access_key: 'AKIAJWVN5Z4FOFT7NLNA',
          secret_key: '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',
          region: 'us-east-1',
        },
        flags: {},
      },
    },
    {
      name: 'repeated data and a flag',
      input: [['allowed_domains=example.com', 'allowed_domains=foo.example.com'], ['-wrap-ttl=2h']],
      expected: {
        data: {
          allowed_domains: ['example.com', 'foo.example.com'],
        },
        flags: {
          wrapTTL: '2h',
        },
      },
    },
    {
      name: 'data with more than one equals sign',
      input: [['foo=bar=baz', 'foo=baz=bop', 'some=value=val'], []],
      expected: {
        data: {
          foo: ['bar=baz', 'baz=bop'],
          some: 'value=val',
        },
        flags: {},
      },
    },
  ];

  testExtractCases.forEach(function(testCase) {
    test(`#extractDataAndFlags: ${testCase.name}`, function(assert) {
      let { data, flags } = extractDataAndFlags(...testCase.input);
      assert.deepEqual(data, testCase.expected.data, 'has expected data');
      assert.deepEqual(flags, testCase.expected.flags, 'has expected flags');
    });
  });

  let testResponseCases = [
    {
      name: 'write response, no content',
      args: [null, 'foo/bar', 'write', {}],
      expectedData: {
        type: 'success',
        content: 'Success! Data written to: foo/bar',
      },
    },
    {
      name: 'delete response, no content',
      args: [null, 'foo/bar', 'delete', {}],
      expectedData: {
        type: 'success',
        content: 'Success! Data deleted (if it existed) at: foo/bar',
      },
    },
    {
      name: 'read, no data, auth, wrap_info',
      args: [{ foo: 'bar', one: 'two' }, 'foo/bar', 'read', {}],
      expectedData: {
        type: 'object',
        content: { foo: 'bar', one: 'two' },
      },
    },
    {
      name: 'read with -format=json flag, no data, auth, wrap_info',
      args: [{ foo: 'bar', one: 'two' }, 'foo/bar', 'read', { format: 'json' }],
      expectedData: {
        type: 'json',
        content: { foo: 'bar', one: 'two' },
      },
    },
    {
      name: 'read with -field flag, no data, auth, wrap_info',
      args: [{ foo: 'bar', one: 'two' }, 'foo/bar', 'read', { field: 'one' }],
      expectedData: {
        type: 'text',
        content: 'two',
      },
    },
    {
      name: 'write, with content',
      args: [{ data: { one: 'two' } }, 'foo/bar', 'write', {}],
      expectedData: {
        type: 'object',
        content: { one: 'two' },
      },
    },
    {
      name: 'with wrap-ttl flag',
      args: [{ wrap_info: { one: 'two' } }, 'foo/bar', 'read', { wrapTTL: '1h' }],
      expectedData: {
        type: 'object',
        content: { one: 'two' },
      },
    },
    {
      name: 'with -format=json flag and wrap-ttl flag',
      args: [{ foo: 'bar', wrap_info: { one: 'two' } }, 'foo/bar', 'read', { format: 'json', wrapTTL: '1h' }],
      expectedData: {
        type: 'json',
        content: { foo: 'bar', wrap_info: { one: 'two' } },
      },
    },
    {
      name: 'with -format=json and -field flags',
      args: [{ foo: 'bar', data: { one: 'two' } }, 'foo/bar', 'read', { format: 'json', field: 'one' }],
      expectedData: {
        type: 'json',
        content: 'two',
      },
    },
    {
      name: 'with -format=json and -field, and -wrap-ttl flags',
      args: [
        { foo: 'bar', wrap_info: { one: 'two' } },
        'foo/bar',
        'read',
        { format: 'json', wrapTTL: '1h', field: 'one' },
      ],
      expectedData: {
        type: 'json',
        content: 'two',
      },
    },
    {
      name: 'with string field flag and wrap-ttl flag',
      args: [{ foo: 'bar', wrap_info: { one: 'two' } }, 'foo/bar', 'read', { field: 'one', wrapTTL: '1h' }],
      expectedData: {
        type: 'text',
        content: 'two',
      },
    },
    {
      name: 'with object field flag and wrap-ttl flag',
      args: [
        { foo: 'bar', wrap_info: { one: { two: 'three' } } },
        'foo/bar',
        'read',
        { field: 'one', wrapTTL: '1h' },
      ],
      expectedData: {
        type: 'object',
        content: { two: 'three' },
      },
    },
    {
      name: 'with response data and string field flag',
      args: [{ foo: 'bar', data: { one: 'two' } }, 'foo/bar', 'read', { field: 'one', wrapTTL: '1h' }],
      expectedData: {
        type: 'text',
        content: 'two',
      },
    },
    {
      name: 'with response data and object field flag ',
      args: [
        { foo: 'bar', data: { one: { two: 'three' } } },
        'foo/bar',
        'read',
        { field: 'one', wrapTTL: '1h' },
      ],
      expectedData: {
        type: 'object',
        content: { two: 'three' },
      },
    },
    {
      name: 'response with data',
      args: [{ foo: 'bar', data: { one: 'two' } }, 'foo/bar', 'read', {}],
      expectedData: {
        type: 'object',
        content: { one: 'two' },
      },
    },
    {
      name: 'with response data, field flag, and field missing',
      args: [{ foo: 'bar', data: { one: 'two' } }, 'foo/bar', 'read', { field: 'foo' }],
      expectedData: {
        type: 'error',
        content: 'Field "foo" not present in secret',
      },
    },
    {
      name: 'with response data and auth block',
      args: [{ data: { one: 'two' }, auth: { three: 'four' } }, 'auth/token/create', 'write', {}],
      expectedData: {
        type: 'object',
        content: { three: 'four' },
      },
    },
    {
      name: 'with -field and -format with an object field',
      args: [{ data: { one: { three: 'two' } } }, 'sys/mounts', 'read', { field: 'one', format: 'json' }],
      expectedData: {
        type: 'json',
        content: { three: 'two' },
      },
    },
    {
      name: 'with -field and -format with a string field',
      args: [{ data: { one: 'two' } }, 'sys/mounts', 'read', { field: 'one', format: 'json' }],
      expectedData: {
        type: 'json',
        content: 'two',
      },
    },
  ];

  testResponseCases.forEach(function(testCase) {
    test(`#logFromResponse: ${testCase.name}`, function(assert) {
      let data = logFromResponse(...testCase.args);
      assert.deepEqual(data, testCase.expectedData);
    });
  });

  let testErrorCases = [
    {
      name: 'AdapterError write',
      args: [{ httpStatus: 404, path: 'v1/sys/foo', errors: [{}] }, 'sys/foo', 'write'],
      expectedContent: 'Error writing to: sys/foo.\nURL: v1/sys/foo\nCode: 404',
    },
    {
      name: 'AdapterError read',
      args: [{ httpStatus: 404, path: 'v1/sys/foo', errors: [{}] }, 'sys/foo', 'read'],
      expectedContent: 'Error reading from: sys/foo.\nURL: v1/sys/foo\nCode: 404',
    },
    {
      name: 'AdapterError list',
      args: [{ httpStatus: 404, path: 'v1/sys/foo', errors: [{}] }, 'sys/foo', 'list'],
      expectedContent: 'Error listing: sys/foo.\nURL: v1/sys/foo\nCode: 404',
    },
    {
      name: 'AdapterError delete',
      args: [{ httpStatus: 404, path: 'v1/sys/foo', errors: [{}] }, 'sys/foo', 'delete'],
      expectedContent: 'Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404',
    },
    {
      name: 'VaultError single error',
      args: [{ httpStatus: 404, path: 'v1/sys/foo', errors: ['no client token'] }, 'sys/foo', 'delete'],
      expectedContent: 'Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404\nErrors:\n  no client token',
    },
    {
      name: 'VaultErrors multiple errors',
      args: [
        { httpStatus: 404, path: 'v1/sys/foo', errors: ['no client token', 'this is an error'] },
        'sys/foo',
        'delete',
      ],
      expectedContent:
        'Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404\nErrors:\n  no client token\n  this is an error',
    },
  ];

  testErrorCases.forEach(function(testCase) {
    test(`#logFromError: ${testCase.name}`, function(assert) {
      let data = logFromError(...testCase.args);
      assert.deepEqual(
        data,
        { type: 'error', content: testCase.expectedContent },
        'returns the expected data'
      );
    });
  });

  const testCommandCases = [
    {
      name: 'errors when command does not include a path',
      args: [],
      expectedContent: 'A path is required to make a request.',
    },
    {
      name: 'errors when write command does not include data and does not have force tag',
      args: ['foo/bar', 'write', {}, []],
      expectedContent: 'Must supply data or use -force',
    },
  ];

  testCommandCases.forEach(function(testCase) {
    test(`#logErrorFromInput: ${testCase.name}`, function(assert) {
      let data = logErrorFromInput(...testCase.args);

      assert.deepEqual(
        data,
        { type: 'error', content: testCase.expectedContent },
        'returns the pcorrect data'
      );
    });
  });
});
