import { Machine } from 'xstate';
import { moduleFor, test } from 'ember-qunit';
import AuthMachineConfig from 'vault/machines/policies-machine';

moduleFor('machine:policies-machine', 'Unit | Machine | policies-machine', {
  beforeEach() {},
  afterEach() {},
});

const policiesMachine = Machine(AuthMachineConfig);

const testCases = [
  {
    currentState: policiesMachine.initialState,
    event: 'CONTINUE',
    params: null,
    expectedResults: {
      value: 'create',
      actions: [{ component: 'wizard/policies-create', level: 'step', type: 'render' }],
    },
  },
  {
    currentState: 'create',
    event: 'CONTINUE',
    params: null,
    expectedResults: {
      value: 'details',
      actions: [{ component: 'wizard/policies-details', level: 'step', type: 'render' }],
    },
  },
  {
    currentState: 'details',
    event: 'CONTINUE',
    expectedResults: {
      value: 'delete',
      actions: [{ component: 'wizard/policies-delete', level: 'step', type: 'render' }],
    },
  },
  {
    currentState: 'delete',
    event: 'CONTINUE',
    expectedResults: {
      value: 'others',
      actions: [{ component: 'wizard/policies-others', level: 'step', type: 'render' }],
    },
  },
  {
    currentState: 'others',
    event: 'CONTINUE',
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
    let result = policiesMachine.transition(testCase.currentState, testCase.event, testCase.params);
    assert.equal(result.value, testCase.expectedResults.value);
    assert.deepEqual(result.actions, testCase.expectedResults.actions);
  });
});
