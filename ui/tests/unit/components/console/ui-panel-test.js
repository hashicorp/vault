import { moduleForComponent, test } from 'ember-qunit';
import sinon from 'sinon';

moduleForComponent('console/ui-panel', 'Unit | Component | console/ui panel', {
  unit: true,
  needs: ['service:auth', 'service:console', 'service:flash-messages'],
});

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
        'region=us-east-1'
      ]
    ]
  },
  {
    name: 'read with field',
    command:`vault read -field=access_key aws/creds/my-role`,
    expected: [
      'read',
      ['-field=access_key'],
      'aws/creds/my-role',
      []
    ]
  },
];

testCommands.forEach(function(testCase) {
  test(`#parseCommand: ${testCase.name}`, function(assert) {
    let panel = this.subject();
    let result = panel.parseCommand(testCase.command);
    assert.deepEqual(result, testCase.expected);
  });
});

test('#parseCommand: invalid commands', function(assert) {
  let panel = this.subject();
  let command = 'vault kv get foo';
  let result = panel.parseCommand(command);
  assert.equal(result, false, 'parseCommand returns false by default');

  assert.throws(() => {
    panel.parseCommand(command, true)
  }, /invalid command/, 'throws on invalid command when `shouldThrow` is true');
});

const testExtractCases = [
  {
    name: 'data fields',
    input: [
      [
        'access_key=AKIAJWVN5Z4FOFT7NLNA',
        'secret_key=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',
        'region=us-east-1'
      ],
      []
    ],
    expected: {
      data: {
        'access_key': 'AKIAJWVN5Z4FOFT7NLNA',
        'secret_key': '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',
        'region': 'us-east-1'
      },
      flags: {}
    }
  },
 {
   name: 'repeated data and a flag',
   input: [
    [
     'allowed_domains=example.com',
     'allowed_domains=foo.example.com',
    ],
    ['-wrap-ttl=2h']
   ],
   expected: {
     data: {
       'allowed_domains': ['example.com', 'foo.example.com']
     },
     flags: {
       'wrapTTL': '2h'
     }
   }
 },
 {
   name: 'data with more than one equals sign',
   input: [
      [
       'foo=bar=baz',
       'foo=baz=bop',
       'some=value=val'
      ],
      []
   ],
   expected: {
     data: {
       foo: ['bar=baz', 'baz=bop'],
       some: 'value=val'
     },
     flags: {}
   }
 },
];

testExtractCases.forEach(function(testCase) {
  test(`#extractDataAndFlags: ${testCase.name}`, function(assert) {
    let panel = this.subject();
    let {data, flags} = panel.extractDataAndFlags(...testCase.input);
    assert.deepEqual(data, testCase.expected.data, 'has expected data');
    assert.deepEqual(flags, testCase.expected.flags, 'has expected flags');
  });
});

