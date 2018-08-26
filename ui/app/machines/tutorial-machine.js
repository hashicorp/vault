export default {
  key: 'tutorial',
  initial: 'idle',
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
      onExit: ['showTutorialWhenAuthenticated'],
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
              on: { TOUNSEAL: 'unseal' },
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
      on: {
        DISMISS: 'dismissed',
        DONE: 'complete',
      },
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
        DISMISS: 'dismissed',
        CONTINUE: 'active',
      },
      onEntry: [
        { type: 'render', level: 'tutorial', component: 'wizard/tutorial-idle' },
        { type: 'render', level: 'feature', component: null },
        { type: 'render', level: 'step', component: null },
        { type: 'render', level: 'detail', component: null },
      ],
    },
    dismissed: {
      on: { CONTINUE: 'idle' },
      onEntry: [
        { type: 'render', level: 'feature', component: null },
        { type: 'render', level: 'tutorial', component: null },
        'handleDismissed',
      ],
    },
    paused: {
      on: { CONTINUE: ['handlePause'] },
      onEntry: { type: 'render', level: 'tutorial', component: 'wizard/tutorial-paused' },
    },
    complete: {
      onEntry: { type: 'render', level: 'tutorial', component: 'wizard/tutorial-complete' },
      on: {
        PAUSE: 'idle',
        DISMISS: 'dismissed',
      },
    },
  },
};
