export default {
  key: 'tutorial',
  initial: 'idle',
  states: {
    active: {
      key: 'feature',
      initial: 'init',
      on: {
        DISMISS: 'dismissed',
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
        init: {
          key: 'init',
          initial: 'setup',
          on: { DONE: 'select' },
          states: {
            setup: {
              on: { CONTINUE: 'save' },
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-setup' },
            },
            save: {
              on: { CONTINUE: 'unseal' },
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-save-keys' },
            },
            unseal: {
              on: { CONTINUE: 'login' },
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-unseal' },
            },
            login: {
              onEntry: { type: 'render', level: 'feature', component: 'wizard/init-login' },
            },
          },
        },
      },
    },
    idle: {
      on: {
        DISMISS: 'dismissed',
        CONTINUE: 'active',
        AUTH: 'active.select',
        UNSEAL: 'active.unseal',
        LOGIN: 'active.login',
      },
      onEntry: { type: 'render', level: 'tutorial', component: 'wizard/tutorial-idle' },
    },
    dismissed: {
      on: { CONTINUE: 'idle' },
      onEntry: [{ type: 'render', level: 'tutorial', component: null }, 'handleDismissed'],
    },
    paused: {
      on: { CONTINUE: ['handlePause'] },
      onEntry: { type: 'render', level: 'tutorial', component: 'wizard/tutorial-paused' },
    },
    complete: {
      on: {
        CONTINUE: 'idle',
        DISMISS: 'dismissed',
      },
    },
  },
};
