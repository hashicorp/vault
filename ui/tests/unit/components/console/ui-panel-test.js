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


test('#processResponse', function(assert) {
 let panel = this.subject();

})