let testResponseCases = [
  {
    name: 'write response, no content',
    args: [null, 'vault write', 'foo/bar', 'write', {}],
    expectedCommand: 'vault write',
    expectedLogArgs: [
      {
        type: 'text',
        content: 'Success! Data written to: foo/bar'
      }
    ]
  },
  {
    name: 'delete response, no content',
    args: [null, 'vault delete', 'foo/bar', 'delete', {}],
    expectedCommand: 'vault delete',
    expectedLogArgs: [
      {
        type: 'text',
        content: 'Success! Data deleted (if it existed) at: foo/bar'
      }
    ]
  },
  {
    name: 'write, with content',
    args: [{data: {one: 'two'}}, 'vault write', 'foo/bar', 'write', {}],
    expectedCommand: 'vault write',
    expectedLogArgs: [
      {
        type: 'object',
        content: { one: 'two'}
      }
    ]
  },
  {
    name: 'with wrap-ttl flag',
    args: [{wrap_info: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {wrapTTL: '1h'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'object',
        content: { one: 'two'}
      }
    ]
  },
  {
    name: 'with -format=json flag and wrap-ttl flag',
    args: [{foo: 'bar', wrap_info: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {format: 'json', wrapTTL: '1h'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'json',
        content: {foo: 'bar', wrap_info: {one: 'two'}}
      }
    ]
  },
  {
    name: 'with -format=json and -field flags',
    args: [{foo: 'bar', data: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {format: 'json', field: 'one'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'json',
        content:'two'
      }
    ]
  },
  {
    name: 'with -format=json and -field, and -wrap-ttl flags',
    args: [{foo: 'bar', wrap_info: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {format: 'json', wrapTTL: '1h', field: 'one'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'json',
        content: 'two'
      }
    ]
  },
  {
    name: 'with string field flag and wrap-ttl flag',
    args: [{foo: 'bar', wrap_info: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {field: 'one', wrapTTL: '1h'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'text',
        content: 'two'
      }
    ]
  },
  {
    name: 'with object field flag and wrap-ttl flag',
    args: [{foo: 'bar', wrap_info: {one: {two: 'three'}}}, 'vault read', 'foo/bar', 'read', {field: 'one', wrapTTL: '1h'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'object',
        content: {two: 'three'}
      }
    ]
  },
  {
    name: 'with response data and string field flag',
    args: [{foo: 'bar', data: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {field: 'one', wrapTTL: '1h'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'text',
        content: 'two'
      }
    ]
  },
  {
    name: 'with response data and object field flag ',
    args: [{foo: 'bar', data: {one: {two: 'three'}}}, 'vault read', 'foo/bar', 'read', {field: 'one', wrapTTL: '1h'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'object',
        content: {two: 'three'}
      }
    ]
  },
  {
    name: 'response with data',
    args: [{foo: 'bar', data: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'object',
        content: {one: 'two'}
      }
    ]
  },
  {
    name: 'with response data, field flag, and field missing',
    args: [{foo: 'bar', data: {one: 'two'}}, 'vault read', 'foo/bar', 'read', {field: 'foo'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'error',
        content: 'Field "foo" not present in secret'
      }
    ]
  },
  {
    name: 'with response data and auth block',
    args: [{data: {one: 'two'}, auth: {three: 'four'}}, 'vault write', 'auth/token/create', 'write', {}],
    expectedCommand: 'vault write',
    expectedLogArgs: [
      {
        type: 'object',
        content: {three: 'four'}
      }
    ],
  },
  {
    name: 'with -field and -format with an object field',
    args: [{data: {one: {three: 'two'}}}, 'vault read', 'sys/mounts', 'read', {field: 'one', format: 'json'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'json',
        content: {three: 'two'}
      }
    ],
  },
  {
    name: 'with -field and -format with a string field',
    args: [{data: {one: 'two'}}, 'vault read', 'sys/mounts', 'read', {field: 'one', format: 'json'}],
    expectedCommand: 'vault read',
    expectedLogArgs: [
      {
        type: 'json',
        content: 'two'
      }
    ],
  }
];

testResponseCases.forEach(function(testCase) {
  test(`#processResponse: ${testCase.name}`, function(assert) {
    let panel = this.subject({
      appendToLog: sinon.spy()
    });
    panel.processResponse(...testCase.args);

    let spy = panel.appendToLog;
    let commandArgs = spy.getCall(spy.callCount - 2).args;
    let appendArgs = spy.lastCall.args;
    assert.deepEqual(commandArgs[0], {type: 'command', content: testCase.expectedCommand}, 'appends command');
    assert.deepEqual(appendArgs, testCase.expectedLogArgs, 'appends output');
  });
});

let testErrorCases = [
  {
    name: 'AdapterError write',
    args: ['command', 'write', 'sys/foo', { httpStatus: 404, path: 'v1/sys/foo', errors: [{}]}],
    expectedContent: "Error writing to: sys/foo.\nURL: v1/sys/foo\nCode: 404"
  },
  {
    name: 'AdapterError read',
    args: ['command', 'read', 'sys/foo', { httpStatus: 404, path: 'v1/sys/foo', errors: [{}]}],
    expectedContent: "Error reading from: sys/foo.\nURL: v1/sys/foo\nCode: 404"
  },
  {
    name: 'AdapterError list',
    args: ['command', 'list', 'sys/foo', { httpStatus: 404, path: 'v1/sys/foo', errors: [{}]}],
    expectedContent: "Error listing: sys/foo.\nURL: v1/sys/foo\nCode: 404"
  },
  {
    name: 'AdapterError delete',
    args: ['command', 'delete', 'sys/foo', { httpStatus: 404, path: 'v1/sys/foo', errors: [{}]}],
    expectedContent: "Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404"
  },
  {
    name: 'VaultError single error',
    args: ['command', 'delete', 'sys/foo', { httpStatus: 404, path: 'v1/sys/foo', errors: ['no client token']}],
    expectedContent: "Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404\nErrors:\n  no client token"
  },
  {
    name: 'VaultErrors multiple errors',
    args: ['command', 'delete', 'sys/foo', { httpStatus: 404, path: 'v1/sys/foo', errors: ['no client token', 'this is an error']}],
    expectedContent: "Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404\nErrors:\n  no client token\n  this is an error"
  }
];

testErrorCases.forEach(function(testCase) {
  test(`#handleServiceError: ${testCase.name}`, function(assert) {
    let panel = this.subject({
      appendToLog: sinon.spy()
    });
    panel.handleServiceError(...testCase.args);

    let spy = panel.appendToLog;
    let appendArgs = spy.lastCall.args;
    assert.deepEqual(appendArgs[0].content, testCase.expectedContent, 'calls appendToLog with the expected error content');
  });
});

const testCommandCases = [
  {
    name: 'errors when command does not include a path',
    command: `list`,
    expectedContent: 'A path is required to make a request.'
  },
  {
    name: 'errors when write command does not include data and does not have force tag',
    command: `write this/is/a/path`,
    expectedContent: 'Must supply data or use -force'
  }
];

testCommandCases.forEach(function(testCase) {
  test(`#executeCommand: ${testCase.name}`, function(assert) {
    let panel = this.subject({
      appendToLog: sinon.spy()
    });
    panel.executeCommand(testCase.command);

    let spy = panel.appendToLog;
    let appendArgs = spy.lastCall.args;
    assert.deepEqual(appendArgs[0].content, testCase.expectedContent, 'calls appendToLog with the expected content');
  });
});
