import { module, test } from 'qunit';
import { Machine } from 'xstate';
import PoliciesMachineConfig from 'vault/machines/policies-machine';

module('Unit | Machine | policies-machine', function() {
  const policiesMachine = Machine(PoliciesMachineConfig);

  const testCases = [
    {
      currentState: policiesMachine.initialState,
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'create',
        actions: [{ component: 'wizard/policies-create', level: 'feature', type: 'render' }],
      },
    },
    {
      currentState: 'create',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'details',
        actions: [{ component: 'wizard/policies-details', level: 'feature', type: 'render' }],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      expectedResults: {
        value: 'delete',
        actions: [{ component: 'wizard/policies-delete', level: 'feature', type: 'render' }],
      },
    },
    {
      currentState: 'delete',
      event: 'CONTINUE',
      expectedResults: {
        value: 'others',
        actions: [{ component: 'wizard/policies-others', level: 'feature', type: 'render' }],
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
    test(`transition: ${testCase.event} for currentState ${testCase.currentState} and componentState ${
      testCase.params
    }`, function(assert) {
      let result = policiesMachine.transition(testCase.currentState, testCase.event, testCase.params);
      assert.equal(result.value, testCase.expectedResults.value);
      assert.deepEqual(result.actions, testCase.expectedResults.actions);
    });
  });
});
