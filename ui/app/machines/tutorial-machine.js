export default {
  key: 'tutorial',
  initial: 'idle',
  states: {
    active: {
      on: {
        DISMISS: 'dismissed',
        AUTH: 'select',
      },
      key: 'feature',
      initial: 'init',
      states: {
        select: {
          on: {
            CONTINUE: 'feature',
          },
        },
        feature: {},
        init: {
          key: 'init',
          initial: 'setup',
          on: { DONE: 'select' },
          states: {
            setup: {
              on: { CONTINUE: 'save' },
              onEntry: { type: 'render', component: 'wizard/init-setup' },
            },
            save: {
              on: { CONTINUE: 'unseal' },
              onEntry: { type: 'render', component: 'wizard/init-save-keys' },
            },
            unseal: {
              on: { CONTINUE: 'login' },
              onEntry: { type: 'render', component: 'wizard/init-unseal' },
            },
            login: {
              onEntry: { type: 'render', component: 'wizard/init-login' },
            },
          },
        },
      },
    },
    idle: {
      on: {
        DISMISS: 'dismissed',
        CONTINUE: 'active',
      },
    },
    dismissed: {
      on: { CONTINUE: 'idle' },
      onEntry: ['handleDismissed'],
    },
    paused: {
      on: { CONTINUE: ['handlePause'] },
    },
    complete: {
      on: {
        CONTINUE: 'idle',
        DISMISS: 'dismissed',
      },
    },
  },
};
