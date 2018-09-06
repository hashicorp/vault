<<<<<<< HEAD
import { module, test } from 'qunit';
import { Machine } from 'xstate';
import ReplicationMachineConfig from 'vault/machines/replication-machine';

module('Unit | Machine | replication-machine', function() {
  const replicationMachine = Machine(ReplicationMachineConfig);

  const testCases = [
    {
      currentState: replicationMachine.initialState,
      event: 'ENABLEREPLICATION',
      params: null,
      expectedResults: {
        value: 'details',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/replication-details' }],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
  ];

  testCases.forEach(testCase => {
    test(`transition: ${testCase.event} for currentState ${testCase.currentState} and componentState ${
      testCase.params
    }`, function(assert) {
      let result = replicationMachine.transition(testCase.currentState, testCase.event, testCase.params);
      assert.equal(result.value, testCase.expectedResults.value);
      assert.deepEqual(result.actions, testCase.expectedResults.actions);
    });
=======
import { Machine } from 'xstate';
import { moduleFor, test } from 'ember-qunit';
import ReplicationMachineConfig from 'vault/machines/replication-machine';

moduleFor('machine:replication-machine', 'Unit | Machine | replication-machine', {
  beforeEach() {},
  afterEach() {},
});

const replicationMachine = Machine(ReplicationMachineConfig);

const testCases = [
  {
    currentState: replicationMachine.initialState,
    event: 'ENABLEREPLICATION',
    params: null,
    expectedResults: {
      value: 'details',
      actions: [{ type: 'render', level: 'feature', component: 'wizard/replication-details' }],
    },
  },
  {
    currentState: 'details',
    event: 'CONTINUE',
    params: null,
    expectedResults: {
      value: 'complete',
      actions: ['completeFeature'],
    },
  },
];

testCases.forEach(testCase => {
  test(`transition: ${testCase.event} for currentState ${testCase.currentState} and componentState ${testCase.params}`, function(
    assert
  ) {
    let result = replicationMachine.transition(testCase.currentState, testCase.event, testCase.params);
    assert.equal(result.value, testCase.expectedResults.value);
    assert.deepEqual(result.actions, testCase.expectedResults.actions);
>>>>>>> add tests for replication machine
  });
});
