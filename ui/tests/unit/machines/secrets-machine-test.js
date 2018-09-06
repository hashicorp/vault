import { Machine } from 'xstate';
import { moduleFor, test } from 'ember-qunit';
import SecretsMachineConfig from 'vault/machines/secrets-machine';

moduleFor('machine:secrets-machine', 'Unit | Machine | secrets-machine', {
  beforeEach() {},
  afterEach() {},
});

const secretsMachine = Machine(SecretsMachineConfig);

const testCases = [
  {
    args: [secretsMachine.initialState, 'CONTINUE'],
    expectedResults: {
      value: 'enable',
      actions: [
        { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        { component: 'wizard/secrets-enable', level: 'step', type: 'render' },
      ],
    },
  },
  {
    args: ['enable', 'CONTINUE', 'pki'],
    expectedResults: {
      value: 'details',
      actions: [
        { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        { component: 'wizard/secrets-details', level: 'step', type: 'render' },
      ],
    },
  },
  {
    args: ['details', 'CONTINUE', 'pki'],
    expectedResults: {
      value: 'role',
      actions: [
        { component: 'wizard/secrets-role', level: 'step', type: 'render' },
        { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
      ],
    },
  },
  {
    args: ['role', 'CONTINUE', 'pki'],
    expectedResults: {
      value: 'displayRole',
      actions: [
        { component: 'wizard/secrets-display-role', level: 'step', type: 'render' },
        { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
      ],
    },
  },
  {
    args: ['displayRole', 'CONTINUE', 'pki'],
    expectedResults: {
      value: 'credentials',
      actions: [
        { component: 'wizard/secrets-credentials', level: 'step', type: 'render' },
        { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
      ],
    },
  },
];

test('it calculates state and actions properly', function(assert) {
  testCases.forEach(testCase => {
    let result = secretsMachine.transition(...testCase.args);
    assert.equal(result.value, testCase.expectedResults.value);
    assert.deepEqual(result.actions, testCase.expectedResults.actions);
  });
});
