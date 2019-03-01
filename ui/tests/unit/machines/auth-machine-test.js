import { module, test } from 'qunit';
import { Machine } from 'xstate';
import AuthMachineConfig from 'vault/machines/auth-machine';

module('Unit | Machine | auth-machine', function() {
  const authMachine = Machine(AuthMachineConfig);

  const testCases = [
    {
      currentState: authMachine.initialState,
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'enable',
        actions: [
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
          { component: 'wizard/auth-enable', level: 'step', type: 'render' },
        ],
      },
    },
    {
      currentState: 'enable',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: 'config',
        actions: [
          { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
          { type: 'render', level: 'step', component: 'wizard/auth-config' },
        ],
      },
    },
    {
      currentState: 'config',
      event: 'CONTINUE',
      expectedResults: {
        value: 'details',
        actions: [
          { component: 'wizard/auth-details', level: 'step', type: 'render' },
          { component: 'wizard/mounts-wizard', level: 'feature', type: 'render' },
        ],
      },
    },
    {
      currentState: 'details',
      event: 'CONTINUE',
      expectedResults: {
        value: 'complete',
        actions: ['completeFeature'],
      },
    },
    {
      currentState: 'details',
      event: 'RESET',
      params: null,
      expectedResults: {
        value: 'idle',
        actions: [
          {
            params: ['vault.cluster.settings.auth.enable'],
            type: 'routeTransition',
          },
          {
            component: 'wizard/mounts-wizard',
            level: 'feature',
            type: 'render',
          },
          {
            component: 'wizard/auth-idle',
            level: 'step',
            type: 'render',
          },
        ],
      },
    },
  ];

  testCases.forEach(testCase => {
    test(`transition: ${testCase.event} for currentState ${testCase.currentState} and componentState ${
      testCase.params
    }`, function(assert) {
      let result = authMachine.transition(testCase.currentState, testCase.event, testCase.params);
      assert.equal(result.value, testCase.expectedResults.value);
      assert.deepEqual(result.actions, testCase.expectedResults.actions);
    });
  });
});
