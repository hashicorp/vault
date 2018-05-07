import { moduleForComponent, test } from 'ember-qunit';
import sinon from 'sinon';
import Ember from 'ember';

moduleForComponent('console/ui-panel', 'Unit | Component | console/ui-panel', {
  unit: true,
  needs: ['service:auth', 'service:console', 'service:flash-messages'],
});

const testCommands = [
  {
    command: `vault write aws/config/root \
    access_key=AKIAJWVN5Z4FOFT7NLNA \
    secret_key=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 \
    region=us-east-1`,
    expected: [
      'write',
      'aws/config/root',
      [
        'access_key=AKIAJWVN5Z4FOFT7NLNA',
        'secret_key=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',
        'region=us-east-1'
      ]
    ]
  },
  {
    command:`vault read aws/creds/my-role -field=access_key`,
    expected: [
      'read',
      'aws/creds/my-role',
      ['-field=access_key']
    ]
  },

];

test('#parseCommand', function(assert) {
  let panel = this.subject();
  testCommands.forEach(test => {
    let result = panel.parseCommand(test.command);
    assert.deepEqual(result, test.expected);
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
      'access_key=AKIAJWVN5Z4FOFT7NLNA',
      'secret_key=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',
      'region=us-east-1'
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
     'allowed_domains=example.com',
     'allowed_domains=foo.example.com',
     '-wrap-ttl=2h',
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
     'foo=bar=baz',
     'foo=baz=bop',
     'some=value=val'
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

test('#extractDataAndFlags', function(assert) {
 let panel = this.subject();
  testExtractCases.forEach(test => {
    let {data, flags} = panel.extractDataAndFlags(test.input);
    assert.deepEqual(data, test.expected.data, `${test.name}: has expected data`);
    assert.deepEqual(flags, test.expected.flags, `${test.name}: has expected flags`);
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
        type: 'text',
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
        type: 'text',
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
];

test('#processResponse', function(assert) {
  let panel = this.subject({
    appendToLog: sinon.spy()
  });
  testResponseCases.forEach(test => {
    panel.processResponse(...test.args);

    let spy = panel.appendToLog;
    let commandArgs = spy.getCall(spy.callCount - 2).args;
    let appendArgs = spy.lastCall.args;
    assert.deepEqual(commandArgs[0], {type: 'command', content: test.expectedCommand}, `${test.name}: calls appendToLog with the expected args`);
    assert.deepEqual(appendArgs, test.expectedLogArgs, `${test.name}: calls appendToLog with the expected args`);
  });
});
