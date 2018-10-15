import { module, test } from 'qunit';
import { Machine } from 'xstate';
import ToolsMachineConfig from 'vault/machines/tools-machine';

module('Unit | Machine | tools-machine', function() {
  const toolsMachine = Machine(ToolsMachineConfig);

  const testCases = [
    {
      currentState: toolsMachine.initialState,
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'wrapped',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-wrapped' }],
      },
    },
    {
      currentState: 'wrapped',
      event: 'LOOKUP',
      params: null,
      expectedResults: {
        value: 'lookup',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-lookup' }],
      },
    },
    {
      currentState: 'lookup',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'info',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-info' }],
      },
    },
    {
      currentState: 'info',
      event: 'REWRAP',
      params: null,
      expectedResults: {
        value: 'rewrap',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-rewrap' }],
      },
    },
    {
      currentState: 'rewrap',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'rewrapped',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-rewrapped' }],
      },
    },
    {
      currentState: 'rewrapped',
      event: 'UNWRAP',
      params: null,
      expectedResults: {
        value: 'unwrap',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-unwrap' }],
      },
    },
    {
      currentState: 'unwrap',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'unwrapped',
        actions: [{ type: 'render', level: 'feature', component: 'wizard/tools-unwrapped' }],
      },
    },
    {
      currentState: 'unwrapped',
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
      let result = toolsMachine.transition(testCase.currentState, testCase.event, testCase.params);
      assert.equal(result.value, testCase.expectedResults.value);
      assert.deepEqual(result.actions, testCase.expectedResults.actions);
    });
  });
});
