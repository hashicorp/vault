export default {
  key: 'tutorial',
  initial: 'idle',
  on: {
    DISMISS: 'dismissed',
    DONE: 'complete',
    PAUSE: 'paused',
  },
  states: {
    init: {
      key: 'init',
      initial: 'idle',
      on: { INITDONE: 'active.select' },
      onEntry: [
        'showTutorialAlways',
        { type: 'render', level: 'tutorial', component: 'wizard/tutorial-idle' },
        { type: 'render', level: 'feature', component: null },
      ],
      onExit: ['showTutorialWhenAuthenticated', 'clearFeatureData'],
      states: {
        idle: {
          on: {
            START: 'active.setup',
            SAVE: 'active.save',
            UNSEAL: 'active.unseal',
            LOGIN: 'active.login',
          },
        },
        active: {
          onEntry: { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
          states: {
            setup: {
              on: { TOSAVE: 'save' },
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-setup' },
            },
            save: {
              on: {
                TOUNSEAL: 'unseal',
                TOLOGIN: 'login',
              },
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-save-keys' },
            },
            unseal: {
              on: { TOLOGIN: 'login' },
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-unseal' },
            },
            login: {
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-login' },
            },
          },
        },
      },
    },
    active: {
      key: 'feature',
      initial: 'select',
      onEntry: { type: 'render', level: 'tutorial', component: 'wizard/tutorial-active' },
      states: {
        select: {
          on: {
            CONTINUE: 'feature',
          },
          onEntry: { type: 'render', level: 'feature', component: 'wizard/features-selection' },
        },
        feature: {},
      },
    },
    idle: {
      on: {
        INIT: 'init.idle',
        AUTH: 'active.select',
        CONTINUE: 'active',
      },
      onEntry: [
        { type: 'render', level: 'feature', component: null },
        { type: 'render', level: 'step', component: null },
        { type: 'render', level: 'detail', component: null },
        { type: 'render', level: 'tutorial', component: 'wizard/tutorial-idle' },
      ],
    },
    dismissed: {
      onEntry: [
        { type: 'render', level: 'tutorial', component: null },
        { type: 'render', level: 'feature', component: null },
        { type: 'render', level: 'step', component: null },
        { type: 'render', level: 'detail', component: null },
        'handleDismissed',
      ],
    },
    paused: {
      on: {
        CONTINUE: 'active.feature',
      },
      onEntry: [
        { type: 'render', level: 'feature', component: null },
        { type: 'render', level: 'step', component: null },
        { type: 'render', level: 'detail', component: null },
        { type: 'render', level: 'tutorial', component: 'wizard/tutorial-paused' },
        'handlePaused',
      ],
      onExit: ['handleResume'],
    },
    complete: {
      onEntry: [
        { type: 'render', level: 'feature', component: null },
        { type: 'render', level: 'step', component: null },
        { type: 'render', level: 'detail', component: null },
        { type: 'render', level: 'tutorial', component: 'wizard/tutorial-complete' },
      ],
    },
  },
};
