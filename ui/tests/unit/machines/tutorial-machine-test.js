import { module, test } from 'qunit';
import { Machine } from 'xstate';
import TutorialMachineConfig from 'vault/machines/tutorial-machine';

module('Unit | Machine | tutorial-machine', function() {
  const tutorialMachine = Machine(TutorialMachineConfig);

  const testCases = [
    {
      currentState: 'init',
      event: 'START',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'setup',
          },
        },
        actions: [
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/init-setup' },
        ],
      },
    },
    {
      currentState: 'init.active.setup',
      event: 'TOSAVE',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'save',
          },
        },
        actions: [{ type: 'render', level: 'feature', component: 'wizard/init-save-keys' }],
      },
    },
    {
      currentState: 'init',
      event: 'SAVE',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'save',
          },
        },
        actions: [
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/init-save-keys' },
        ],
      },
    },
    {
      currentState: 'init.active.save',
      event: 'TOUNSEAL',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'unseal',
          },
        },
        actions: [{ type: 'render', level: 'feature', component: 'wizard/init-unseal' }],
      },
    },
    {
      currentState: 'init',
      event: 'UNSEAL',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'unseal',
          },
        },
        actions: [
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/init-unseal' },
        ],
      },
    },
    {
      currentState: 'init.active.unseal',
      event: 'TOLOGIN',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'login',
          },
        },
        actions: [{ type: 'render', level: 'feature', component: 'wizard/init-login' }],
      },
    },
    {
      currentState: 'init',
      event: 'LOGIN',
      params: null,
      expectedResults: {
        value: {
          init: {
            active: 'login',
          },
        },
        actions: [
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/init-login' },
        ],
      },
    },
    {
      currentState: 'init.active.login',
      event: 'INITDONE',
      params: null,
      expectedResults: {
        value: {
          active: 'select',
        },
        actions: [
          'showTutorialWhenAuthenticated',
          'clearFeatureData',
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/features-selection' },
        ],
      },
    },
    {
      currentState: 'active.select',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: {
          active: 'feature',
        },
        actions: [],
      },
    },
    {
      currentState: 'active.feature',
      event: 'DISMISS',
      params: null,
      expectedResults: {
        value: 'dismissed',
        actions: [
          { type: 'render', level: 'tutorial', component: null },
          { type: 'render', level: 'feature', component: null },
          { type: 'render', level: 'step', component: null },
          { type: 'render', level: 'detail', component: null },
          'handleDismissed',
        ],
      },
    },
    {
      currentState: 'active.feature',
      event: 'DONE',
      params: null,
      expectedResults: {
        value: 'complete',
        actions: [
          { type: 'render', level: 'feature', component: null },
          { type: 'render', level: 'step', component: null },
          { type: 'render', level: 'detail', component: null },
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-complete' },
        ],
      },
    },
    {
      currentState: 'active.feature',
      event: 'PAUSE',
      params: null,
      expectedResults: {
        value: 'paused',
        actions: [
          { type: 'render', level: 'feature', component: null },
          { type: 'render', level: 'step', component: null },
          { type: 'render', level: 'detail', component: null },
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-paused' },
          'handlePaused',
        ],
      },
    },
    {
      currentState: 'paused',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: {
          active: 'feature',
        },
        actions: ['handleResume', { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' }],
      },
    },
    {
      currentState: 'idle',
      event: 'INIT',
      params: null,
      expectedResults: {
        value: {
          init: 'idle',
        },
        actions: [
          'showTutorialAlways',
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-idle' },
          { type: 'render', level: 'feature', component: null },
        ],
      },
    },
    {
      currentState: 'idle',
      event: 'AUTH',
      params: null,
      expectedResults: {
        value: {
          active: 'select',
        },
        actions: [
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/features-selection' },
        ],
      },
    },
    {
      currentState: 'idle',
      event: 'CONTINUE',
      params: null,
      expectedResults: {
        value: {
          active: 'select',
        },
        actions: [
          { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          { type: 'render', level: 'feature', component: 'wizard/features-selection' },
        ],
      },
    },
  ];

  testCases.forEach(testCase => {
    test(`transition: ${testCase.event} for currentState ${testCase.currentState} and componentState ${
      testCase.params
    }`, function(assert) {
      let result = tutorialMachine.transition(testCase.currentState, testCase.event, testCase.params);
      assert.deepEqual(result.value, testCase.expectedResults.value);
      assert.deepEqual(result.actions, testCase.expectedResults.actions);
    });
  });
});
